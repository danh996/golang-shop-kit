package grpc_client

import (
	"fmt"
	"time"

	"gitlab.com/canco1/canco-kit/registry"
	"gitlab.com/canco1/canco-kit/tracing"

	consul "github.com/hashicorp/consul/api"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/codes"
)

type GRPCClient struct {
	isEnableHystrix bool
	serviceName     string
	tracer          *tracing.OpenTracer
	consul          *consul.Client
	consulAddress   string
	endpoint        string
}

func NewGRPCClient(
	isEnableHystrix bool,
	tracer *tracing.OpenTracer,
	endpoint string,
) *GRPCClient {
	return &GRPCClient{
		isEnableHystrix: isEnableHystrix,
		tracer:          tracer,
		endpoint:        endpoint,
	}
}

func (c *GRPCClient) DialWithConsul(consul *consul.Client, consulAddress string) (*grpc.ClientConn, error) {
	optsRetry := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(50 * time.Millisecond)),
		grpc_retry.WithCodes(codes.Unavailable),
		grpc_retry.WithMax(3),
		grpc_retry.WithPerRetryTimeout(3 * time.Second),
	}

	sIntOpt := grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
		grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(c.tracer.Tracer)),
		grpc_prometheus.StreamClientInterceptor,
		grpc_retry.StreamClientInterceptor(optsRetry...),
	))

	grpc_prometheus.EnableClientHandlingTimeHistogram()

	uIntOpt := grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
		UnaryClientInterceptor(c.isEnableHystrix),
		grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(c.tracer.Tracer)),
		grpc_prometheus.UnaryClientInterceptor,
		grpc_retry.UnaryClientInterceptor(optsRetry...),
	))

	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.FailFast(false)),
		sIntOpt,
		uIntOpt,
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
	}
	registry.Init()
	return grpc.Dial(fmt.Sprintf("%s://%s/%s", "consul", c.consulAddress, c.serviceName), opts...)
}

func (c *GRPCClient) Dial() (*grpc.ClientConn, error) {
	optsRetry := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(50 * time.Millisecond)),
		grpc_retry.WithCodes(codes.Unavailable),
		grpc_retry.WithMax(3),
		grpc_retry.WithPerRetryTimeout(3 * time.Second),
	}

	sIntOpt := grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
		grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(c.tracer.Tracer)),
		grpc_prometheus.StreamClientInterceptor,
		grpc_retry.StreamClientInterceptor(optsRetry...),
	))

	grpc_prometheus.EnableClientHandlingTimeHistogram()

	uIntOpt := grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
		UnaryClientInterceptor(c.isEnableHystrix),
		grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(c.tracer.Tracer)),
		grpc_prometheus.UnaryClientInterceptor,
		grpc_retry.UnaryClientInterceptor(optsRetry...),
	))

	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.FailFast(false)),
		sIntOpt,
		uIntOpt,
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
	}
	registry.Init()
	return grpc.Dial(c.endpoint, opts...)
}

func (c *GRPCClient) ClientStreamDial() (*grpc.ClientConn, error) {
	optsRetry := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(50 * time.Millisecond)),
		grpc_retry.WithCodes(codes.Unavailable),
		grpc_retry.WithMax(3),
		grpc_retry.WithPerRetryTimeout(3 * time.Second),
	}

	sIntOpt := grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
		grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(c.tracer.Tracer)),
		grpc_prometheus.StreamClientInterceptor,
	))

	grpc_prometheus.EnableClientHandlingTimeHistogram()

	uIntOpt := grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
		UnaryClientInterceptor(c.isEnableHystrix),
		grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(c.tracer.Tracer)),
		grpc_prometheus.UnaryClientInterceptor,
		grpc_retry.UnaryClientInterceptor(optsRetry...),
	))

	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.FailFast(false)),
		sIntOpt,
		uIntOpt,
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
	}
	registry.Init()
	return grpc.Dial(c.endpoint, opts...)
}
