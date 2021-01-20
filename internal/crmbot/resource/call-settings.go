package resource

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callSettingsHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID

	answer := tgbotapi.NewMessage(chatID, emoji.Wrench+" *Settings* "+emoji.Gear)
	answer.ParseMode = "MarkDown"
	answer.ReplyMarkup = keyboards.Settings
	bot.Bot.Send(answer)
}
