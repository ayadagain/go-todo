package ctx

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type ServiceCtx interface {
	Conf() *Config
	GrpcClient() *grpc.ClientConn
}

type defaultServiceCtx struct {
	conf     *Config
	grpcConn *grpc.ClientConn
}

func (ctx *defaultServiceCtx) Conf() *Config {
	return ctx.conf
}

func (ctx *defaultServiceCtx) GrpcClient() *grpc.ClientConn {
	return ctx.grpcConn
}

func NewGrpcConn(ctx Config) *grpc.ClientConn {
	conn, err := grpc.NewClient(ctx.grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	return conn
}

func NewDefaultServiceContext() ServiceCtx {
	config := loadConf()
	hgrpc := NewGrpcConn(*config)

	return &defaultServiceCtx{
		conf:     config,
		grpcConn: hgrpc,
	}
}
