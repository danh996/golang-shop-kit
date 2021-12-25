package main

import (
	"context"
	"fmt"
	"net"

	"gitlab.com/canco1/canco-kit/example/pb"
	"gitlab.com/canco1/canco-kit/grpc_server"
	"gitlab.com/canco1/canco-kit/registry"
	"gitlab.com/canco1/canco-kit/tracing"

	pbHealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	isEnableHystrix = true
	serviceName     = "Server"
	host            = ""
	port            = 11011
	consulAddress   = "localhost:8500"
	tracerAddress   = "localhost:6831"
	prometheusPort  = 11013
)

func main() {
	// client to connect consul server
	consulClient, err := registry.NewClient(consulAddress)
	if err != nil {
		panic(err)
	}

	tracer, err := tracing.NewOpenTracer(serviceName, tracerAddress)
	if err != nil {
		panic(err)
	}
	server := grpc_server.NewGRPCServer(serviceName, host, string(port))

	server.
		WithConsul(consulClient).
		WithTracer(tracer.Tracer).
		// WithAuthFunc(grpc_server.MappingRequestInfo).
		// EnablePrometheus(string(prometheusPort)).
		InitServer()

	delivery := NewDelivery()
	pbHealth.RegisterHealthServer(server.Server, grpc_server.NewHealthService())
	pb.RegisterCaculatorServiceServer(server.Server, delivery)

	register := registry.NewConsulRegister(consulClient, serviceName, port, []string{
		"Server",
	})

	if _, err := register.Register(); err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("%s:%d", host, port))
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}

	reflection.Register(server.Server)

	fmt.Printf("Server listening in port: %d", port)

	if err := server.Server.Serve(ln); err != nil {
		panic(err)
	}
}

type Delivery struct {
	pb.UnimplementedCaculatorServiceServer
}

func (d *Delivery) Add(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	fmt.Println(req.GetA())
	return &pb.Response{
		Result: req.GetA() + req.GetB(),
	}, nil
}

func (d *Delivery) AddDate(ctx context.Context, req *pb.AddDateRequest) (*emptypb.Empty, error) {
	fmt.Printf("%v", req.Time.AsTime())
	return &emptypb.Empty{}, nil
}

func NewDelivery() pb.CaculatorServiceServer {
	return &Delivery{}
}
