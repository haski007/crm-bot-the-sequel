package resource

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
)

const (
	getCashAddAmount Step = iota
	getCashAddComment
)

func (bot *CrmBotService) callCashAddHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	userID := update.CallbackQuery.From.ID

	// ---> adding to queue
	OpsQueue[userID] = &Operation{
		Name: OperationType_CashAdd,
		Step: getCashAddAmount,
	}

	delete := tgbotapi.NewDeleteMessage(chatID, messageID)
	bot.Bot.Send(delete)

	message := "Введите сумму изменения в UAH:"
	bot.Reply(chatID, message)
}

func (bot *CrmBotService) hookCashAdd(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	op := OpsQueue[userID]

	switch op.Step {
	case getCashAddAmount:
		amount, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			if err != nil {
				bot.Reply(chatID, "Неверный тип данных! "+emoji.NoEntry+"\n*Попробуйте ещё раз!*")
				return
			}
		}
		var txType model.TxType
		if amount > 0 {
			txType = model.TxAddCash
		} else {
			txType = model.TxGetCash
		}

		OpsQueue[userID].Data = model.Transaction{
			ID: uuid.New().String(),
			Author: fmt.Sprintf("%s %s (%s)",
				update.Message.From.FirstName,
				update.Message.From.LastName,
				update.Message.From.UserName),
			Amount:    model.NewMoney(amount),
			Type:      txType,
			CreatedAt: time.Now(),
			Comment:   "",
		}
		OpsQueue[userID].Step++
		bot.Reply(chatID, "Введите комментарий:")

	case getCashAddComment:
		comment := update.Message.Text

		tx := OpsQueue[userID].Data.(model.Transaction)
		delete(OpsQueue, userID)

		tx.Comment = comment
		tx.CreatedAt = time.Now()

		// ---> Save transaction
		if err := bot.TransactionRepository.Add(tx); err != nil {
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			bot.ReportToTheCreator(
				fmt.Sprintf("[hookCashAdd] TransactionRepository.Add| err: %s", err))
			return
		}

		// ---> Change main cashbox
		if err := bot.CashRepository.ChangeAmount(tx.Amount); err != nil {
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			bot.ReportToTheCreator(
				fmt.Sprintf("[hookCashAdd] CashRepository.ChangeAmount | err: %s", err))
			return
		}
		message := fmt.Sprintf("Транзакция успешна %s\n"+
			"tx\\_id:\n*%s*\n",
			emoji.Check,
			tx.ID)
		answer := tgbotapi.NewMessage(chatID, message)
		answer.ReplyMarkup = keyboards.MainMenu
		answer.ParseMode = "Markdown"
		bot.Bot.Send(answer)
	}
}
