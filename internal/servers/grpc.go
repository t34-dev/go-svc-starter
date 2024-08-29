package servers

import (
	"context"
	"errors"
	"fmt"
	adapterservice "github.com/t34-dev/go-svc-starter/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math/rand"
	"net"
	"time"
)

var _ adapterservice.RandomServiceServer = &server{}

type server struct {
	adapterservice.UnimplementedRandomServiceServer
}

func (s *server) GetLen(ctx context.Context, req *adapterservice.TxtRequest) (*adapterservice.TxtResponse, error) {
	if req == nil {
		return nil, errors.New("GetLen: TxtRequest is empty")
	}

	return &adapterservice.TxtResponse{
		Number: uint32(len(req.GetText())),
	}, nil
}

func (s *server) GetPerson(context.Context, *adapterservice.EmptyRequest) (*adapterservice.Person, error) {
	children := []string{
		"zak1", "zak2", "zak3",
	}
	return &adapterservice.Person{
		Name:     "T34",
		Age:      22,
		Children: children,
		Parent: &adapterservice.Parent{
			Name: "parent",
		},
	}, nil
}

func (s *server) GetCurrentTime(ctx context.Context, req *adapterservice.EmptyRequest) (*adapterservice.TimeResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &adapterservice.TimeResponse{Time: timestamppb.Now()}, nil
	}
}

func (s *server) GetRandomNumber(ctx context.Context, req *adapterservice.EmptyRequest) (*adapterservice.NumberResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &adapterservice.NumberResponse{Number: int32(rand.Intn(100))}, nil
	}
}

func (s *server) GetRandomQuote(ctx context.Context, req *adapterservice.EmptyRequest) (*adapterservice.QuoteResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		quotes := []string{
			"The only way to do great work is to love what you do.",
			"Life is what happens when you're busy making other plans.",
			"The future belongs to those who believe in the beauty of their dreams.",
		}
		return &adapterservice.QuoteResponse{Quote: quotes[rand.Intn(len(quotes))]}, nil
	}
}
func (s *server) PerformLongOperation(req *adapterservice.LongOperationRequest, stream adapterservice.RandomService_PerformLongOperationServer) error {
	log.Println("Starting long operation")

	for i := 0; i < 3; i++ {
		select {
		case <-stream.Context().Done():
			log.Println("Client disconnected, stopping operation")
			return stream.Context().Err()
		default:
			log.Printf("Processing step %d/10\n", i+1)
			time.Sleep(time.Second)
			if err := stream.Send(&adapterservice.LongOperationResponse{
				Status:   adapterservice.LongOperationResponse_IN_PROGRESS,
				Message:  "Processing...",
				Progress: int32((i + 1) * 10),
			}); err != nil {
				log.Printf("Error sending update to client: %v\n", err)
				return err
			}
			log.Printf("Sent progress update: %d%%\n", (i+1)*10)
		}
	}

	log.Println("Long operation completed successfully")
	return stream.Send(&adapterservice.LongOperationResponse{
		Status:  adapterservice.LongOperationResponse_COMPLETED,
		Message: "Operation completed",
		Result:  "Long operation result",
	})
}

func GrpcServe(ctx context.Context) error {
	grpcServer := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
	)
	reflection.Register(grpcServer)

	adapterservice.RegisterRandomServiceServer(grpcServer, &server{})

	list, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		return err
	}

	fmt.Printf("%-10s http://%s/\n", "gRPC:", grpcAddress)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	return grpcServer.Serve(list)
}

func ShutdownGrpcServe(ctx context.Context) error {
	log.Println("Shutting down Grpc server...")
	return nil
}
