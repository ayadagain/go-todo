package ctx

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type KafkaProducer interface {
	Run()
	Shutdown() func()
}

type ServiceCtx interface {
	Config() *Config
	KafkaProducer() *kafka.Producer
	TcpListener() net.Listener
	Start(*grpc.Server, KafkaProducer)
	ShutdownHook(shutdownFuncs ...func())
	GrpcClient() *grpc.ClientConn
}

type defaultServiceCtx struct {
	config        *Config
	kafkaProducer *kafka.Producer
	tcpListener   net.Listener
	grpcClient    *grpc.ClientConn
}

func (ctx *defaultServiceCtx) Config() *Config {
	return ctx.config
}

func (ctx *defaultServiceCtx) GrpcClient() *grpc.ClientConn {
	return ctx.grpcClient
}

func (ctx *defaultServiceCtx) KafkaProducer() *kafka.Producer {
	return ctx.kafkaProducer
}

func (ctx *defaultServiceCtx) TcpListener() net.Listener {
	return ctx.tcpListener
}

func (ctx *defaultServiceCtx) Start(server *grpc.Server, kafkaProducer KafkaProducer) {
	fmt.Println("Serving gRPC Server on :", ctx.tcpListener.Addr().String())
	go func() {
		err := server.Serve(ctx.tcpListener)
		if err != nil {
			fmt.Println("Error serving gRPC Server:", err)
		}
	}()

	go kafkaProducer.Run()
}

func (ctx *defaultServiceCtx) ShutdownHook(shutdownFuncs ...func()) {
	func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		ctx.shutdown()
		for _, f := range shutdownFuncs {
			f()
		}
	}()
}

func (ctx *defaultServiceCtx) shutdown() {
	err := ctx.tcpListener.Close()
	if err != nil {
		return
	}
}

func NewKafkaConn(ctx Config) *kafka.Producer {
	prod, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": ctx.kafkaServer,
	})

	if err != nil {
		log.Fatalf("failed to create kafka producer: %v", err)
	}

	return prod
}

func NewTcpListener(ctx Config) net.Listener {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", ctx.todoGrpcPort))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	return listener
}

func NewGrpcConn(ctx Config) *grpc.ClientConn {
	conn, err := grpc.NewClient(ctx.paymentGrpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	return conn
}

func NewDefaultServiceCtx() ServiceCtx {
	config := loadConfig()
	kafkaProducer := NewKafkaConn(*config)
	tcpListener := NewTcpListener(*config)
	grpcClient := NewGrpcConn(*config)

	return &defaultServiceCtx{
		config:        config,
		kafkaProducer: kafkaProducer,
		tcpListener:   tcpListener,
		grpcClient:    grpcClient,
	}
}
