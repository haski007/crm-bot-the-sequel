package resource

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository/mongodb"
	"github.com/caarlos0/env"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type CrmBotService struct {
	Bot            *tgbotapi.BotAPI
	Cfg            *config.Bot
	ChatRepository *mongodb.ChatRepository
}

func NewCrmBotService() (*CrmBotService, error) {
	var err error

	bot := &CrmBotService{}

	/*
	** ---> Bot configs
	 */
	bot.Cfg = &config.Bot{}
	if err := env.Parse(bot.Cfg); err != nil {
		logrus.Fatalf("[env Parse] Bot config err: %s", err)
	}

	/*
	** ---> mongo Collection
	 */
	bot.ChatRepository = &mongodb.ChatRepository{}
	bot.ChatRepository.InitDBConn()

	bot.Bot, err = tgbotapi.NewBotAPI(bot.Cfg.GetToken().String())
	if err != nil {
		return nil, err
	}

	bot.Bot.Debug = true

	logrus.Printf("Authorized on account %s", bot.Bot.Self.UserName)

	return bot, nil
}