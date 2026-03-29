package cmd

import (
	"strings"

	"github.com/badAkne/catalog-service/internal/app/builder"
	"github.com/urfave/cli/v2"
)

const (
	cmdWebServerUsage = "Starts the web (REST) server"

	cmdWebServerDescription = `
Initializes and starts web-server, that listens specified port
for incoming REST requests.
`
)

func WebServer() *cli.Command {

	return &cli.Command{
		Name:            "web-server",
		Aliases:         []string{"web", "http"},
		Usage:           cmdWebServerUsage,
		Description:     strings.TrimSpace(cmdMigrateDescription),
		Action:          cmdWebServer,
		HideHelpCommand: true,
	}
}

func cmdWebServer(cCtx *cli.Context) error {
	var app = builder.NewBuilder(cCtx)

	app.BuildConfig()

	app.BuildRepoConnPostgres()

	app.BuilderRepoCategory()
	app.BuilderRepoProduct()

	app.BuildServiceCategory()
	app.BuildServiceProduct()

	app.BuildHandlerHttpCategory()
	app.BuildHandlerHttpProduct()

	app.BuildProcHttp()

	app.Run()

	return nil
}
