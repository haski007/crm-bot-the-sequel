package resource

import (
	"fmt"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func (bot *CrmBotService) Reply(chatID int64, message string) {
	resp := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Bot.Send(resp)
	if err != nil {
		logrus.Printf("[send message /help] err: %s", err)
		return
	}
}

func (bot *CrmBotService) Errorf(chatID int64, format string, data ...interface{}) {
	message := fmt.Sprintf(emoji.Failed+" "+format+" "+emoji.Failed, data...)
	answer := tgbotapi.NewMessage(chatID, message)
	answer.ReplyMarkup = keyboards.MainMenu
	answer.ParseMode = "MarkDown"
	bot.Bot.Send(answer)
}
