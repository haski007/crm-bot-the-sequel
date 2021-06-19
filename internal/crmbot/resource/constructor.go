package resource

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository/mongodb"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/service/auth"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type CrmBotService struct {
	Bot *tgbotapi.BotAPI
	Cfg *config.Bot

	AuthService *auth.AuthService

	ProductRepository     *mongodb.ProductRepository
	CategoryRepository    *mongodb.CategoryRepository
	SupplierRepository    *mongodb.SupplierRepository
	TransactionRepository *mongodb.TransactionRepository
	CashRepository        *mongodb.CashRepository
	UserRepository        *mongodb.UserRepository
	RevisionRepository    *mongodb.RevisionRepository
}

func NewCrmBotService() (*CrmBotService, error) {
	var err error

	var cfg config.Config
	if err := cfg.Parse(config.ConfigFile); err != nil {
		logrus.Fatalf("[cfg.Parse] err: %s", err)
	}

	bot := &CrmBotService{}

	/*
	** ---> Bot configs
	 */
	bot.Cfg = &cfg.Bot

	/*
	** ---> mongo Collection
	 */

	mgoSession := mongodb.InitDBConnection(cfg.MongoDB)

	bot.ProductRepository = &mongodb.ProductRepository{}
	bot.ProductRepository.InitConn(mgoSession, cfg.MongoDB.DBName)

	bot.CategoryRepository = &mongodb.CategoryRepository{}
	bot.CategoryRepository.InitConn(mgoSession, cfg.MongoDB.DBName)

	bot.SupplierRepository = &mongodb.SupplierRepository{}
	bot.SupplierRepository.InitConn(mgoSession, cfg.MongoDB.DBName)

	bot.TransactionRepository = &mongodb.TransactionRepository{}
	bot.TransactionRepository.InitConn(mgoSession, cfg.MongoDB.DBName)

	bot.CashRepository = &mongodb.CashRepository{}
	bot.CashRepository.InitConn(mgoSession, cfg.MongoDB.DBName)

	bot.UserRepository = &mongodb.UserRepository{}
	bot.UserRepository.InitConn(mgoSession, cfg.MongoDB.DBName)

	bot.RevisionRepository = &mongodb.RevisionRepository{}
	bot.RevisionRepository.InitConn(mgoSession, cfg.MongoDB.DBName)

	// ---> Init Bot
	bot.Bot, err = tgbotapi.NewBotAPI(bot.Cfg.GetToken().String())
	if err != nil {
		return nil, err
	}

	// ---> AuthService
	bot.AuthService, err = auth.NewAuthService(auth.Config{
		MgoSession:    mgoSession,
		Mongo:         cfg.MongoDB,
		UsersCollName: "users",
	})
	if err != nil {
		logrus.Fatalf("[NewAuthService] err: %s", err)
	}

	bot.Bot.Debug = true

	logrus.Printf("Authorized on account %s", bot.Bot.Self.UserName)

	return bot, nil
}
