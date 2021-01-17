package resource

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) commandTestHandler(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	message := fmt.Sprintf("Language = %s", update.Message.From.LanguageCode)

	bot.Reply(chatID, message)
}
