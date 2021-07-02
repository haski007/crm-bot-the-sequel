package resource

import (
	"fmt"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callCashHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	var totalMoney model.Money
	totalMoney, err := bot.getCashAmount()
	if err != nil {
		bot.ReportToTheCreator(fmt.Sprintf("[callCashHandler] getCashAmount | err: %s", err))
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		return
	}

	var moneyInGoods model.Money
	moneyInGoods, err = bot.getMoneyInGoods()
	if err != nil {
		bot.ReportToTheCreator(fmt.Sprintf("[callCashHandler] getMoneyInGoods | err: %s", err))
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		return
	}

	var monthMoneySpent model.Money
	var monthMoneyProfit model.Money
	currMonthTxs, err := bot.getMonthTransactions()
	if err != nil {
		bot.ReportToTheCreator(fmt.Sprintf("[callCashHandler] getMonthMoneySpent | err: %s", err))
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		return
	}

	for _, tx := range currMonthTxs {
		if tx.Type == model.TxAddCash {
			monthMoneyProfit += tx.Amount
		} else if tx.Type == model.TxGetCash {
			monthMoneySpent += tx.Amount
		}
	}

	message := fmt.Sprintf("%s *Каса* %s\n"+
		"Денег в кассе: *%.2f*\n"+
		"Денег в товаре: *%.2f*\n"+
		"Траты за этот месяц: *%.2f*\n"+
		"Прибыль за этот месяц: *%.2f*\n",
		emoji.FaceMoney,
		emoji.FaceMoney,
		totalMoney,
		moneyInGoods,
		monthMoneySpent,
		monthMoneyProfit,
	)
	answer := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		messageID,
		message,
		keyboards.Cash)
	answer.ParseMode = config.MarkdownParseMode
	bot.Bot.Send(answer)
}

func (bot *CrmBotService) getCashAmount() (model.Money, error) {
	var cash model.Money
	if err := bot.CashRepository.GetAmount(&cash); err != nil {
		return model.Money(0), err
	}
	return cash, nil
}

func (bot *CrmBotService) getMoneyInGoods() (model.Money, error) {
	var products []*model.Product

	if err := bot.ProductRepository.FindAll(&products); err != nil {
		return model.Money(0), err
	}

	var result model.Money
	for _, pro := range products {
		result += model.NewMoney(pro.PurchasingPrice.Float64() * float64(pro.Quantity))
	}

	return result, nil
}

func (bot *CrmBotService) getMonthTransactions() ([]*model.Transaction, error) {
	var transactions []*model.Transaction

	if err := bot.TransactionRepository.GetCurrentMonth(&transactions); err != nil {
		return []*model.Transaction{}, err
	}

	return transactions, nil
}
