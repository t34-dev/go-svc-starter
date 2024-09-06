package servers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/t34-dev/go-svc-starter/pkg/api/random_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SwaggerServe(ctx context.Context) error {
	grpcGatewayMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := random_v1.RegisterRandomServiceHandlerFromEndpoint(ctx, grpcGatewayMux, grpcAddress, opts)
	if err != nil {
		return err
	}

	// SWAGGER
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/swagger/", serveSwagger)

	combinedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/swagger/") {
			httpMux.ServeHTTP(w, r)
		} else {
			grpcGatewayMux.ServeHTTP(w, r)
		}
	})

	fmt.Printf("%-10s http://%s/swagger/\n", "Swagger:", httpAddress)

	server := &http.Server{
		Addr:    ":8081",
		Handler: combinedHandler,
	}

	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}

func ShutdownSwaggerServe(_ context.Context) error {
	log.Println("Shutting down Swagger Implementation...")
	return nil
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
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
}
