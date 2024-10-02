package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/t34-dev/go-svc-starter/internal/logger"
	"github.com/t34-dev/go-svc-starter/internal/tracing"
	"github.com/t34-dev/go-svc-starter/internal/validator"
	"github.com/t34-dev/go-svc-starter/pkg/api/common_v1"
	"github.com/t34-dev/go-utils/pkg/logs"
	"github.com/t34-dev/go-utils/pkg/sys/validate"
	"github.com/tidwall/gjson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"net/http"
)

const (
	grpcPort    = 50052
	serviceName = "other_service"
)

type server struct {
	common_v1.UnimplementedCommonV1Server
}

func (s *server) GetPost(ctx context.Context, req *common_v1.PostRequest) (*common_v1.PostResponse, error) {
	id := int64(req.GetId())
	ctx, finish := tracing.TraceFunc(ctx, "server.GetPost", map[string]interface{}{"id": id})
	defer finish()

	// validate
	err := validate.Validate(ctx, validator.ValidateID(id))
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch post: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// parse JSON using gjson
	result := gjson.ParseBytes(body)

	return &common_v1.PostResponse{
		UserId: result.Get("userId").Int(),
		Id:     result.Get("id").Int(),
		Title:  result.Get("title").String(),
		Body:   result.Get("body").String(),
	}, nil
}

func main() {
	if err := logger.SetLogLevel("info"); err != nil {
		log.Println(err)
	}
	logs.Init(logger.GetCore(logger.GetAtomicLevel()))
	tracing.Init(logs.Logger(), serviceName)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

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
