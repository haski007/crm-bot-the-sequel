package resource

import (
	"fmt"
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) commandTransactionRemoveHandler(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	if !bot.AuthService.IsAdmin(userID) {
		bot.Errorf(chatID, fmt.Sprintf("%s", repository.ErrYouHaveNoRights))
		return
	}

	if len(update.Message.Text) < len(update.Message.Text)-36 {
		bot.Errorf(chatID, "Wrong type of command!")
		return
	}
	txID := strings.ReplaceAll(update.Message.Text[len(update.Message.Text)-36:], "_", "-")

	var tx model.Transaction
	if err := bot.TransactionRepository.GetTxByID(txID, &tx); err != nil {
		if err == repository.ErrDocDoesNotExist {
			bot.Errorf(chatID,
				"No transaction with such ID: \"%s\"", txID)
			return
		}
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		bot.ReportToTheCreator(
			fmt.Sprintf("[commandTransactionRemoveHandler] TransactionRepository.RemoveByID | err: %s", err))
		return
	}

	if err := bot.CashRepository.ChangeAmount(-tx.Amount); err != nil {
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		bot.ReportToTheCreator(
			fmt.Sprintf("[commandTransactionRemoveHandler] CashRepository.ChangeAmount | err: %s", err))
		return
	}

	if err := bot.TransactionRepository.RemoveByID(txID); err != nil {
		if err == repository.ErrDocDoesNotExist {
			bot.Errorf(chatID,
				"No transaction with such ID: \"%s\"", txID)
			return
		}
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		bot.ReportToTheCreator(
			fmt.Sprintf("[commandTransactionRemoveHandler] TransactionRepository.RemoveByID | err: %s", err))
		return
	}

	message := "Транзакция успешно удалёна " + emoji.Basket
	answer := tgbotapi.NewMessage(chatID, message)
	answer.ReplyMarkup = keyboards.MainMenu
	answer.ParseMode = "MarkDown"
	bot.Bot.Send(answer)
}
