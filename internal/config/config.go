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

// status const

const (
	REGISTERED = "REGISTERED"
	INVALID    = "INVALID"
	PROCESSING = "PROCESSING"
	PROCESSED  = "PROCESSED"
)

// flags

var Flags struct {
	ServerAddress        string
	DatabaseDSN          string
	AccrualSystemAddress string
}

// env vars

var EnvVar struct {
	ServerAddress        string `env:"RUN_ADDRESS" envDefault:":9090"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://127.0.0.1:8080"`
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

func NewConfig(srvAddr, databaseDSN, accrualSystemAddress string) *Config {
	cfg := make(map[string]string)
	cfg["server_address_str"] = srvAddr
	cfg["database_dsn_str"] = databaseDSN
	cfg["accrual_system_address"] = accrualSystemAddress
	return &Config{
		cfg: cfg,
	}
}
