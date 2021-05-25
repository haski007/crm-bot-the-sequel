package resource

import (
	"fmt"

	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callRevision(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	message := fmt.Sprintf("%s *Ревизия* %s", emoji.MagnifyingGlass, emoji.MagnifyingGlass)

	delete := tgbotapi.NewDeleteMessage(chatID, messageID)
	bot.Bot.Send(delete)

	message := "Введите название продукта:"
	bot.Reply(chatID, message)
}
