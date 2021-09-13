package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/swaggo/http-swagger"
	"github.com/v8platform/ras-grpc-gw/pkg/gen/access/service"
	"github.com/v8platform/ras-grpc-gw/swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// getOpenAPIHandler serves an OpenAPI UI.
// Adapted from https://github.com/philips/grpc-gateway-example/blob/a269bcb5931ca92be0ceae6130ac27ae89582ecc/cmd/serve.go#L63
func getOpenAPIHandler() http.Handler {

	// mime.AddExtensionType(".svg", "image/svg+xml")
	// Use subdirectory in embedded files
	subFS, err := fs.Sub(swagger.SwaggerDocs, ".")
	if err != nil {
		panic("couldn't create sub filesystem: " + err.Error())
	}
	return http.FileServer(http.FS(subFS))
}

// Run runs the gRPC-Gateway, dialling the provided address.
func Run(dialAddr string) error {
	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	// Create a client connection to the gRPC Server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	conn, err := grpc.DialContext(
		context.Background(),
		dialAddr,
		grpc.WithInsecure(),
		// grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(insecure.CertPool, "")),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}

	gwmux := runtime.NewServeMux()
	err = service.RegisterAuthServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	// redocHandler.EnsureDefaults()

	// redocHandler.AddSpec("Access service", "https://petstore.swagger.io/v2/swagger.json")

	apiHandler := getOpenAPIHandler()

	port := os.Getenv("PORT")
	if port == "" {
		port = "11000"
	}
	gatewayAddr := "0.0.0.0:" + port

	docs := httpSwagger.Handler(httpSwagger.URL("/swagger/docs/docs.swagger.json"))

	gwServer := &http.Server{
		Addr: gatewayAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				gwmux.ServeHTTP(w, r)
				return
			}
			if strings.HasPrefix(r.URL.Path, "/swagger") {
				apiHandler.ServeHTTP(w, r)
				return
			}

			if strings.HasPrefix(r.URL.Path, "/docs") {
				docs.ServeHTTP(w, r)
				return
			}
		}),
	}
	// Empty parameters mean use the TLS Config specified with the server.
	// if strings.ToLower(os.Getenv("SERVE_HTTP")) == "true" {
	log.Info("Serving gRPC-Gateway and OpenAPI Documentation on http://", gatewayAddr)
	return fmt.Errorf("serving gRPC-Gateway server: %w", gwServer.ListenAndServe())
	// }

	// gwServer.TLSConfig = &tls.Config{
	// 	Certificates: []tls.Certificate{insecure.Cert},
	// }
	// log.Info("Serving gRPC-Gateway and OpenAPI Documentation on https://", gatewayAddr)
	// return fmt.Errorf("serving gRPC-Gateway server: %w", gwServer.ListenAndServeTLS("", ""))
}
