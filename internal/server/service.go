package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/grpclog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lithammer/shortuuid/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"github.com/urfave/negroni"
)

// InterruptSignals are the default interrupt signals to catch
var InterruptSignals = []os.Signal{
	// syscall.SIGSTOP, // Not supporsed on Windows
	syscall.SIGINT,
	syscall.SIGTERM,
	syscall.SIGQUIT,
}

// Service represents the microservice
type Service struct {
	GRPCServer         *grpc.Server
	HTTPServer         *http.Server
	httpHandler        HTTPHandlerFunc
	errorHandler       runtime.ErrorHandlerFunc
	annotators         []AnnotatorFunc
	swagger            *SwaggerOpts
	staticDir          string
	muxOptions         []runtime.ServeMuxOption
	mux                *runtime.ServeMux
	middleware         []negroni.Handler
	routes             []Route
	streamInterceptors []grpc.StreamServerInterceptor
	unaryInterceptors  []grpc.UnaryServerInterceptor
	shutdownFunc       func()
	shutdownTimeout    time.Duration
	preShutdownDelay   time.Duration
	interruptSignals   []os.Signal
	grpcServerOptions  []grpc.ServerOption
	grpcDialOptions    []grpc.DialOption

	reverseProxyRegister ReverseProxyRegisterFunc
	grpcServiceRegister  GRPCServiceRegisterFunc
}

const (
	// the default timeout before the server shutdown abruptly
	defaultShutdownTimeout = 30 * time.Second
	// the default time waiting for running goroutines to finish their jobs before the shutdown starts
	defaultPreShutdownDelay = 1 * time.Second
)

// ReverseProxyFunc is the callback that the caller should implement to steps to reverse-proxy the HTTP/1 requests to gRPC
type ReverseProxyFunc func(ctx context.Context, mux *runtime.ServeMux, grpcHostAndPort string, opts []grpc.DialOption) error

// GRPCServiceRegisterFunc is the callback that the caller should implement to steps to reverse-proxy the HTTP/1 requests to gRPC
type GRPCServiceRegisterFunc func(server *grpc.Server)

// ReverseProxyRegisterFunc is the callback that the caller should implement to steps to reverse-proxy the HTTP/1 requests to gRPC
type ReverseProxyRegisterFunc func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error

// HTTPHandlerFunc is the http middleware handler function
type HTTPHandlerFunc func(*runtime.ServeMux) http.Handler

// AnnotatorFunc is the annotator function is for injecting meta data from http request into gRPC context
type AnnotatorFunc func(context.Context, *http.Request) metadata.MD

// DefaultHTTPHandler is the default http handler which does nothing
func DefaultHTTPHandler(mux *runtime.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	})
}

func defaultService() *Service {
	s := Service{}
	s.httpHandler = DefaultHTTPHandler
	s.errorHandler = runtime.DefaultHTTPErrorHandler
	s.shutdownFunc = func() {}
	s.shutdownTimeout = defaultShutdownTimeout
	s.preShutdownDelay = defaultPreShutdownDelay

	s.swagger = &SwaggerOpts{
		Up: false,
	}

	// default interrupt signals to catch, you can use InterruptSignal option to append more
	s.interruptSignals = InterruptSignals

	s.streamInterceptors = []grpc.StreamServerInterceptor{}
	s.unaryInterceptors = []grpc.UnaryServerInterceptor{}

	// install validator cond
	s.streamInterceptors = append(s.streamInterceptors, grpc_validator.StreamServerInterceptor())
	s.unaryInterceptors = append(s.unaryInterceptors, grpc_validator.UnaryServerInterceptor())

	// install prometheus cond
	s.streamInterceptors = append(s.streamInterceptors, grpc_prometheus.StreamServerInterceptor)
	s.unaryInterceptors = append(s.unaryInterceptors, grpc_prometheus.UnaryServerInterceptor)

	// install panic handler which will turn panics into gRPC errors
	s.streamInterceptors = append(s.streamInterceptors, grpc_recovery.StreamServerInterceptor())
	s.unaryInterceptors = append(s.unaryInterceptors, grpc_recovery.UnaryServerInterceptor())

	// add /metrics HTTP/1 endpoint
	routeMetrics := Route{
		Method: "GET",
		Path:   "/metrics",
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			promhttp.Handler().ServeHTTP(w, r)
		},
	}
	s.routes = append(s.routes, routeMetrics)

	return &s
}

// NewService creates a new microservice
func NewService(opts ...Option) *Service {
	s := defaultService()

	s.apply(opts...)

	// default dial option is using insecure connection
	if len(s.grpcDialOptions) == 0 {
		s.grpcDialOptions = append(s.grpcDialOptions, grpc.WithInsecure())
	}

	// init gateway mux
	s.muxOptions = append(s.muxOptions, runtime.WithErrorHandler(s.errorHandler))

	for _, annotator := range s.annotators {
		s.muxOptions = append(s.muxOptions, runtime.WithMetadata(annotator))
	}

	s.mux = runtime.NewServeMux(s.muxOptions...)

	s.grpcServerOptions = append(s.grpcServerOptions, grpc_middleware.WithStreamServerChain(s.streamInterceptors...))
	s.grpcServerOptions = append(s.grpcServerOptions, grpc_middleware.WithUnaryServerChain(s.unaryInterceptors...))

	s.GRPCServer = grpc.NewServer(
		s.grpcServerOptions...,
	)

	if s.HTTPServer == nil {
		s.HTTPServer = &http.Server{}
	}

	return s
}

// Getpid gets the process id of server
func (s *Service) Getpid() int {
	return os.Getpid()
}

// Start starts the microservice with listening on the ports
func (s *Service) Start(httpPort uint, grpcPort uint) error {

	// intercept interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, s.interruptSignals...)

	// channels to receive error
	errChan1 := make(chan error, 1)
	errChan2 := make(chan error, 1)

	// start gRPC server
	go func() {
		fmt.Printf("Starting gPRC server listening on %d\n", grpcPort)
		errChan1 <- s.startGRPCServer(grpcPort)
	}()

	// start HTTP/1.0 gateway server
	go func() {
		fmt.Printf("Starting http server listening on %d\n", httpPort)
		errChan2 <- s.startGRPCGateway(httpPort, grpcPort)
	}()

	// wait for context cancellation or shutdown signal
	select {
	// if gRPC server fail to start
	case err := <-errChan1:
		return err

	// if http server fail to start
	case err := <-errChan2:
		return err

	// if we received an interrupt signal
	case sig := <-sigChan:
		fmt.Printf("Interrupt signal received: %v\n", sig)
		s.Stop()
		return nil
	}
}

func (s *Service) startGRPCServer(grpcPort uint) error {

	s.grpcServiceRegister(s.GRPCServer)

	// register reflection service on gRPC server.
	reflection.Register(s.GRPCServer)

	grpcHost := fmt.Sprintf(":%d", grpcPort)
	lis, err := net.Listen("tcp", grpcHost)
	if err != nil {
		return err
	}

	return s.GRPCServer.Serve(lis)
}

func (s *Service) startGRPCGateway(httpPort uint, grpcPort uint) error {
	if s.swagger.Up {
		// s.swagger.EnsureDefaults()

		routeDocs := Route{
			Method: "GET",
			Path:   s.swagger.Path,
			Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
				s.swagger.ServeHTTP(w, r, pathParams)
			},
		}

		serveFile := getOpenAPIHandler()
		routeSwaggerFile := Route{
			Method: "GET",
			Path:   s.swagger.SpecURL,
			Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
				serveFile.ServeHTTP(w, r)
			},
		}
		s.AddRoutes(routeDocs, routeSwaggerFile)

	}

	endpoint := fmt.Sprintf(":%d", grpcPort)

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
	}()

	err = s.reverseProxyRegister(context.Background(), s.mux, conn)
	if err != nil {
		return err
	}

	// apply routes
	for _, route := range s.routes {
		err := s.mux.HandlePath(route.Method, route.Path, route.Handler)
		if err != nil {
			return err
		}
	}

	// create common middleware to be shared across routes
	common := negroni.New(
		s.middleware...,
	)

	s.HTTPServer.Addr = fmt.Sprintf(":%d", httpPort)
	s.HTTPServer.Handler = common.With(
		negroni.Wrap(s.mux),
	)

	s.HTTPServer.RegisterOnShutdown(s.shutdownFunc)

	return s.HTTPServer.ListenAndServe()
}

// Stop stops the microservice gracefully
func (s *Service) Stop() {
	// disable keep-alives on existing connections
	s.HTTPServer.SetKeepAlivesEnabled(false)

	// we wait for a duration of preShutdownDelay for running goroutines to finish their jobs
	if s.preShutdownDelay > 0 {
		fmt.Printf("Waiting for %v before shutdown starts\n", s.preShutdownDelay)
		time.Sleep(s.preShutdownDelay)
	}

	// gracefully stop gRPC server first
	s.GRPCServer.GracefulStop()

	var ctx, cancel = context.WithTimeout(
		context.Background(),
		s.shutdownTimeout,
	)
	defer cancel()

	// gracefully stop http server
	err := s.HTTPServer.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
}

// AddRoutes adds additional routes
func (s *Service) AddRoutes(routes ...Route) {
	s.routes = append(s.routes, routes...)
}

// HasRoute checks if a route already exists
func (s *Service) HasRoute(route Route) bool {
	for _, r := range s.routes {
		if r.Method == route.Method && r.Path == route.Path {
			return true
		}
	}

	return false
}
