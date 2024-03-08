package config

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrBadValue = errors.New("bad value")
)

type Config struct {
	Retries    int           `toml:"retries"`
	RetryDelay time.Duration `toml:"retryDelay"`
	Timeout    time.Duration `toml:"timeout"`
}

func (c *Config) Validate() error {
	var errs error

	if c.Retries < 0 {
		errs = errors.Join(errs, fmt.Errorf("retries must be at least 0, get %d: %w",
			c.Retries,
			ErrBadValue,
		))
	}
	if c.RetryDelay <= 0 {
		errs = errors.Join(errs, fmt.Errorf("retryDelay must be greater than zero, get %s: %w",
			c.RetryDelay.String(),
			ErrBadValue,
		))
	}
	if c.Timeout < 0 {
		errs = errors.Join(errs, fmt.Errorf("timeout must be at least zero, get %s: %w",
			c.Timeout.String(),
			ErrBadValue,
		))
	}

	return errs
}
