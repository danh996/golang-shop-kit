package main

const (
	isEnableHystrix   = true
	serviceName       = "Client"
	serverServiceName = "Server"
	port              = 11011
	consulAddress     = "localhost:8500"
	tracerAddress     = "localhost:6831"
	endpoint          = "localhost:8282"
)

// func main() {
// 	ctx := context.Background()
// 	// client to connect consul server
// 	consulClient, err := registry.NewClient(consulAddress)
// 	if err != nil {
// 		panic(err)
// 	}

// 	tracer, err := tracing.NewOpenTracer(serviceName, tracerAddress)
// 	if err != nil {
// 		panic(err)
// 	}

// 	client := grpc_client.NewGRPCClient(
// 		isEnableHystrix,
// 		serverServiceName,
// 		tracer,
// 		consulClient,
// 		consulAddress,
// 		endpoint,
// 	)
// 	conn, err := client.DialEndpoint()

// 	if err != nil {
// 		panic(err)
// 	}

// 	caculatorClient := pb.NewCaculatorServiceClient(conn)

// 	healthCheckClient := pbHealth.NewHealthClient(conn)

// 	// caculate random with 5 second
// 	for {

// 		result, err := caculatorClient.Add(ctx, &pb.Request{
// 			A: 1,
// 			B: 2,
// 		})

// 		if err != nil {
// 			panic(err)
// 		}
// 		if _, err := healthCheckClient.Check(ctx, &pbHealth.HealthCheckRequest{}); err != nil {
// 			panic(err)
// 		} else {
// 			fmt.Println("healcheck ok")
// 		}

// 		fmt.Println(timestamppb.Now())

// 		if _, err := caculatorClient.AddDate(ctx, &pb.AddDateRequest{
// 			Time: timestamppb.Now(),
// 		}); err != nil {
// 			panic(err)
// 		}

// 		fmt.Printf("result: %d\n", result.Result)
// 		time.Sleep(5 * time.Second)
// 	}
// }
