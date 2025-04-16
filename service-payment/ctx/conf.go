package ctx

import "os"

const (
	paymentGrpcPort  = "PAYMENT_GRPC_PORT"
	kafkaServer      = "KAFKA_SERVER"
	kafkaGroupId     = "KAFKA_GROUP_ID"
	kafkaOffsetReset = "KAFKA_OFFSET_RESET"
	kafkaTopic       = "KAFKA_TOPIC"
	mongoUri         = "MONGO_URI"
)

type Config struct {
	paymentGrpcPort  string
	kafkaServer      string
	kafkaGroupId     string
	kafkaOffsetReset string
	kafkaTopic       string
	mongoUri         string
}

func loadConfig() *Config {
	return &Config{
		paymentGrpcPort:  os.Getenv(paymentGrpcPort),
		kafkaServer:      os.Getenv(kafkaServer),
		kafkaGroupId:     os.Getenv(kafkaGroupId),
		kafkaOffsetReset: os.Getenv(kafkaOffsetReset),
		kafkaTopic:       os.Getenv(kafkaTopic),
		mongoUri:         os.Getenv(mongoUri),
	}
}

func (c *Config) getPaymentGrpcPort() string {
	return c.paymentGrpcPort
}

func (c *Config) getKafkaServer() string {
	return c.kafkaServer
}

func (c *Config) getKafkaGroupId() string {
	return c.kafkaGroupId
}

func (c *Config) getKafkaOffsetReset() string {
	return c.kafkaOffsetReset
}

func (c *Config) getKafkaTopic() string {
	return c.kafkaTopic
}

func (c *Config) getMongoUri() string {
	return c.mongoUri
}
