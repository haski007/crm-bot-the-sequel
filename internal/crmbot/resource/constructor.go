package resource

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository/mongodb"
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type CrmBotService struct {
	Bot *tgbotapi.BotAPI
	Cfg *config.Bot

	ProductRepository     *mongodb.ProductRepository
	CategoryRepository    *mongodb.CategoryRepository
	SupplierRepository    *mongodb.SupplierRepository
	TransactionRepository *mongodb.TransactionRepository
	CashRepository        *mongodb.CashRepository
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

	bot.ProductRepository = &mongodb.ProductRepository{}
	bot.ProductRepository.InitConn()

	bot.CategoryRepository = &mongodb.CategoryRepository{}
	bot.CategoryRepository.InitConn()

	bot.SupplierRepository = &mongodb.SupplierRepository{}
	bot.SupplierRepository.InitConn()

	bot.TransactionRepository = &mongodb.TransactionRepository{}
	bot.TransactionRepository.InitConn()

	bot.CashRepository = &mongodb.CashRepository{}
	bot.CashRepository.InitConn()

	bot.Bot, err = tgbotapi.NewBotAPI(bot.Cfg.GetToken().String())
	if err != nil {
		return nil, err
	}

	bot.Bot.Debug = true

	logrus.Printf("Authorized on account %s", bot.Bot.Self.UserName)

	return bot, nil
}
