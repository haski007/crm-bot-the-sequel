package resource

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callProductSettingsHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	answer := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		messageID,
		emoji.Wrench+" *Настройка продуктов* "+emoji.Gear,
		keyboards.ProductSettings)
	answer.ParseMode = config.MarkdownParseMode
	bot.Bot.Send(answer)
}
