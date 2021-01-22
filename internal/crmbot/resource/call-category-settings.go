package resource

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callCategorySettingsHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	answer := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		messageID,
		emoji.Wrench+" *Настройка категорий* "+emoji.Gear,
		keyboards.CategorySettings)
	answer.ParseMode = "MarkDown"
	bot.Bot.Send(answer)
	//bot.Reply(chatID, fmt.Sprintln(messageID))
}
