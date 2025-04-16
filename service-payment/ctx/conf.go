package ctx

const (
	paymentGrpcPort  = "PAYMENT_GRPC_PORT"
	kafkaServer      = "KAFKA_SERVER"
	kafkaGroupId     = "KAFKA_GROUP_ID"
	kafkaOffsetReset = "KAFKA_OFFSET_RESET"
	kafkaTopic       = "KAFKA_TOPIC"
	mongoUri         = "MONGO_URI"
)

type Config struct{}
