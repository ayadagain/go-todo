package ctx

import "os"

const (
	todoGrpcPort    = "TODO_GRPC_PORT"
	kafkaServer     = "KAFKA_SERVER"
	paymentGrpcHost = "PAYMENT_GRPC_HOST"
	kafkaTopic      = "KAFKA_TOPIC"
)

type Config struct {
	todoGrpcPort    string
	kafkaServer     string
	paymentGrpcHost string
	kafkaTopic      string
}

func loadConfig() *Config {
	return &Config{
		todoGrpcPort:    os.Getenv(todoGrpcPort),
		kafkaServer:     os.Getenv(kafkaServer),
		paymentGrpcHost: os.Getenv(paymentGrpcHost),
		kafkaTopic:      os.Getenv(kafkaTopic),
	}
}

func (c *Config) getTodoGrpcPort() string {
	return c.todoGrpcPort
}

func (c *Config) getKafkaServer() string {
	return c.kafkaServer
}

func (c *Config) getPaymentGrpcHost() string {
	return c.paymentGrpcHost
}

func (c *Config) GetKafkaTopic() string {
	return c.kafkaTopic
}
