package resource

import (
	"errors"
	"fmt"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"strconv"
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"

	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callTransactionsGetAllHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	userID := update.CallbackQuery.From.ID

	OpsQueue[userID] = &Operation{
		Name: OperationType_TransactionsGetAll,
		Step: 0,
		Data: nil,
	}

	delete := tgbotapi.NewDeleteMessage(chatID, messageID)
	bot.Bot.Send(delete)

	message := "За сколько дней нужна информация?"
	bot.Reply(chatID, message)
}

func (bot *CrmBotService) hookTransactionsGetAll(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	days, err := strconv.Atoi(update.Message.Text)
	if err != nil {
		bot.Reply(chatID, "Неверный тип данных! "+emoji.NoEntry+"\n*Попробуйте ещё раз!*")
		return
	}
	delete(OpsQueue, userID)

	var transactions []model.Transaction

	if err := bot.TransactionRepository.GetForLastDays(days, &transactions); err != nil {
		if errors.Is(err, repository.ErrDocDoesNotExist) {
			bot.Errorf(chatID, "Пока нет проведённых транзакций в базе")
			return
		}
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		bot.ReportToTheCreator(
			fmt.Sprintf("[callTransactionsGetAllHandler] TransactionRepository.GetForLastDays | err: %s", err))
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

	answer := tgbotapi.NewMessage(chatID, message)
	answer.ParseMode = config.MarkdownParseMode
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Bot.Send(answer)
}
