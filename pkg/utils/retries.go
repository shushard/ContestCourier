package utils

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

func WithRetries(
	ctx context.Context,
	retriesCount uint,
	timeSleep time.Duration,
	function func() error,
	logger *slog.Logger,
) error {
	var errs error
	var iteration uint
	for {
		newErr := function()
		if newErr == nil {
			return nil
		}
		if !errors.Is(errs, newErr) {
			errs = errors.Join(errs, newErr)
		}

		iteration++
		if iteration >= retriesCount && retriesCount != 0 {
			return fmt.Errorf("still getting errors after %d retries, errors: %w", retriesCount, errs)
		}

		logger.Warn("RETRIES Get error, wait for next retry...",
			"error", newErr,
			"try", iteration,
			"delay", timeSleep,
		)

		timer := time.NewTimer(timeSleep)
		select {
		case <-ctx.Done():
			return context.Canceled
		case <-timer.C:
		}
	}
}
