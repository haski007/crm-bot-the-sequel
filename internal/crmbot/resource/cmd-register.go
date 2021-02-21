package resource

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/google/uuid"

	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) commandRegisterHandler(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	args := update.Message.CommandArguments()
	if len(args) == 0 {
		bot.Errorf(chatID, "Нужен пароль для регистрации!")
		return
	}

	cfg := config.Auth{}
	if err := env.Parse(&cfg); err != nil {
		logrus.Fatalf("[env Parse] Bot config err: %s", err)
	}

	password := strings.Fields(args)[0]

	var role string
	switch password {
	case cfg.AdminPassword:
		role = config.RoleAdmin

	default:
		bot.Errorf(chatID, "Неправильный пароль!")
		return
	}

	if err := bot.UserRepository.AddUser(model.User{
		ID:        uuid.New().String(),
		TgID:      update.Message.From.ID,
		FirstName: update.Message.From.FirstName,
		LastName:  update.Message.From.LastName,
		Username:  "@" + update.Message.From.UserName,
		Role:      role,
	}); err != nil {
		if errors.Is(err, repository.ErrDocAlreadyExists) {
			bot.Errorf(chatID, "Вы уже зарегистрированы")
			return
		}
		bot.ReportToTheCreator(fmt.Sprintf("[commandRegisterHandler] UserRepository.AddUser | err: %s", err))
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		return
	}
	bot.Reply(chatID, fmt.Sprintf("You are succesfully registered as *%s* %s", role, emoji.Check))
}
