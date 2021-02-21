package resource

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callTransactionsGetAllHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	var transactions []model.Transaction

	if err := bot.TransactionRepository.GetAll(&transactions); err != nil {
		if errors.Is(err, repository.ErrDocDoesNotExist) {
			bot.Errorf(chatID, "Пока нет проведённых транзакций в базе")
			return
		}
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		bot.ReportToTheCreator(
			fmt.Sprintf("[callTransactionsGetAllHandler] TransactionRepository.GetAll | err: %s", err))
		return
	}

	var message string

	for i, tx := range transactions {
		message += fmt.Sprintf("Транзакция #%d\n"+
			"Тип транзакции: %s\n"+
			"Сумма: *%.2f UAH*\n"+
			"Автор: %s\n"+
			"Комментарий: *%s*\n"+
			"Дата: %s\n"+
			"/remove\\_tx\\_%s\n"+
			"----------------------\n",
			i+1,
			tx.Type,
			tx.Amount,
			tx.Author,
			tx.Comment,
			tx.CreatedAt.String(),
			strings.ReplaceAll(tx.ID, "-", "\\_"),
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
