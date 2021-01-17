package crmbot

import (
	"fmt"

	"github.com/sirupsen/logrus"

	//"time"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/resource"
	"github.com/Haski007/crm-bot-the-sequel/pkg/graceshut"

	//"github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func StartBot(bot *resource.CrmBotService) {
	defer func() {
		if errR := recover(); errR != nil {
			bot.ReportToTheCreator(fmt.Sprintf("[Main panic] err: %+v\n", errR))
		}
	}()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.Bot.GetUpdatesChan(u)
	if err != nil {
		logrus.Errorf("[GetUpdatedChan] err: %s", err)
		return
	}

	go bot.HandleRoutes(updates)

	graceshut.Loop()
}
