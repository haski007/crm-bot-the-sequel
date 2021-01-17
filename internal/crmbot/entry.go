package crmbot

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/resource"
	"github.com/Haski007/crm-bot-the-sequel/pkg/factory"
	"github.com/Haski007/crm-bot-the-sequel/pkg/run"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func Run(args *run.Args) error {
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("[Load .env] err: %s", err)
	}

	botService, err := resource.NewCrmBotService()
	if err != nil {
		logrus.Fatalf("[NewCrmBotService] err: %s", err)
	}

	factory.InitLog(args.LogLevel)

	StartBot(botService)
	return nil
}
