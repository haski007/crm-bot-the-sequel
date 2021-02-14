package resource

import (
	"fmt"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callStockHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	var totalQuantity int
	var totalPurPrice model.Money
	var totalBidPrice model.Money
	var lastRevisionDate string

	// ---> Total quantity
	if res, err := bot.ProductRepository.GetFieldSum("quantity"); err != nil {
		bot.ReportToTheCreator(fmt.Sprintf("[callStockHandler] GetFieldSum | err: %s", err))
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		return
	} else {
		totalQuantity = int(res)
	}

	// ---> Total purchase price
	if res, err := bot.ProductRepository.GetFieldSum("pur_price"); err != nil {
		bot.ReportToTheCreator(fmt.Sprintf("[callStockHandler] GetFieldSum | err: %s", err))
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		return
	} else {
		totalPurPrice = model.NewMoney(res)
	}

	// ---> Total bid price
	if res, err := bot.ProductRepository.GetFieldSum("bid_price"); err != nil {
		bot.ReportToTheCreator(fmt.Sprintf("[callStockHandler] GetFieldSum | err: %s", err))
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		return
	} else {
		totalBidPrice = model.NewMoney(res)
	}

	message := fmt.Sprintf("%s *Склад* %s\n"+
		"Всего продуктов на складе: *%d*\n"+
		"Цена закупки всех продуктов: *%.2f*\n"+
		"Цена продажи всех продуктов: *%.2f*\n"+
		"Последняя ревизия: *%s*\n",
		emoji.Box,
		emoji.Box,
		totalQuantity,
		totalPurPrice,
		totalBidPrice,
		lastRevisionDate,
	)

	answer := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		messageID,
		message,
		keyboards.Stock)
	answer.ParseMode = "MarkDown"
	bot.Bot.Send(answer)
}
