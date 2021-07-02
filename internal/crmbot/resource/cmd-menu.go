package resource

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) commandMenuHandler(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	message := emoji.House + " *Главное меню* " + emoji.HouseWithGarden

	answer := tgbotapi.NewMessage(chatID, message)
	answer.ReplyMarkup = keyboards.MainMenu
	answer.ParseMode = config.MarkdownParseMode
	bot.Bot.Send(answer)
}
