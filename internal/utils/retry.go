package utils

import (
	_ "github.com/cenkalti/backoff"
	"github.com/cenkalti/backoff/v4"
	"time"
)

type Retry struct {
	InitInterval time.Duration
	RetryTimeout time.Duration

	MaxRetryTimes uint64
}

func (r *Retry) NewExponentialBackOff() backoff.BackOff {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = r.InitInterval
	bo.MaxElapsedTime = r.RetryTimeout

	return backoff.WithMaxRetries(bo, r.MaxRetryTimes)
}
