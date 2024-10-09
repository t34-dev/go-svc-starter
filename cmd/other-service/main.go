package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	jsonplaceholder "github.com/t34-dev/go-svc-starter/internal/client/json_placeholder"
	jsonplaceholderImpl "github.com/t34-dev/go-svc-starter/internal/client/json_placeholder/impl"
	"github.com/t34-dev/go-svc-starter/internal/logger"
	"github.com/t34-dev/go-svc-starter/internal/validator"
	"github.com/t34-dev/go-svc-starter/pkg/api/common_v1"
	"github.com/t34-dev/go-utils/pkg/logs"
	"github.com/t34-dev/go-utils/pkg/sys/validate"
	"github.com/t34-dev/go-utils/pkg/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

const (
	grpcPort    = 50052
	serviceName = "other_service"
)

var service jsonplaceholder.JSONPlaceholderService

type server struct {
	common_v1.UnimplementedCommonV1Server
}

func (s *server) GetPost(ctx context.Context, req *common_v1.PostRequest) (*common_v1.PostResponse, error) {
	id := req.GetId()
	ctx, finish := trace.TraceFunc(ctx, "server.GetPost", map[string]interface{}{"id": id})
	defer finish()

	// validate
	err := validate.Validate(ctx, validator.ValidateID(id))
	if err != nil {
		return nil, err
	}

	post, err := service.GetPost(ctx, int(id))
	if err != nil {
		return nil, err
	}
	return &common_v1.PostResponse{
		UserId: int64(post.UserID),
		Id:     int64(post.ID),
		Title:  post.Title,
		Body:   post.Body,
	}, nil
}

func main() {
	if err := logger.SetLogLevel("info"); err != nil {
		log.Println(err)
	}
	logs.Init(logger.GetCore(logger.GetAtomicLevel(), "logs/other_service.log"))
	tracing.Init(logs.Logger(), serviceName)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// SERVICE
	service = jsonplaceholderImpl.NewService(
		jsonplaceholderImpl.WithTimeout(time.Second*10),
		jsonplaceholderImpl.WithUserAgent("MyApp/1.0"),
	)

	// SERVER
	s := grpc.NewServer(
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())),
	)
	reflection.Register(s)
	common_v1.RegisterCommonV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
