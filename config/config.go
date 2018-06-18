package config

import (
	"time"
	"github.com/kelseyhightower/envconfig"
)

type Timeouts struct {
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `envconfig:"IDLE_TIMEOUT"`
}

type Config struct {
	Env  string `envconfig:"ENV"`
	Port uint   `envconfig:"PORT"`

	Timeouts Timeouts
}

// NewConfig returns a new empty config
func NewConfig() *Config {
	return &Config{
		// We can add defaults here
		// You can use the envconfig library to set the defaults as well
		// But it only works for Exported fields
		Env:  "development",
		Port: 5000,
		Timeouts: Timeouts{
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
			IdleTimeout:  time.Second * 60,
		},
	}
}

// LoadFromEnv load configuration from env vars
func (c *Config) LoadFromEnv() error {
	return envconfig.Process("", c)
}
