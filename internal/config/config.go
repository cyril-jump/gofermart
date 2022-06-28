package config

import (
	"github.com/cyril-jump/gofermart/internal/utils/errs"
)

// context const
type contextKey string

const (
	TokenKey = contextKey("token")
)

func (c contextKey) String() string {
	return string(c)
}

// flags

var Flags struct {
	ServerAddress        string
	DatabaseDSN          string
	AccrualSystemAddress string
}

// env vars

var EnvVar struct {
	ServerAddress        string `env:"RUN_ADDRESS" envDefault:":8080"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseDSN          string `env:"DATABASE_URI" envDefault:"postgres://dmosk:dmosk@localhost:5432/dmosk?sslmode=disable"`
}

// config

type Config struct {
	cfg map[string]string
}

func (c Config) Get(key string) (string, error) {
	if _, ok := c.cfg[key]; !ok {
		return "", errs.ErrNotFound
	}
	return c.cfg[key], nil
}

//constructor

func NewConfig(srvAddr, databaseDSN string) *Config {
	cfg := make(map[string]string)
	cfg["server_address_str"] = srvAddr
	cfg["database_dsn_str"] = databaseDSN
	return &Config{
		cfg: cfg,
	}
}
