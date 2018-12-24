package main

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

var config Config

func main() {
	var configFile string
	app := cli.NewApp()
	app.Name = "bldr_package_sync"
	app.Usage = "CLI Application to manage the sync process from upstream habitat builders"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Load configuration from `FILE`",
			Destination: &configFile,
			Value:       "./config.toml",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "sync",
			Aliases: []string{"s"},
			Usage:   "Run the upstream sync process",
			Action: func(c *cli.Context) error {
				log.Debug("Launching the sync process with config file: " + configFile)
				if _, err := toml.DecodeFile(configFile, &config); err != nil {
					log.Error(err)
					return err
				}
				log.Info(config)
				syncer := Syncer{config: config}
				log.Info(syncer)
				return syncer.run()
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
