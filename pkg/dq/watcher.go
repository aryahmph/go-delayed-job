package dq

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

type watcher struct {
	client *redis.Client
}

func NewWatcher(client *redis.Client) Watcher {
	return &watcher{
		client: client,
	}
}

func (dq *watcher) Watch(ctx *WatcherContext) {
	ticker := time.NewTicker(ctx.Interval)
	defer func() {
		ticker.Stop()
	}()

	newCtx, cancel := context.WithCancel(ctx.Context)

	go func() {
		for {
			select {
			case <-newCtx.Done():
				slog.Warn(fmt.Sprintf("[delay queue] stopped watch %s queue", ctx.Name))
				return
			case <-ticker.C:
				keys, values, err := dq.pulls(newCtx, ctx.Name, ctx.Limit)
				if err != nil {
					slog.Error(fmt.Sprintf("[delay queue] queue %s error pull jobs, err: %v", ctx.Name, err))
					continue
				}

				for i, key := range keys {
					go func(newCtx context.Context, key int64, value string) {
						ctx.Handler(&JobDecoder{
							Name:  ctx.Name,
							Key:   key,
							Value: value,
							Commit: func(decoder *JobDecoder) {
								err := dq.remove(newCtx, ctx.Name, decoder.Value)
								if err != nil {
									slog.Error(fmt.Sprintf("[delay queue] queue %s error remove jobs, err: %v", ctx.Name, err))
								}
							},
						})
					}(newCtx, key, values[i])
				}

			}
		}
	}()

	slog.Info(fmt.Sprintf("[delay queue] watch %s queue is running", ctx.Name))

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm

	cancel()
	slog.Info(fmt.Sprintf("[delay queue] cancelled watch %s queue", ctx.Name))
}

func (dq *watcher) pulls(ctx context.Context, queueName string, limit int64) (keys []int64, values []string, err error) {
	now := time.Now().UnixMilli()
	members, err := dq.client.ZRangeByScoreWithScores(ctx, queueName,
		&redis.ZRangeBy{
			Min:   "-inf",
			Max:   fmt.Sprintf("%d", now),
			Count: limit,
		},
	).Result()
	if err != nil {
		return keys, values, err
	}

	for _, member := range members {
		keys = append(keys, int64(member.Score))
		values = append(values, member.Member)
	}

	return
}

func (dq *watcher) remove(ctx context.Context, queueName string, value string) error {
	return dq.client.ZRem(ctx, queueName, value).Err()
}
