package dq

import (
	"context"
	"time"
)

type (
	Producer interface {
		Add(ctx context.Context, job *JobContext) error
	}

	Watcher interface {
		Watch(ctx *WatcherContext)
	}
)

type (
	JobContext struct {
		QueueName string
		Value     string
		ExpiredAt int64
	}

	WatcherContext struct {
		Handler  JobProcessorFunc
		Context  context.Context
		Name     string
		Interval time.Duration
		Limit    int64
	}
)
