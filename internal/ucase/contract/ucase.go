package contract

import (
	"context"
	"go-delayed-job/internal/appctx"
)

type JobProcessor interface {
	Serve(ctx context.Context, data *appctx.WatcherData) error
}
