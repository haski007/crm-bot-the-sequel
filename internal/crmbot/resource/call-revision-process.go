package resource

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/kr/pretty"
	"github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callRevisionProcess(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	message := pretty.Sprint(bot.RevisionRepository.RemoveByID("zadadwd"))

	logrus.Fatalln(message)
	answer := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		messageID,
		message,
		keyboards.Revision)
	answer.ParseMode = config.MarkdownParseMode
	bot.Bot.Send(answer)
}
