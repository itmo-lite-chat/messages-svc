package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL" env-required:"true"`

	GrpcServer grpcServer

	DBConfig dbConfig
}

type grpcServer struct {
	Host string `env:"GRPC_HOST" env-default:"localhost"`
	Port int    `env:"GRPC_PORT" env-default:"9999"`
}

// DBConfig настройки базы данных
type dbConfig struct {
	Host     string `env:"POSTGRES_HOST"     env-required:"true"`
	Port     int    `env:"POSTGRES_PORT"     env-required:"true"`
	Name     string `env:"POSTGRES_DB"       env-required:"true"`
	Username string `env:"POSTGRES_USER"     env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`

	ApplyMigration bool `env:"APPLY_MIGRATION" env-required:"true"`

	ConnectionTTL      int `env:"POSTGRES_CONN_TTL"      env-required:"true"`
	MaxOpenConnections int `env:"POSTGRES_MAX_OPEN_CONN" env-required:"true"`
	MaxIdleConnections int `env:"POSTGRES_MAX_IDLE_CONN" env-required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, errors.Wrap(err, "config error")
	}

	return &cfg, nil
}
