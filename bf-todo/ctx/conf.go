package ctx

import (
	"os"
)

const (
	secretKey = "SECRET_KEY"
	grpcHost  = "TODO_GRPC_HOST"
)

type Config struct {
	secretKey string
	grpcHost  string
}

func loadConf() *Config {
	return &Config{
		secretKey: os.Getenv(secretKey),
		grpcHost:  os.Getenv(grpcHost),
	}
}

func (c *Config) GetSecretKey() string {
	return c.secretKey
}
