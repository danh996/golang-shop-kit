package grpc_server

import (
	"context"
	"expvar"
	"fmt"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	consul "github.com/hashicorp/consul/api"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

const MaxMsgSize = 1024 * 1024 * 1024

type GRPCServer struct {
	ServiceName string
	Host        string
	Port        string

	consul       *consul.Client
	tracer       opentracing.Tracer
	authFunction grpc_auth.AuthFunc
	Server       *grpc.Server
}

// NewGRPCServer ...
func NewGRPCServer(serviceName string, host string, port string) *GRPCServer {
	return &GRPCServer{
		ServiceName: serviceName,
		Host:        host,
		Port:        port,
	}
}

func (s *GRPCServer) InitServer() *GRPCServer {
	filter := grpc_opentracing.WithFilterFunc(func(ctx context.Context, fullMethodName string) bool {
		return fullMethodName != "/grpc.health.v1.Health/Check"
	})

	s.Server = grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(s.tracer), filter),
			grpc_prometheus.StreamServerInterceptor,
			// grpc_zap.StreamServerInterceptor(zapLogger),

			// auth function in middleware file
			// grpc_auth.StreamServerInterceptor(s.authFunction),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(s.tracer), filter),
			grpc_prometheus.UnaryServerInterceptor,
			// grpc_zap.UnaryServerInterceptor(zapLogger),
			// grpc_auth.UnaryServerInterceptor(s.authFunction),
		)),
		grpc.MaxRecvMsgSize(MaxMsgSize),
		grpc.MaxSendMsgSize(MaxMsgSize),
		grpc.MaxMsgSize(MaxMsgSize),
	)
	return s
}

func (s *GRPCServer) WithConsul(consul *consul.Client) *GRPCServer {
	s.consul = consul
	return s
}

func (s *GRPCServer) WithTracer(tracer opentracing.Tracer) *GRPCServer {
	s.tracer = tracer
	return s
}

func (s *GRPCServer) WithAuthFunc(authFunction grpc_auth.AuthFunc) *GRPCServer {
	s.authFunction = authFunction
	return s
}

func (s *GRPCServer) EnablePrometheus(port string) *GRPCServer {
	grpc_prometheus.Register(s.Server)
	mux := http.NewServeMux()

	mux.Handle("/debug/vars", expvar.Handler())
	mux.Handle("/metrics", promhttp.Handler())

	go func() {
		http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	}()

	return s
}
