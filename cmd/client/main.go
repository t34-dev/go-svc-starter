package main

import (
	"context"
	"fmt"
	grpcpool "github.com/t34-dev/go-grpc-pool"
	"io"
	"log"
	"time"

	"github.com/t34-dev/go-svc-starter/pkg/api/random_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Define a factory function to create gRPC connections
	factory := func() (*grpc.ClientConn, error) {
		//opts := []grpc.DialOption{
		//	grpc.WithTransportCredentials(insecure.NewCredentials()),
		//}
		//return grpc.NewClient("localhost"+constants.Address, opts...)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		return grpc.DialContext(ctx, "localhost:50051",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock())
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

	c := random_v1.NewRandomServiceClient(conn.GetConn())

	ctx := context.Background()

	// Requests
	timeResp, err := c.GetCurrentTime(ctx, &random_v1.EmptyRequest{})
	if err != nil {
		log.Fatalf("could not get time: %v", err)
	}
	fmt.Printf("Current time: %v\n", timeResp.GetTime().AsTime())

	numberResp, err := c.GetRandomNumber(ctx, &random_v1.EmptyRequest{})
	if err != nil {
		log.Fatalf("could not get random number: %v", err)
	}
	fmt.Printf("Random number: %d\n", numberResp.GetNumber())

	quoteResp, err := c.GetRandomQuote(ctx, &random_v1.EmptyRequest{})
	if err != nil {
		log.Fatalf("could not get random quote: %v", err)
	}
	fmt.Printf("Random quote: %s\n", quoteResp.GetQuote())

	txt := "test text"
	number, err := c.GetLen(ctx, &random_v1.TxtRequest{
		Text: txt,
	})
	if err != nil {
		log.Fatalf("could not get GetLen: %v", err)
	}
	fmt.Printf("GetLen [%s] - %d\n", txt, number.GetNumber())

	stream, err := c.PerformLongOperation(ctx, &random_v1.LongOperationRequest{})
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
			if resp.Status == random_v1.LongOperationResponse_COMPLETED {
				fmt.Printf("Result: %s\n", resp.Result)
				done <- true
				return
			}
		}
	}()

	<-done
}
