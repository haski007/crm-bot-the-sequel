package resource

import (
	"fmt"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) commandSupplierRemove(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	if len(update.Message.Text) < 18 {
		bot.Errorf(chatID, "Wrong type of command!")
		return
	}
	supplierID := strings.ReplaceAll(update.Message.Text[len(update.Message.Text)-36:], "_", "-")

	if err := bot.ProductRepository.RemoveAllBySupplierID(supplierID); err != nil {
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		bot.ReportToTheCreator(
			fmt.Sprintf("[commandSupplierRemove] ProductRepository.RemoveAllBySupplierID | err: %s", err))
		return
	}

	if err := bot.SupplierRepository.RemoveByID(supplierID); err != nil {
		if err == repository.ErrDocDoesNotExist {
			bot.Errorf(chatID,
				"No supplier with such ID: \"%s\"", supplierID)
			return
		}
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		bot.ReportToTheCreator(
			fmt.Sprintf("[SupplierRepository.RemoveByID] supplierID: {%s} | err: %s", supplierID, err))
	}

	message := "Поставщик успешно удалён " + emoji.Basket
	answer := tgbotapi.NewMessage(chatID, message)
	answer.ReplyMarkup = keyboards.MainMenu
	answer.ParseMode = config.MarkdownParseMode
	bot.Bot.Send(answer)
}
