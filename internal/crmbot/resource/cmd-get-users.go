package resource

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) commandGetUsersHandler(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	var users []model.User
	if err := bot.UserRepository.GetAll(&users); err != nil {
		if errors.Is(err, repository.ErrDocDoesNotExist) {
			bot.Errorf(chatID, "Нет зарегистрированных пользователей")
			return
		}
		bot.ReportToTheCreator(fmt.Sprintf("[commandGetUsersHandler] UserRepository.GetAllUsers | err: %s", err))
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		return
	}

	var message string
	for i, user := range users {
		message += fmt.Sprintf("User #%d\n"+
			"Telegram ID: %d\n"+
			"First Name: *%s*\n"+
			"Last Name: *%s*\n"+
			"Username %s\n"+
			"Role \"*%s*\"\n"+
			"/remove\\_user\\_%s\n"+
			"________________________",
			i+1,
			user.TgID,
			user.FirstName,
			user.LastName,
			user.Username,
			user.Role,
			strings.ReplaceAll(user.ID, "-", "\\_"),
		)
	}

	answer := tgbotapi.NewMessage(chatID, message)
	answer.ParseMode = "MarkDown"
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Bot.Send(answer)
}
