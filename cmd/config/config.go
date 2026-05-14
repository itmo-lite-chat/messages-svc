package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL" env-required:"true"`

	GrpcServer grpcServer

	DBConfig mongoConfig
}

type grpcServer struct {
	Host string `env:"GRPC_HOST" env-default:"localhost"`
	Port int    `env:"GRPC_PORT" env-default:"9999"`
}

type mongoConfig struct {
	Host     string `env:"MONGO_HOST"     env-required:"true"`
	Port     int    `env:"MONGO_PORT"     env-required:"true"`
	DB       string `env:"MONGO_DB"       env-required:"true"`
	User     string `env:"MONGO_USER"     env-required:"true"`
	Password string `env:"MONGO_PASSWORD" env-required:"true"`

	MongoMaxPoolSize int           `env:"MONGO_MAX_POOL_SIZE" env-required:"true"`
	MongoMinPoolSize int           `env:"MONGO_MIN_POOL_SIZE" env-required:"true"`
	MongoConnTimeout time.Duration `env:"MONGO_CONN_TIMEOUT"  env-required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, errors.Wrap(err, "config error")
	}

	return &cfg, nil
}
