package job

import (
	"context"
	"go-delayed-job/internal/consts"
	"go-delayed-job/internal/handler"
	"go-delayed-job/internal/ucase/job"

	"go-delayed-job/pkg/dq"
	"go-delayed-job/pkg/redis"
)

func RunJobExpireOrder(ctx context.Context) {
	redisClient := redis.NewRedis()
	dqWatcher := dq.NewWatcher(redisClient)

	ucase := job.NewExpireOrder()

	dqWatcher.Watch(&dq.WatcherContext{
		Handler:  handler.DQWatcherHandler(ucase),
		Context:  ctx,
		Name:     consts.OrderExpirationJobQueueName,
		Interval: consts.OrderExpirationJobInterval,
		Limit:    consts.OrderExpirationJobLimit,
	})
}
