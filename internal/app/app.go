package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/t34-dev/go-svc-starter/internal/config"
	"github.com/t34-dev/go-svc-starter/internal/interceptor"
	"github.com/t34-dev/go-svc-starter/internal/logger"
	"github.com/t34-dev/go-svc-starter/internal/metric"
	"github.com/t34-dev/go-svc-starter/pkg/api/access_v1"
	"github.com/t34-dev/go-svc-starter/pkg/api/auth_v1"
	"github.com/t34-dev/go-svc-starter/pkg/api/common_v1"
	"github.com/t34-dev/go-utils/pkg/closer"
	"github.com/t34-dev/go-utils/pkg/logs"
	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	serviceProvider  *serviceProvider
	grpcServer       *grpc.Server
	httpServer       *http.Server
	prometheusServer *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctxMain context.Context) error {
	ctx, cancel := context.WithCancel(ctxMain)
	defer cancel()

	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	// Channels for errors and signals
	errChan := make(chan error)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println(strings.Repeat("=", 50))

	servers := []struct {
		name      string
		startFunc func(context.Context) error
	}{
		{"gRPC", a.runGRPCServer},
		{"HTTP", a.runHTTPServer},
		{"Prometheus", a.runPrometheus},
	}

	// Start all servers
	for _, server := range servers {
		go func(name string, startFn func(context.Context) error) {
			err := startFn(ctx)
			if err != nil {
				errChan <- fmt.Errorf("%s server error: %w", name, err)
			}
		}(server.name, server.startFunc)
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Println(strings.Repeat("=", 50))

	// Waiting for termination signal or error
	var err error
	select {
	case err = <-errChan:
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("Error occurred: %v", err)
		}
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
	}

	log.Println("Starting graceful shutdown...")
	cancel()

	for range servers {
		if shutdownErr := <-errChan; shutdownErr != nil {
			log.Printf("Shutdown error: %v", shutdownErr)
			if err == nil {
				err = shutdownErr
			}
		}
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.checkConfig,
		a.initLogger,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initPrometheus,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) checkConfig(_ context.Context) error {
	if config.Grpc().Host() == "" || config.Grpc().Port() == "" {
		return errors.New("Grpc Address is not set, current Address: " + config.Grpc().Address())
	}
	if config.Prometheus().Host() == "" || config.Prometheus().Port() == "" {
		return errors.New("Prometheus Address is not set, current Address: " + config.Prometheus().Address())
	}
	if config.Http().Host() == "" || config.Http().Port() == "" {
		return errors.New("Grpc Http is not set, current Address: " + config.Http().Address())
	}
	return nil
}
func (a *App) initConfig(ctx context.Context) error {
	err, resultChan := config.New(ctx, "configs/dev.yaml", ".env")
	if err != nil {
		return err
	}

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case result := <-resultChan:
				if result.Error != nil {
					logs.Error(fmt.Sprintf("error from Watch: %s", result.Error))
				} else {
					err = logger.SetLogLevel(config.App().LogLevel())
					if err != nil {
						logs.Error(fmt.Sprintf("error from SetLogLevel: %s", err))
					} else {
						marshal, _ := json.Marshal(config.GetAllConfig())
						if err != nil {
							logs.Error(fmt.Sprintf("config marshal error: %s", err))
						} else {
							logs.Warn(fmt.Sprintf("config is updated"), zap.String("newConfig", string(marshal)))
						}
					}
				}
			case <-ticker.C:
				logs.Debug(fmt.Sprintf("FROM: %s %s %s", config.App().Name(), config.Grpc().Port(), config.App().LogLevel()))
			}
		}
	}()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	if err := logger.SetLogLevel(config.App().LogLevel()); err != nil {
		return err
	}
	logs.Init(logger.GetCore(logger.GetAtomicLevel()))
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	creds, err := credentials.NewServerTLSFromFile("cert/service.pem", "cert/service.key")
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.MetricsInterceptor,
				interceptor.GrpcValidateInterceptor,
				interceptor.ErrorCodesInterceptor,
			),
		),
	}
	if config.App().IsTSL() {
		opts = append(opts, grpc.Creds(creds))
	} else {
		opts = append(opts, grpc.Creds(insecure.NewCredentials()))
	}

	a.grpcServer = grpc.NewServer(opts...)

	reflection.Register(a.grpcServer)
	common_v1.RegisterCommonV1Server(a.grpcServer, a.serviceProvider.GrpcImpl(ctx).Common)
	access_v1.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.GrpcImpl(ctx).Access)
	auth_v1.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.GrpcImpl(ctx).Auth)

	return nil
}

func (a *App) initPrometheus(ctx context.Context) error {
	err := metric.Init(ctx)
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	a.prometheusServer = &http.Server{
		Addr:    config.Prometheus().Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	grpcGatewayMux := runtime.NewServeMux()

	creds, err := credentials.NewClientTLSFromFile("cert/service.pem", "")
	if err != nil {
		return fmt.Errorf("failed to load client TLS credentials: %v", err)
	}
	opts := []grpc.DialOption{}

	if config.App().IsTSL() {
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	err = common_v1.RegisterCommonV1HandlerFromEndpoint(ctx, grpcGatewayMux, config.Grpc().Address(), opts)
	if err != nil {
		return err
	}

	// HTTP
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".json") {
			http.ServeFile(w, r, "pkg/api/api.swagger.json")
		} else {
			html := `
            <!DOCTYPE html>
            <html lang="en">
            <head>
                <meta charset="UTF-8">
                <title>Swagger UI</title>
                <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.1.0/swagger-ui.css" >
                <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.1.0/swagger-ui-bundle.js"> </script>
            </head>
            <body>
                <div id="swagger-ui"></div>
                <script>
                    window.onload = function() {
                        SwaggerUIBundle({
                            url: "/swagger/common_v1.swagger.json",
                            dom_id: '#swagger-ui',
                            presets: [
                                SwaggerUIBundle.presets.apis,
                                SwaggerUIBundle.SwaggerUIStandalonePreset
                            ],
                            layout: "BaseLayout"
                        })
                    }
                </script>
            </body>
            </html>
        `
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(html))
		}
	})

	combinedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/swagger/") {
			httpMux.ServeHTTP(w, r)
		} else {
			grpcGatewayMux.ServeHTTP(w, r)
		}
	})

	a.httpServer = &http.Server{
		Addr:    config.Http().Address(),
		Handler: combinedHandler,
	}

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	blue := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("%-20s %s\n", blue("gRPC:"), "http://"+config.Grpc().Address()+"/")

	list, err := net.Listen("tcp", config.Grpc().Address())
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		a.grpcServer.GracefulStop()
	}()

	if err = a.grpcServer.Serve(list); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return errors.New("gRPC server closed")
}
func (a *App) runPrometheus(ctx context.Context) error {
	blue := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("%-20s %s\n", blue("Prometheus:"), "http://"+config.Prometheus().Address()+"/metrics")

	go func() {
		<-ctx.Done()
		_ = a.prometheusServer.Shutdown(ctx)
	}()

	return a.prometheusServer.ListenAndServe()
}

func (a *App) runHTTPServer(ctx context.Context) error {
	blue := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("%-20s %s\n", blue("Swagger:"), "http://"+config.Http().Address()+"/swagger/")

	go func() {
		<-ctx.Done()
		_ = a.httpServer.Shutdown(ctx)
	}()

	return a.httpServer.ListenAndServe()
}
