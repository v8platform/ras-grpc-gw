package main

import (
	"github.com/urfave/cli/v2"
	ras "github.com/v8platform/ras-grpc-gw/internal/server"
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

			host := "localhost:1545"
			if c.Args().Present() {
				host = c.Args().First()
			}

			server := ras.NewRASServer(host)

			if err := server.Serve(c.String("bind")); err != nil {
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
