package resource

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func (bot *CrmBotService) callRevisionProcess(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

}
