package main

import (
	"context"
	"fmt"
	grpcpool "github.com/t34-dev/go-grpc-pool"
	adapterservice "github.com/t34-dev/go-svc-starter/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

func main() {
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
		log.Fatalf("Failed to create connection pool: %v", err)
	}
	defer grpcPool.Close()
	conn, err := grpcPool.Get()
	//conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Free()

	c := adapterservice.NewRandomServiceClient(conn.GetConn())

	ctx := context.Background()

	// Requests
	timeResp, err := c.GetCurrentTime(ctx, &adapterservice.EmptyRequest{})
	if err != nil {
		log.Fatalf("could not get time: %v", err)
	}
	fmt.Printf("Current time: %v\n", timeResp.GetTime().AsTime())

	numberResp, err := c.GetRandomNumber(ctx, &adapterservice.EmptyRequest{})
	if err != nil {
		log.Fatalf("could not get random number: %v", err)
	}
	fmt.Printf("Random number: %d\n", numberResp.GetNumber())

	quoteResp, err := c.GetRandomQuote(ctx, &adapterservice.EmptyRequest{})
	if err != nil {
		log.Fatalf("could not get random quote: %v", err)
	}
	fmt.Printf("Random quote: %s\n", quoteResp.GetQuote())

	txt := "test text"
	number, err := c.GetLen(ctx, &adapterservice.TxtRequest{
		Text: txt,
	})
	if err != nil {
		log.Fatalf("could not get GetLen: %v", err)
	}
	fmt.Printf("GetLen [%s] - %d\n", txt, number.GetNumber())

	stream, err := c.PerformLongOperation(ctx, &adapterservice.LongOperationRequest{})
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
			if resp.Status == adapterservice.LongOperationResponse_COMPLETED {
				fmt.Printf("Result: %s\n", resp.Result)
				done <- true
				return
			}
		}
	}()

	<-done
}
