package utils

import (
	"github.com/cenkalti/backoff/v4"
	"time"
)

// RetryConfig configuration for the retry mechanism such as retry interval, randomization factor, multiplier
func RetryConfig() (exponentialBackOffInstance *backoff.ExponentialBackOff) {
	exponentialBackOffInstance = backoff.NewExponentialBackOff()

	exponentialBackOffInstance.InitialInterval = 500 * time.Millisecond
	exponentialBackOffInstance.Multiplier = 1.5
	exponentialBackOffInstance.RandomizationFactor = 0.5
	return
}
