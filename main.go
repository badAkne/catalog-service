package main

import (
	"fmt"
	"os"

	"github.com/badAkne/catalog-service/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	flag := &cli.BoolFlag{
		Name:    "no-json",
		Value:   true,
		Usage:   "Человеко-читаемый формат для логов вместо JSON",
		Aliases: []string{"nj"},
	}

	app := cli.App{
		Name:     "catalog-service",
		Version:  "1.0",
		Usage:    "catalog-service [global options] command [command options]",
		Commands: []*cli.Command{cmd.Migrate()},
		Flags:    []cli.Flag{flag},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

}
