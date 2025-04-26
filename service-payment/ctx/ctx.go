package ctx

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type KafkaConsumer interface {
	Run()
	Shutdown()
}

type ServiceCtx interface {
	Config() *Config
	KafkaConsumer() *kafka.Consumer
	TcpListener() net.Listener
	MongoClient() *mongo.Client
	MongoCollection() *mongo.Collection
	Start(*grpc.Server, ...KafkaConsumer)
	ShutdownHook(shutdownFuncs ...func() error)
}

type defaultServiceCtx struct {
	config          *Config
	kafkaConsumer   *kafka.Consumer
	mongoClient     *mongo.Client
	mongoCollection *mongo.Collection
	tcpListener     net.Listener
}

func (ctx *defaultServiceCtx) Config() *Config {
	return ctx.config

}

func (ctx *defaultServiceCtx) KafkaConsumer() *kafka.Consumer {
	return ctx.kafkaConsumer
}

func (ctx *defaultServiceCtx) MongoClient() *mongo.Client {
	return ctx.mongoClient
}

func (ctx *defaultServiceCtx) MongoCollection() *mongo.Collection {
	return ctx.mongoCollection
}

func (ctx *defaultServiceCtx) TcpListener() net.Listener {
	return ctx.tcpListener
}

func (ctx *defaultServiceCtx) Start(server *grpc.Server, consumers ...KafkaConsumer) {
	fmt.Println("Serving gRPC Server on :", ctx.tcpListener.Addr().String())
	go func() {
		err := server.Serve(ctx.tcpListener)
		if err != nil {
			fmt.Println("Error serving gRPC Server:", err)
		}
	}()

	for _, consumer := range consumers {
		go consumer.Run()
	}
}

func (ctx *defaultServiceCtx) ShutdownHook(shutdownFuncs ...func() error) {
	func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		ctx.shutdown()
		for _, f := range shutdownFuncs {
			err := f()
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()
}

func NewKafkaConn(ctx Config) *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": ctx.kafkaServer,
		"group.id":          ctx.kafkaGroupId,
		"auto.offset.reset": ctx.kafkaOffsetReset,
	})

	if err != nil {
		panic("failed to create kafka consumer: " + err.Error())
	}

	metadata, err := c.GetMetadata(nil, true, 5000)
	if err != nil {
		fmt.Printf("Failed to get Kafka metadata: %s\n", err.Error())
	} else {
		fmt.Printf("Successfully connected to Kafka cluster with %d broker(s)\n", len(metadata.Brokers))
	}

	return c
}

func NewMongoConn(ctx Config) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(ctx.mongoUri))

	if err != nil {
		panic("failed to create MongoDB client: " + err.Error())
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		panic("failed to ping MongoDB: " + err.Error())
	}

	return client
}

func NewTcpListener(ctx Config) net.Listener {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", ctx.grpcPort))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	return listener
}

func NewDefaultServiceCtx() ServiceCtx {
	config := loadConfig()

	kafkaConsumer := NewKafkaConn(*config)
	mongoClient := NewMongoConn(*config)
	localMongoCollection := mongoClient.Database(config.mongoDatabase).Collection(config.mongoCollection)
	tcpListener := NewTcpListener(*config)

	err := kafkaConsumer.SubscribeTopics([]string{config.kafkaTopic}, nil)

	if err != nil {
		panic(err)
	}

	return &defaultServiceCtx{
		config:          config,
		kafkaConsumer:   kafkaConsumer,
		mongoClient:     mongoClient,
		mongoCollection: localMongoCollection,
		tcpListener:     tcpListener,
	}
}

func (ctx *defaultServiceCtx) shutdown() {
	if ctx.mongoClient != nil {
		err := ctx.mongoClient.Disconnect(context.Background())
		if err != nil {
			log.Println("error while trying to disconnect from mongo: ", err)
		}
		if ctx.kafkaConsumer != nil {
			err = ctx.kafkaConsumer.Close()

			if err != nil {
				log.Println("error while trying to close kafka consumer: ", err)
			}
		}
	}
}
