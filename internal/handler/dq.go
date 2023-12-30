package handler

import (
	"context"
	"go-delayed-job/internal/appctx"
	ucaseContract "go-delayed-job/internal/ucase/contract"
	"go-delayed-job/pkg/dq"
)

func DQWatcherHandler(jobHandler ucaseContract.JobProcessor) dq.JobProcessorFunc {
	return func(decoder *dq.JobDecoder) {
		err := jobHandler.Serve(context.Background(), &appctx.WatcherData{
			Name:  decoder.Name,
			Key:   decoder.Key,
			Value: decoder.Value,
			Commit: func() {
				decoder.Commit(decoder)
			},
		})
		if err != nil {
			return
		}
	}
}
