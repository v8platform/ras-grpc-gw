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
			&cli.StringFlag{
				Name:  "bind",
				Value: ":3002",
				Usage: "host:port to bind grpc server",
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

			handlers := grpc_v1.NewHandlers(services)
			interceptors := grpc_v1.NewInterceptors(services)
			rasHandlers := grpc_v1.NewRasHandlers(services)

			svr := server.NewServer(
				server.WithHandlers(handlers...),
				server.WithHandlers(rasHandlers...),
				server.WithChainInterceptor(interceptors...),
			)

			if err := svr.Serve(c.String("bind")); err != nil {
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
