package retry

import (
	"time"
)

func WithExponentialBackoff(funcToRetry func() error, initialDelay time.Duration, maxWait time.Duration) error {
	return withExponentialBackoff(funcToRetry, initialDelay, maxWait, time.Sleep)
}

func withExponentialBackoff(funcToRetry func() error, initialDelay time.Duration, maxWait time.Duration, sleep func(delay time.Duration)) error {
	retryDelay := initialDelay
	startTime := time.Now()
	retries := 0

	for {
		err := funcToRetry()
		if err == nil {
			break
		}

		if time.Since(startTime) >= maxWait {
			return err
		}

		sleep(retryDelay)
		retries++
		retryDelay *= 2
	}

	return nil
}
