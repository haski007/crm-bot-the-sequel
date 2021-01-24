package resource

import (
	"fmt"
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) commandProductRemove(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	if len(update.Message.Text) < len(update.Message.Text)-36 {
		bot.Errorf(chatID, "Wrong type of command!")
		return
	}

	productID := strings.ReplaceAll(update.Message.Text[len(update.Message.Text)-36:], "_", "-")

	if err := bot.ProductRepository.RemoveByID(productID); err != nil {
		if err == repository.ErrDocDoesNotExist {
			bot.Errorf(chatID,
				"No product with such ID: \"%s\"", productID)
			return
		}
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		bot.ReportToTheCreator(
			fmt.Sprintf("[commandProductRemove] ProductRepository.RemoveByID | err: %s", err))
	}

	message := "Продукт успешно удалён " + emoji.Basket
	answer := tgbotapi.NewMessage(chatID, message)
	answer.ReplyMarkup = keyboards.MainMenu
	answer.ParseMode = "MarkDown"
	bot.Bot.Send(answer)
}
