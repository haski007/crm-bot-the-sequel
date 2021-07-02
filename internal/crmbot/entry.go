package crmbot

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/resource"
	"github.com/Haski007/crm-bot-the-sequel/pkg/run"
	"github.com/sirupsen/logrus"
)

func Run(args *run.Args) error {
	botService, err := resource.NewCrmBotService()
	if err != nil {
		logrus.Fatalf("[NewCrmBotService] err: %s", err)
	}

	//factory.InitLog(args.LogLevel)

	StartBot(botService)
	return nil
}
