# GoLang Delayed Job

## Overview

GoLang Delayed Job is a simple and efficient job scheduling library for Go (Golang). It enables you to schedule tasks to
run at a later time, providing a way to handle delayed and background jobs in your applications.

## Features
- Job Scheduling: Schedule tasks to run at a specified time in the future. 
- Background Jobs: Execute tasks in the background without blocking the main application flow. 
- Retry Mechanism: Handle job failures by implementing a retry mechanism. 
- Easy Integration: Seamless integration into your Go projects with a simple API.

## Idea
I use Redis, because Redis has a [SortedSet](https://redis.io/docs/data-types/sorted-sets/) data type, which is a collection where the contents must be unique, and the data position will always be sorted based on the value we want.
 The values entered in the SortedSet are the execution time to be run, and related data. 

Example:
```text
1703950502729 {"id": "57"}
```

Next, to get the data, I check with a predetermined time interval, for example, every 3 seconds will get any job that must be executed.

## Usage

```go
type expireOrder struct {
}

// Implement contract.JobProcessor interface
func NewExpireOrder() contract.JobProcessor {
	return &expireOrder{}
}

func (o *expireOrder) Serve(ctx context.Context, data *appctx.WatcherData) error {
	data.Commit() // Acknowledge data, so it does not perform the retry mechanism
}
```

```go
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
```

## Example
Check out the [example](cmd/job/expire_order.go) for a demonstrating the usage of GoLang Delayed Job.

### How to Run
```bash
$ docker compose up -d
$ go run main.go playground # Generate jobs to demonstrate
$ go run main.go job:expire-order # Run job watcher
```