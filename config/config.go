package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Timeouts struct {
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"5s" desc:"tcp connection read timeout"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"5s" desc:"tcp connection write timeout"`
	IdleTimeout  time.Duration `envconfig:"IDLE_TIMEOUT" default:"60s" desc:"tcp connection idle timeout"`
}

type Config struct {
	Env           string `envconfig:"ENV" default:"development" required:"true" desc:"development, test or production"`
	Port          uint   `envconfig:"PORT" default:"5000" required:"true" desc:"HTTP port to listen"`
	DatabaseDSN   string `envconfig:"DATABASE_DSN" default:":memory:" required:"true" desc:"database connection DSN"`
	RunMigrations bool   `envconfig:"RUN_MIGRATIONS" default:"true" desc:"should run migrations on startup?"`

	Timeouts Timeouts
}

// FromEnv load configuration from env vars
func FromEnv() (*Config, error) {
	c := &Config{}
	return c, errors.WithStack(envconfig.Process("", c))
}
