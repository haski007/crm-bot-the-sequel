package resource

import (
	"fmt"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) commandUserRemove(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	clientID := update.Message.From.ID

	if !bot.AuthService.IsAdmin(clientID) {
		bot.Errorf(chatID, fmt.Sprintf("%s", repository.ErrYouHaveNoRights))
		return
	}

	if len(update.Message.Text) < len(update.Message.Text)-36 {
		bot.Errorf(chatID, "Wrong type of command!")
		return
	}
	userID := strings.ReplaceAll(update.Message.Text[len(update.Message.Text)-36:], "_", "-")

	if err := bot.UserRepository.RemoveByID(userID); err != nil {
		if err == repository.ErrDocDoesNotExist {
			bot.Errorf(chatID,
				"No user with such ID: \"%s\"", userID)
			return
		}
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		bot.ReportToTheCreator(
			fmt.Sprintf("[commandUserRemove] UserRepository.RemoveByID | err: %s", err))
		return
	}

	message := "Юзер успешно удалён " + emoji.Basket
	answer := tgbotapi.NewMessage(chatID, message)
	answer.ReplyMarkup = keyboards.MainMenu
	answer.ParseMode = config.MarkdownParseMode
	bot.Bot.Send(answer)
}
