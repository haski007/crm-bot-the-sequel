package resource

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callCategoryAddHandler(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	messageID := update.Message.MessageID

	tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		messageID,
		emoji.Wrench+" *Settings* "+emoji.Gear,
		keyboards.CategorySettings)
}
