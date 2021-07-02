package resource

import (
	"fmt"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository/mongodb"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

const (
	revisionCheckProductsStep Step = iota
	revisionCheckCashStep
)

func (bot *CrmBotService) callRevisionProcess(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	userID := update.CallbackQuery.From.ID

	delete := tgbotapi.NewDeleteMessage(chatID, messageID)
	bot.Bot.Send(delete)

	message := fmt.Sprintf("Сейчас я буду выдавать тебе продукты один за одним...\n" +
		"Твоя задача - указать настоящее количество на складе\n" +
		"У тебя всё получится " + emoji.FaceWinking)
	bot.Reply(chatID, message)

	if err := bot.startRevisionWorkflow(startNewRevisionInfo{
		chatID: chatID,
		userID: userID,
		author: fmt.Sprintf("%s %s (@%s)",
			update.CallbackQuery.From.FirstName,
			update.CallbackQuery.From.LastName,
			update.CallbackQuery.From.UserName),
	}); err != nil {
		bot.ReportToTheCreator(fmt.Sprintf("[callRevisionProcess] startRevisionWorkflow | err: %s", err))
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		return
	}
}

func (bot *CrmBotService) hookRevisionProcess(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	input := update.Message.Text

	op := OpsQueue[userID]
	revInfo := op.Data.(revisionProcessData)

	var revision model.Revision
	if err := bot.RevisionRepository.FindByID(revInfo.revisionID, &revision); err != nil {
		switch err {
		case repository.ErrDocDoesNotExist:
			bot.Errorf(chatID,
				"No such revision in database, revision_id=%s", revInfo.revisionID)
			delete(OpsQueue, userID)
			return
		default:
			bot.ReportToTheCreator(fmt.Sprintf("[hookRevisionProcess] RevisionRepository.FindByID | err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			delete(OpsQueue, userID)
			return
		}
	}

	switch op.Step {
	case revisionCheckProductsStep:
		floatValue, err := strconv.ParseFloat(input, 64)
		if err != nil {
			answer := tgbotapi.NewMessage(
				chatID,
				fmt.Sprintf("\"%s\" - это похоже на число?! %s\n*Давай ещё раз, я верю в тебя*\n",
					update.Message.Text,
					emoji.NoEntry))
			answer.ParseMode = config.MarkdownParseMode
			answer.ReplyMarkup = keyboards.RevisionProduct

			bot.Bot.Send(answer)
			return
		}

		var newQuantity float64
		if floatValue == -1 {
			newQuantity = revInfo.products[revInfo.incr].Quantity
		} else {
			newQuantity = floatValue
		}

		revInfo.checked = append(revInfo.checked, check{
			productID:   revInfo.products[revInfo.incr].ID,
			newQuantity: newQuantity,
		})
		revInfo.incr++
		message := revInfo.getProductString()
		if message == "" {
			op.Step++

			message, err = revInfo.getCashString(bot.CashRepository)
			if err != nil {
				bot.ReportToTheCreator(fmt.Sprintf("[hookRevisionProcess] getCashString | err: %s", err))
				bot.Errorf(chatID,
					"Internal Server Error | write to @pdemian to get some help")
				delete(OpsQueue, userID)
				return
			}
		}
		op.Data = revInfo
		OpsQueue[userID] = op

		answer := tgbotapi.NewMessage(chatID, message)
		answer.ParseMode = config.MarkdownParseMode
		answer.ReplyMarkup = keyboards.RevisionProduct

		bot.Bot.Send(answer)
	case revisionCheckCashStep:
		floatValue, err := strconv.ParseFloat(input, 64)
		if err != nil {
			answer := tgbotapi.NewMessage(
				chatID,
				fmt.Sprintf("\"%s\" - и как мне это перевести в деньги?! %s\n*Давай ещё раз, я верю в тебя*\n",
					update.Message.Text,
					emoji.NoEntry))
			answer.ParseMode = config.MarkdownParseMode
			answer.ReplyMarkup = keyboards.RevisionProduct

			bot.Bot.Send(answer)
			return
		}
		delete(OpsQueue, userID)

		var newCashValue = revInfo.prevCashValue
		if floatValue != -1 {
			newCashValue = model.NewMoney(floatValue)
		}

		for _, prod := range revInfo.products {
			for _, v := range revInfo.checked {
				if v.productID == prod.ID {
					revision.ProductsCost = revision.ProductsCost.
						Add(v.newQuantity * prod.BidPrice.Float64())
					if err := bot.ProductRepository.UpdateField(prod.ID, "quantity", v.newQuantity); err != nil {
						bot.ReportToTheCreator(fmt.Sprintf("[hookRevisionProcess] ProductRepository.SetQuantity | err: %s", err))
						bot.Errorf(chatID,
							"Internal Server Error | write to @pdemian to get some help")
						return
					}
				}
			}
		}

		// ---> First row
		var prevRevision model.Revision
		if err := bot.RevisionRepository.FindPreLast(&prevRevision); err != nil {
			switch err {
			case repository.ErrDocDoesNotExist:
				// no action
				bot.Reply(chatID, "Кстати, это ваша первая ревизия!\n"+
					"С дебютом "+emoji.PartyPopper)
			default:
				bot.ReportToTheCreator(fmt.Sprintf("[hookRevisionProcess] RevisionRepository.FindLast | err: %s", err))
				bot.Errorf(chatID,
					"Internal Server Error | write to @pdemian to get some help")
				return
			}
		}

		var stockTransactions []model.TransactionStock
		if err := bot.TransactionRepository.GetStockTxAfterDate(prevRevision.CreatedAt, &stockTransactions); err != nil {
			switch err {
			case repository.ErrDocDoesNotExist:
				// no action
			default:
				bot.ReportToTheCreator(fmt.Sprintf("[hookRevisionProcess] TransactionRepository.GetAllStockTx | err: %s", err))
				bot.Errorf(chatID,
					"Internal Server Error | write to @pdemian to get some help")
				return
			}
		}

		var transactions []model.Transaction
		if err := bot.TransactionRepository.GetTxAfterDate(prevRevision.CreatedAt, &transactions); err != nil {
			switch err {
			case repository.ErrDocDoesNotExist:
				// no action
			default:
				bot.ReportToTheCreator(fmt.Sprintf("[hookRevisionProcess] TransactionRepository.GetAllStockTx | err: %s", err))
				bot.Errorf(chatID,
					"Internal Server Error | write to @pdemian to get some help")
				return
			}
		}

		var allStockTxAmount model.Money
		for _, tx := range stockTransactions {
			allStockTxAmount = allStockTxAmount.Add(tx.ProductPrice.Float64() * tx.Amount)
		}

		var allTxAmount model.Money
		for _, tx := range transactions {
			if !strings.Contains(tx.Comment, "Ревизия за") {
				allTxAmount = allTxAmount.Add(tx.Amount.Float64())
			}
		}

		firstNum := prevRevision.ProductsCost + allStockTxAmount + prevRevision.Cash + allTxAmount

		firstMessage := fmt.Sprintf("Стоимость склада за прошлую ревизию: *%.2f*\n"+
			"Все транзакции со складом с %s: *%.2f UAH*\n"+
			"Касса на период последней ревизии *%.2f UAH*\n"+
			"Все кассовые транзакции с %s: *%.2f UAH\n*"+
			"%.2f + %.2f + %.2f + %.2f = %.2f UAH\n",
			prevRevision.ProductsCost,
			prevRevision.CreatedAt.Format("2006-01-02 15:04:05"), allStockTxAmount,
			prevRevision.Cash,
			prevRevision.CreatedAt.Format("2006-01-02 15:04:05"), allTxAmount,
			prevRevision.ProductsCost, allStockTxAmount, prevRevision.Cash, allTxAmount, firstNum)

		bot.Reply(chatID, firstMessage)

		// ---> Second row

		secondNum := revision.ProductsCost + revision.Cash

		secondMessage := fmt.Sprintf("Стоимость склада за эту ревизию: *%.2f*\n"+
			"Касса за эту ревизию *%.2f UAH*\n"+
			"%.2f + %.2f = %.2f",
			revision.ProductsCost,
			revision.Cash,
			revision.ProductsCost, revision.Cash, secondNum)

		bot.Reply(chatID, secondMessage)

		// ---> Third row
		thirdNum := firstNum + secondNum

		var indicator string
		if thirdNum <= 0 {
			indicator = emoji.GreenCircle
		} else {
			indicator = emoji.RedTrianle
		}

		lastMessage := fmt.Sprintf("%s %.2f + %.2f = %.2f %s", indicator,
			firstNum,
			secondNum,
			thirdNum, indicator)

		revision.Cash = newCashValue

		cashDiff := newCashValue - revInfo.prevCashValue
		if err := bot.changeCashAmount(
			cashDiff,
			revision.Author,
			fmt.Sprintf("Ревизия за %s", prevRevision.CreatedAt.Format("2006-01-02"))); err != nil {
			bot.ReportToTheCreator(fmt.Sprintf("[hookRevisionProcess] changeCashAmount | err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
		}

		revision.UpdatedAt = time.Now()
		revision.Status = model.RevisionCompleted
		if err := bot.RevisionRepository.Update(revision); err != nil {
			switch err {
			case repository.ErrDocDoesNotExist:

			default:
				bot.ReportToTheCreator(fmt.Sprintf("[hookRevisionProcess] RevisionRepository.Update | err: %s", err))
				bot.Errorf(chatID,
					"Internal Server Error | write to @pdemian to get some help")
				return
			}
		}

		answer := tgbotapi.NewMessage(chatID, lastMessage)
		answer.ParseMode = config.MarkdownParseMode
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Bot.Send(answer)
	}
}

type revisionProcessData struct {
	revisionID    string
	incr          int
	products      []*model.Product
	checked       []check
	prevCashValue model.Money
}

type check struct {
	productID   string
	newQuantity float64
}

func (d *revisionProcessData) getProductString() string {
	if d.incr > len(d.products)-1 {
		return ""
	}
	product := fmt.Sprintf("Продукт #*%d*:\n"+
		"Название: *%s*\n"+
		"Количество на складе: *%.2f*\n"+
		"Единица измерения: *%s*\n\n"+
		"Введите новое количество для этого товара или *-1* для того чтобы оставить прежнее...",
		d.incr+1,
		d.products[d.incr].Title,
		d.products[d.incr].Quantity,
		d.products[d.incr].Unit)

	return product
}

func (d *revisionProcessData) getCashString(rep *mongodb.CashRepository) (string, error) {
	var amount model.Money

	if err := rep.GetAmount(&amount); err != nil {
		return "", fmt.Errorf("GetAmount | err: %s", err)
	}

	d.prevCashValue = amount

	return fmt.Sprintf("В кассе сейчас вот столько деняг - *%.2f UAH*\n"+
			"Можешь задать новое значение или оставить как было введя *-1*", amount),
		nil
}

type startNewRevisionInfo struct {
	chatID int64
	userID int
	author string
}

func (bot CrmBotService) startRevisionWorkflow(info startNewRevisionInfo) error {
	generatedUUID := uuid.New().String()
	newRevision := model.Revision{
		ID:        generatedUUID,
		Result:    "",
		Status:    model.RevisionProcessing,
		Author:    info.author,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	if err := bot.RevisionRepository.Add(newRevision); err != nil {
		return fmt.Errorf("RevisionRepository.Add | err: %s", err)
	}

	var revInfo revisionProcessData
	revInfo.revisionID = generatedUUID

	// ---> Get all products
	if err := bot.ProductRepository.FindAll(&revInfo.products); err != nil {
		return fmt.Errorf("ProductRepository.FindAll | err: %s", err)
	}

	if len(revInfo.products) == 0 {
		return fmt.Errorf("Пока нет продуктов в базе данных ")
	}

	OpsQueue[info.userID] = &Operation{
		Name: OperationType_RevisionProcess,
		Step: revisionCheckProductsStep,
		Data: revInfo,
	}

	// ---> Print First product
	message := revInfo.getProductString()
	answer := tgbotapi.NewMessage(info.chatID, message)
	answer.ParseMode = config.MarkdownParseMode
	answer.ReplyMarkup = keyboards.RevisionProduct

	time.Sleep(time.Second / 2)
	bot.Bot.Send(answer)

	return nil
}
