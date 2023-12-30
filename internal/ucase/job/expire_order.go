package job

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-delayed-job/internal/appctx"
	"go-delayed-job/internal/presentations"
	ucaseContract "go-delayed-job/internal/ucase/contract"
	"log/slog"
	"math/rand"
)

type expireOrder struct {
}

func NewExpireOrder() ucaseContract.JobProcessor {
	return &expireOrder{}
}

func (o *expireOrder) Serve(ctx context.Context, data *appctx.WatcherData) error {
	attrs := []slog.Attr{
		slog.Any("event.name", "JobExpireOrder"),
		slog.Any("event.data", fmt.Sprintf("key: %d, val: %s", data.Key, data.Value)),
	}

	/*-------------------------------
	| STEP 1: Decode Request
	* -------------------------------*/
	var payload presentations.OrderExpirePayload
	err := json.Unmarshal([]byte(data.Value), &payload)
	if err != nil {
		attrs = append(attrs, slog.Any("event.status", "Failed"))
		slog.LogAttrs(ctx, slog.LevelError, fmt.Sprintf("%s Failed", "JobExpireOrder"), attrs...)

		return err
	}

	attrs = append(attrs, slog.Any("event.status", "Success"))

	/*-------------------------------
	| STEP 2: Try error
	* -------------------------------*/
	if rand.Intn(2) == 0 {
		attrs = append(attrs, slog.Any("event.status", "Failed"))
		slog.LogAttrs(ctx, slog.LevelError, fmt.Sprintf("%s Failed", "JobExpireOrder"), attrs...)

		return errors.New("something wrong")
	}
	attrs = append(attrs, slog.Any("event.status", "Success"))

	/*-------------------------------
	| STEP 3: ....
	* -------------------------------*/

	data.Commit() // Acknowledge data, so it does not perform the retry mechanism
	slog.LogAttrs(ctx, slog.LevelInfo, fmt.Sprintf("%s Success", "JobExpireOrder"), attrs...)

	return nil
}
