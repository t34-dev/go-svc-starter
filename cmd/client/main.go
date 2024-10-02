package main

import (
	"context"
	"fmt"
	grpcpool "github.com/t34-dev/go-grpc-pool"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
	"time"

	"github.com/t34-dev/go-svc-starter/pkg/api/common_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const isTSL = false

func main() {
	creds, err := credentials.NewClientTLSFromFile("cert/service.pem", "")
	if err != nil {
		log.Fatalf("failed to load client TLS credentials: %v", err)
	}
	opts := []grpc.DialOption{
		grpc.WithBlock(),
	}

	if isTSL {
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Define a factory function to create gRPC connections
	factory := func() (*grpc.ClientConn, error) {
		//opts := []grpc.DialOption{
		//	grpc.WithTransportCredentials(insecure.NewCredentials()),
		//}
		//return grpc.NewClient("localhost"+constants.Address, opts...)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		return grpc.DialContext(ctx, "localhost:50051", opts...)
	}

	// Create a new connection pool
	grpcPool, err := grpcpool.NewPool(factory, grpcpool.PoolOptions{
		MinConn: 2,
		MaxConn: 30,
	})
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}
	defer grpcPool.Close()
	conn, err := grpcPool.Get()
	// conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Free()

	c := common_v1.NewCommonV1Client(conn.GetConn())

	ctx := context.Background()

	// Requests
	timeResp, err := c.GetTime(ctx, &common_v1.TimeRequest{
		Error: false,
	})
	if err != nil {
		log.Fatalf("could not get time: %v", err)
	}
	fmt.Printf("Current time: %v\n", timeResp.GetTime().AsTime())

	stream, err := c.LongOperation(ctx, &common_v1.LongOperationRequest{})
	if err != nil {
		log.Fatalf("could not perform long operation: %v", err)
	}

	done := make(chan bool)
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf("error receiving response: %v", err)
			}
			fmt.Printf("Status: %v, Message: %s, Progress: %d%%\n", resp.Status, resp.Message, resp.Progress)
			if resp.Status == common_v1.LongOperationResponse_COMPLETED {
				fmt.Printf("Result: %s\n", resp.Result)
				done <- true
				return
			}
		}
	}()

	<-done
}
