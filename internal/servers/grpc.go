package servers

import (
	"context"
	"errors"
	"fmt"
	"github.com/t34-dev/go-svc-starter/internal/interceptor"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/t34-dev/go-svc-starter/pkg/api/random_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ random_v1.RandomServiceServer = &server{}

type server struct {
	random_v1.UnimplementedRandomServiceServer
}

func (s *server) GetPing(context.Context, *random_v1.EmptyRequest) (*random_v1.PongResponse, error) {
	return &random_v1.PongResponse{
		Text: "pong",
	}, nil
}

func (s *server) GetLen(_ context.Context, req *random_v1.TxtRequest) (*random_v1.TxtResponse, error) {
	if req == nil {
		return nil, errors.New("GetLen: TxtRequest is empty")
	}

	return &random_v1.TxtResponse{
		Number: uint32(len(req.GetText())),
	}, nil
}

func (s *server) GetPerson(context.Context, *random_v1.EmptyRequest) (*random_v1.Person, error) {
	children := []string{
		"zak1", "zak2", "zak3",
	}
	return &random_v1.Person{
		Name:     "T34",
		Age:      22,
		Children: children,
		Parent: &random_v1.Parent{
			Name: "parent",
		},
	}, nil
}

func (s *server) GetCurrentTime(ctx context.Context, _ *random_v1.EmptyRequest) (*random_v1.TimeResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &random_v1.TimeResponse{Time: timestamppb.Now()}, nil
	}
}

func (s *server) GetRandomNumber(ctx context.Context, _ *random_v1.EmptyRequest) (*random_v1.NumberResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &random_v1.NumberResponse{Number: int32(rand.Intn(100))}, nil
	}
}

func (s *server) GetRandomQuote(ctx context.Context, _ *random_v1.EmptyRequest) (*random_v1.QuoteResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		quotes := []string{
			"The only way to do great work is to love what you do.",
			"Life is what happens when you're busy making other plans.",
			"The future belongs to those who believe in the beauty of their dreams.",
		}
		return &random_v1.QuoteResponse{Quote: quotes[rand.Intn(len(quotes))]}, nil
	}
}

func (s *server) PerformLongOperation(_ *random_v1.LongOperationRequest, stream random_v1.RandomService_PerformLongOperationServer) error {
	log.Println("Starting long operation")

	for i := 0; i < 3; i++ {
		select {
		case <-stream.Context().Done():
			log.Println("Client disconnected, stopping operation")
			return stream.Context().Err()
		default:
			log.Printf("Processing step %d/10\n", i+1)
			time.Sleep(time.Second)
			if err := stream.Send(&random_v1.LongOperationResponse{
				Status:   random_v1.LongOperationResponse_IN_PROGRESS,
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
	return stream.Send(&random_v1.LongOperationResponse{
		Status:  random_v1.LongOperationResponse_COMPLETED,
		Message: "Operation completed",
		Result:  "Long operation result",
	})
}

func GrpcServe(ctx context.Context) error {
	grpcServer := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)
	reflection.Register(grpcServer)

	random_v1.RegisterRandomServiceServer(grpcServer, &server{})

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

// ShutdownGrpcServe Shutdown Grpc Serve
func ShutdownGrpcServe(_ context.Context) error {
	log.Println("Shutting down Grpc server...")
	return nil
}
