package main

import (
	"fmt"
	"os"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot"
	"github.com/Haski007/crm-bot-the-sequel/pkg/run"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var Version string

func main() {
	app := cli.App{
		Name:    "crm-bot-the-sequel",
		Usage:   "Crm bot v2.0",
		Version: Version,
		Action: func(c *cli.Context) error {
			if err := crmbot.Run(&run.Args{
				LogLevel: run.LogLevel(c.String("info")),
			}); err != nil {
				return fmt.Errorf("run: %w", err)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
