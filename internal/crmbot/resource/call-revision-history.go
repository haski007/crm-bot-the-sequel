package resource

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callRevisionHistory(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	var message string

	message = "revision history"
	answer := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		messageID,
		message,
		keyboards.MainMenu)
	bot.Bot.Send(answer)
}
