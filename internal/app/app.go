package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/t34-dev/go-svc-starter/internal/api/random"
	"github.com/t34-dev/go-svc-starter/internal/config"
	"github.com/t34-dev/go-svc-starter/internal/interceptor"
	"github.com/t34-dev/go-svc-starter/pkg/api/random_v1"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/olezhek28/platform_common/pkg/closer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
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
		a.initServiceProvider,
		a.initGRPCServer,
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

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)

	reflection.Register(a.grpcServer)

	random_v1.RegisterRandomServiceServer(a.grpcServer, random.NewImplementedRandom())

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	grpcGatewayMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := random_v1.RegisterRandomServiceHandlerFromEndpoint(ctx, grpcGatewayMux, a.serviceProvider.GRPCConfig().Address(), opts)
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
                            url: "/swagger/random_v1.swagger.json",
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
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: combinedHandler,
	}

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	fmt.Printf("%-10s http://%s/\n", "gRPC:", a.serviceProvider.GRPCConfig().Address())

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
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

func (a *App) runHTTPServer(ctx context.Context) error {
	fmt.Printf("%-10s http://%s/swagger/\n", "Swagger:", a.serviceProvider.HTTPConfig().Address())

	go func() {
		<-ctx.Done()
		_ = a.httpServer.Shutdown(ctx)
	}()

	return a.httpServer.ListenAndServe()
}
