package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/urfave/cli/v2"
	"github.com/v8platform/ras-grpc-gw/internal/config"
	appCtx "github.com/v8platform/ras-grpc-gw/internal/context"
	grpc_v1 "github.com/v8platform/ras-grpc-gw/internal/delivery/grpc/v1"
	"github.com/v8platform/ras-grpc-gw/internal/repository"
	"github.com/v8platform/ras-grpc-gw/internal/server"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	"github.com/v8platform/ras-grpc-gw/pkg/auth"
	"github.com/v8platform/ras-grpc-gw/pkg/cache"
	"github.com/v8platform/ras-grpc-gw/pkg/docs/rapidoc"
	"github.com/v8platform/ras-grpc-gw/pkg/hash"
	client2 "github.com/v8platform/ras-grpc-gw/pkg/ras_client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net/http"
	"os"
	"time"
)

//go:generate buf generate --template ../buf.gen.yaml ../rasgrpcgwapis --output ../
// nolint: gochecknoglobals
var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

func main() {

	app := &cli.App{
		Name:    "ras-grpc-gw",
		Version: version,
		Authors: []*cli.Author{
			{
				Name: "Aleksey Khorev",
			},
		},
		UsageText:   "ras-grpc-wg [OPTIONS] [HOST:PORT]",
		Copyright:   "(c) 2021 Khorevaa",
		Description: "GRPC gateway for RAS 1S.Enterprise",
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:  "grpc-port",
				Value: 3000,
				Usage: "port to bind grpc server",
			},
			&cli.UintFlag{
				Name:  "http-port",
				Value: 3001,
				Usage: "port to bind http server",
			},
		},
		Action: func(c *cli.Context) error {

			var (
				repositories *repository.Repositories
				cacheEngine  cache.Cache
				tokenManager auth.TokenManager
				err          error
			)

			var cfg config.Config

			if c, err := config.NewConfigFrom(config.DefaultConfig); err != nil {
				return err
			} else if err := c.Unpack(&cfg); err != nil {
				return err
			}

			if repositories, err = repository.CreateRepository(cfg.Database); err != nil {
				return err
			}
			if cacheEngine, err = cache.New(cfg.Cache); err != nil {
				return err
			}

			if tokenManager, err = auth.NewTokenManager(cfg.Secret); err != nil {
				return err
			}

			services, err := service.NewServices(service.Options{
				Repositories: repositories,
				Cache:        cacheEngine,
				Hasher:       hash.NewSHA1Hasher(cfg.SHA1Salt),
				TokenManager: tokenManager,
			})

			if err != nil {
				return err
			}

			client := client2.NewClient(
				os.Getenv("RAS_HOST"),
				client2.EndpointUUID(func(ctx context.Context) (uuid.UUID, bool) {
					uuidStr, ok := appCtx.EndpointFromContext(ctx)
					if !ok {
						return uuid.Nil, false
					}
					return uuid.MustParse(uuidStr), true
				}),
				client2.AnnotateContext(
					client2.AnnotateRequestMetadataGrpc()),
				client2.RequestTimeout(60*time.Second),
				client2.RequestInterceptor(
					client2.OverwriteClusterIdFromContextMetadata(),
					client2.AddClusterAuthFromContextMetadata(),
					client2.AddInfobaseAuthFromContextMetadata(),
					grpc_v1.SendEndpointID,
				),
			)

			gRPCServiceRegisterFunc, reverseProxyRegisterFunc := grpc_v1.RegisterServerServices(services, client)

			interceptors := grpc_v1.NewInterceptors(services)

			svr := server.NewService(
				server.HTTPHandler(func(mux *runtime.ServeMux) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						token, err := auth.BearerExtractor(r)
						if err != nil {
							http.Error(w, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err).Error(), http.StatusUnauthorized)
							return
						}

						tokenInfo, err := services.TokenManager.Validate(token, "access")
						if err != nil {
							http.Error(w, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err).Error(), http.StatusUnauthorized)
							return
						}

						ctx := r.Context()

						if len(tokenInfo) > 0 {

							// user, err := services.Users.GetByUUID(ctx, tokenInfo)
							// if err != nil {
							// 	http.Error(w, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err).Error(), http.StatusUnauthorized)
							// 	return
							// }
							// ctx = appCtx.UserToContext(ctx, user)

						}
						mux.ServeHTTP(w, r.WithContext(ctx))
					})
				}),
				server.UnaryInterceptor(interceptors...),
				server.GRPCServiceRegister(gRPCServiceRegisterFunc),
				server.ReverseProxyRegister(reverseProxyRegisterFunc),
				server.MuxOption(runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
					switch key {
					case "X-Endpoint":
						return key, true
					default:
						return runtime.DefaultHeaderMatcher(key)
					}
				})),
				server.MuxOption(runtime.WithOutgoingHeaderMatcher(func(key string) (string, bool) {
					switch key {
					case "x-endpoint":
						return "X-Endpoint", true
					default:
						return runtime.DefaultHeaderMatcher(key)
					}
				})),
				server.MuxOption(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
					MarshalOptions: protojson.MarshalOptions{
						UseProtoNames:   true,
						EmitUnpopulated: true,
					},
					UnmarshalOptions: protojson.UnmarshalOptions{
						DiscardUnknown: true,
					},
				})),
				server.Swagger(&server.SwaggerOpts{
					Up:      true,
					Path:    "/docs/**",
					SpecURL: "/docs.swagger.json",
					Handler: rapidoc.New(
						rapidoc.URL("/docs.swagger.json"),
						rapidoc.Style(rapidoc.RenderStyle_View),
						rapidoc.Layout(rapidoc.Layout_Row),
					),
				}),
				// server.MuxOption(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONBuiltin{})),
			)

			if err := svr.Start(c.Uint("http-port"), c.Uint("grpc-port")); err != nil {
				return err
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
