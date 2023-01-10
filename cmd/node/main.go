package main

import (
	"os"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-node-monitoring/config"
	"github.com/multiversx/mx-chain-node-monitoring/monitoring"
	"github.com/urfave/cli"
)

var (
	cliHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}
VERSION:
   {{.Version}}
   {{end}}
`
	log = logger.GetOrCreate("eventNotifier")

	generalConfigFile = cli.StringFlag{
		Name:  "general-config",
		Usage: "The path for the general config",
		Value: "./config/config.toml",
	}
)

func main() {
	app := cli.NewApp()
	cli.AppHelpTemplate = cliHelpTemplate
	app.Name = "MultiversX Node Monitoring"
	app.Flags = []cli.Flag{
		generalConfigFile,
	}
	app.Authors = []cli.Author{
		{
			Name:  "The MultiversX Team",
			Email: "contact@multiversx.com",
		},
	}
	app.Action = startNodeMonitoring

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startNodeMonitoring(ctx *cli.Context) error {
	log.Info("starting node monitoring tool...")

	flagsConfig, err := getFlagsConfig(ctx)
	if err != nil {
		return err
	}

	cfg, err := config.LoadConfig(flagsConfig.GeneralConfigPath)
	if err != nil {
		return err
	}
	cfg.Flags = flagsConfig

	runner, err := monitoring.NewMonitoringRunner(cfg)
	if err != nil {
		return err
	}

	err = runner.Start()
	if err != nil {
		return err
	}

	return nil
}

func getFlagsConfig(ctx *cli.Context) (*config.FlagsConfig, error) {
	flagsConfig := &config.FlagsConfig{}

	flagsConfig.GeneralConfigPath = ctx.GlobalString(generalConfigFile.Name)

	return flagsConfig, nil
}
