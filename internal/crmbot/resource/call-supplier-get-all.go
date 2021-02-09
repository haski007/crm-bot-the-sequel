package resource

import (
	"fmt"
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callSupplierGetAllHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	var suppliers []*model.Supplier

	if err := bot.SupplierRepository.FindAll(&suppliers); err != nil {
		bot.ReportToTheCreator(fmt.Sprintf("[SupplierRepository.FindAll] err: %s", err))
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		return
	}

	if len(suppliers) == 0 {
		bot.Errorf(chatID,
			"В базе данных пока нет поставщиков")
		return
	}

	var message string

	for i, supplier := range suppliers {
		message += fmt.Sprintf("Поставщик №%d\nФИО: *%s*\nНомер: *%s*\nОписание: \"%s\"\n/remove\\_supplier\\_%s\n/edit\\_supplier\\_%s\n------------------\n",
			i+1,
			supplier.Name,
			supplier.Phone,
			supplier.Description,
			strings.ReplaceAll(supplier.ID, "-", "\\_"),
			strings.ReplaceAll(supplier.ID, "-", "\\_"),
		)
	}

	answer := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		messageID,
		message,
		keyboards.MainMenu)
	answer.ParseMode = "MarkDown"
	bot.Bot.Send(answer)
}
