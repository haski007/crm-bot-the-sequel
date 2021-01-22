package resource

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callHomeHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	message := emoji.House + " *Главное меню* " + emoji.HouseWithGarden

	answer := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		messageID,
		message,
		keyboards.MainMenu)
	answer.ParseMode = "MarkDown"
	bot.Bot.Send(answer)
}
