package dq

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type producer struct {
	client *redis.Client
}

func NewProducer(client *redis.Client) Producer {
	return &producer{client: client}
}

func (dq *producer) Add(ctx context.Context, job *JobContext) error {
	return dq.client.ZAdd(ctx, job.QueueName, redis.Z{
		Score:  float64(job.ExpiredAt),
		Member: job.Value,
	}).Err()
}
