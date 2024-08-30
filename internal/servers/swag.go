package servers

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	adapterservice "github.com/t34-dev/go-svc-starter/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"strings"
)

func SwaggerServe(ctx context.Context) error {
	// GRPC
	grpcGatewayMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := adapterservice.RegisterRandomServiceHandlerFromEndpoint(ctx, grpcGatewayMux, grpcAddress, opts)
	if err != nil {
		return err
	}

	// SWAGGER
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/swagger/", serveSwagger)
	httpMux.HandleFunc("/openapi/", serveOpenAPI)

	combinedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/swagger/") || strings.HasPrefix(r.URL.Path, "/openapi/") {
			httpMux.ServeHTTP(w, r)
		} else {
			grpcGatewayMux.ServeHTTP(w, r)
		}
	})

	fmt.Printf("%-10s http://%s/swagger/\n", "Swagger:", httpAddress)
	fmt.Printf("%-10s http://%s/openapi/\n", "Openapi:", httpAddress)

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

func ShutdownSwaggerServe(ctx context.Context) error {
	log.Println("Shutting down Swagger server...")
	return nil
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".json") {
		http.ServeFile(w, r, "pkg/api/v1/adapter_service.swagger.json")
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
                            url: "/swagger/adapterservice.swagger.json",
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

func serveOpenAPI(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".yaml") {
		http.ServeFile(w, r, "pkg/api/v1/openapi.yaml")
	} else {
		html := `
            <!DOCTYPE html>
            <html lang="en">
            <head>
                <meta charset="UTF-8">
                <title>OpenApi UI</title>
                <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.1.0/swagger-ui.css" >
                <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.1.0/swagger-ui-bundle.js"> </script>
            </head>
            <body>
                <div id="swagger-ui"></div>
                <script>
                    window.onload = function() {
                        SwaggerUIBundle({
                            url: "/openapi/openapi.yaml",
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
