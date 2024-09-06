package servers

import (
	"context"
	"fmt"
	"github.com/t34-dev/go-svc-starter/internal/api/random"
	"github.com/t34-dev/go-svc-starter/internal/interceptor"
	"github.com/t34-dev/go-svc-starter/pkg/api/random_v1"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func GrpcServe(ctx context.Context) error {
	grpcServer := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)
	reflection.Register(grpcServer)

	random_v1.RegisterRandomServiceServer(grpcServer, random.NewImplementedRandom())

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
	log.Println("Shutting down Grpc ImplementedRandom...")
	return nil
}
