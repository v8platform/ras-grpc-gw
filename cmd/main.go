package main

import (
	"github.com/urfave/cli/v2"
	"github.com/v8platform/ras-grpc-gw/internal/database/pudgedb"
	grpc_v1 "github.com/v8platform/ras-grpc-gw/internal/delivery/grpc/v1"
	"github.com/v8platform/ras-grpc-gw/internal/repository"
	"github.com/v8platform/ras-grpc-gw/internal/server"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	"github.com/v8platform/ras-grpc-gw/pkg/auth"
	"github.com/v8platform/ras-grpc-gw/pkg/cache"
	"github.com/v8platform/ras-grpc-gw/pkg/docs/rapidoc"
	"github.com/v8platform/ras-grpc-gw/pkg/hash"
	"log"
	"os"
)

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

			// host := "localhost:1545"
			// if c.Args().Present() {
			// 	host = c.Args().First()
			// }
			manager, err := auth.NewTokenManager("sercet")
			if err != nil {
				return err
			}

			db := pudgedb.New("./pudgedb")

			services, err := service.NewServices(service.Options{
				Repositories: repository.NewPudgeRepositories(db),
				Cache:        cache.NewMemoryCache(),
				Hasher:       hash.NewSHA1Hasher("salt"),
				TokenManager: manager,
			})

			if err != nil {
				return err
			}

			gRPCServiceRegisterFunc, reverseProxyRegisterFunc := grpc_v1.RegisterServerServices(services)

			interceptors := grpc_v1.NewInterceptors(services)

			svr := server.NewService(
				server.UnaryInterceptor(interceptors...),
				server.GRPCServiceRegister(gRPCServiceRegisterFunc),
				server.ReverseProxyRegister(reverseProxyRegisterFunc),
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
