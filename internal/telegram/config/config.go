package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	ErrMissing  = errors.New("missing")
	ErrBadValue = errors.New("bad value")
)

type Config struct {
	Subscribers            []int64 `toml:"subscribers"`
	BotTokenEnvariableName string  `toml:"botTokenEnvariableName"`
	BotToken               string  `json:"-"` // manual

	Timeout time.Duration `toml:"timeout"`
}

func (c *Config) Validate() error {
	var errs error

	if len(c.Subscribers) == 0 {
		errs = errors.Join(errs, fmt.Errorf("subscribers %w", ErrMissing))
	}
	if c.BotToken == "" {
		errs = errors.Join(errs, fmt.Errorf("bot token %w", ErrMissing))
	}
	if c.Timeout < 0 {
		errs = errors.Join(errs, fmt.Errorf("timeout must be at least zero, get %s: %w",
			c.Timeout.String(),
			ErrBadValue,
		))
	}

	return errs
}

func (c *Config) PostLoad() error {
	if c.BotTokenEnvariableName == "" {
		return fmt.Errorf("botTokenEnvariableName %w", ErrMissing)
	}

	var ok bool
	c.BotToken, ok = os.LookupEnv(c.BotTokenEnvariableName)
	if !ok {
		return fmt.Errorf("envariable %s %w", c.BotTokenEnvariableName, ErrMissing)
	}

	return nil
}

func New(path string) (*Config, error) {
	conf := new(Config)
	err := cleanenv.ReadConfig(path, conf)
	if err != nil {
		return nil, fmt.Errorf("can't read config: %w", err)
	}

	err = conf.PostLoad()
	if err != nil {
		return nil, fmt.Errorf("can't do postload: %w", err)
	}

	err = conf.Validate()
	if err != nil {
		return nil, fmt.Errorf("not valid: %w", err)
	}

	return conf, nil
}
