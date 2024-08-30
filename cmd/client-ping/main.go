package main

import (
	"context"
	"fmt"
	"log"
	"time"

	adapterservice "github.com/t34-dev/go-svc-starter/pkg/api/v1"

	grpcpool "github.com/t34-dev/go-grpc-pool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Initialize logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Define a factory function to create gRPC connections
	factory := func() (*grpc.ClientConn, error) {
		return grpc.Dial("localhost:50051",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
			grpc.WithTimeout(5*time.Second))
	}

	// Create a new connection pool
	grpcPool, err := grpcpool.NewPool(factory, grpcpool.PoolOptions{
		MinConn: 2,
		MaxConn: 30,
	})
	if err != nil {
		logger.Fatal("Failed to create connection pool", zap.Error(err))
	}
	defer grpcPool.Close()

	// Slice of test strings
	testStrings := []string{
		"Hello, World!",
		"gRPC is awesome",
		"Periodic client",
		"Test string",
		"Another test",
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for i := 0; ; i++ {
		<-ticker.C // Wait for the next tick

		conn, err := grpcPool.Get()
		if err != nil {
			logger.Error("Failed to get connection from pool", zap.Error(err))
			continue
		}

		client := adapterservice.NewRandomServiceClient(conn.GetConn())

		// Select a string from the slice using the remainder of division by the slice length
		testString := testStrings[i%len(testStrings)]

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		response, err := client.GetLen(ctx, &adapterservice.TxtRequest{Text: testString})
		cancel()

		if err != nil {
			logger.Error("Error calling GetLen", zap.Error(err), zap.String("testString", testString))
			conn.Free()
			continue
		}

		fmt.Printf("String: %s, Length: %d, pool:%d\n", testString, response.GetNumber(), grpcPool.GetStats().CurrentConnections)

		conn.Free() // Return the connection to the pool
	}
}
