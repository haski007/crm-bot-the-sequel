package resource

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callCategorySettingsHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID

}
