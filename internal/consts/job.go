package consts

import "time"

const (
	OrderExpirationJobQueueName = "order-expiration"
	OrderExpirationJobInterval  = time.Duration(2) * time.Second
	OrderExpirationJobLimit     = 100
)
