package playground

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go-delayed-job/internal/consts"
	"go-delayed-job/internal/presentations"

	"go-delayed-job/pkg/dq"
	"go-delayed-job/pkg/redis"
)

func RunPlayground(ctx context.Context) {
	redisClient := redis.NewRedis()
	dqProducer := dq.NewProducer(redisClient)

	for i := 0; i < 500; i++ {
		order := presentations.OrderExpirePayload{ID: fmt.Sprintf("%d", i)}
		payload, _ := json.Marshal(order)

		randDuration := time.Duration(rand.Intn(60)) * time.Second

		err := dqProducer.Add(ctx, &dq.JobContext{
			QueueName: consts.OrderExpirationJobQueueName,
			Value:     string(payload),
			ExpiredAt: time.Now().Add(randDuration).UnixMilli(),
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
}
