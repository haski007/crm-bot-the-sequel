package resource

import (
	"fmt"
	"strconv"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
)

const (
	getProductTitle Step = iota
	getProductDescription
	getProductPurchasingPrice
	getProductBidPrice
	getProductCategory
	getProductSupplier
	getProductUnit
)

func (bot *CrmBotService) callProductAddHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	userID := update.CallbackQuery.From.ID

	// ---> adding to queue
	OpsQueue[userID] = &Operation{
		Name: OperationType_ProductAdd,
		Step: 0,
	}

	delete := tgbotapi.NewDeleteMessage(chatID, messageID)
	bot.Bot.Send(delete)

	message := "Введите название продукта:"
	bot.Reply(chatID, message)
}

func (bot *CrmBotService) hookProductAdd(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	op := OpsQueue[userID]

	switch op.Step {
	case getProductTitle:
		OpsQueue[userID].Data = model.Product{
			ID:              uuid.New().String(),
			Title:           update.Message.Text,
			Description:     "",
			PurchasingPrice: 0,
			BidPrice:        0,
			Quantity:        0,
			CategoryID:      "",
			SupplierID:      "",
			Unit:            "",
		}
		OpsQueue[userID].Step++
		bot.Reply(chatID, "Введите описание для продукта:")
	case getProductDescription:
		product := OpsQueue[userID].Data.(model.Product)
		product.Description = update.Message.Text
		OpsQueue[userID].Data = product

		OpsQueue[userID].Step++
		bot.Reply(chatID, "Введите цену закупки:")
	case getProductPurchasingPrice:
		purchasingPrice, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			bot.Reply(chatID, "Неверный тип данных! "+emoji.NoEntry+"\n*Попробуйте ещё раз!*")
			return
		}

		product := OpsQueue[userID].Data.(model.Product)
		product.PurchasingPrice = model.NewMoney(purchasingPrice)
		OpsQueue[userID].Data = product

		OpsQueue[userID].Step++
		bot.Reply(chatID, "Введите цену продажи:")
	case getProductBidPrice:
		bitPrice, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			bot.Reply(chatID, "Неверный тип данных! "+emoji.NoEntry+"\n*Попробуйте ещё раз!*")
			return
		}

		product := OpsQueue[userID].Data.(model.Product)
		product.BidPrice = model.NewMoney(bitPrice)
		OpsQueue[userID].Data = product

		OpsQueue[userID].Step++

		message := "Выберите категорию продукта:"
		answer := tgbotapi.NewMessage(chatID, message)

		var categories []string

		if err := bot.CategoryRepository.DistinctCategories(&categories); err != nil {
			bot.ReportToTheCreator(fmt.Sprintf("[hookProductAdd] DistinctCategories | err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			delete(OpsQueue, userID)
			return
		}

		if len(categories) == 0 {
			message := "Oops!"
			answer := tgbotapi.NewMessage(chatID, message)
			answer.ReplyMarkup = tgbotapi.NewHideKeyboard(false)
			bot.Bot.Send(answer)

			bot.Errorf(chatID, "Нет поставщиков!")
			delete(OpsQueue, userID)
			return
		}
		answer.ReplyMarkup = keyboards.MarkupByArray(categories)
		bot.Bot.Send(answer)
	case getProductCategory:
		product := OpsQueue[userID].Data.(model.Product)
		categoryTitle := update.Message.Text

		var category model.Category
		if err := bot.CategoryRepository.FindByTitle(categoryTitle, &category); err != nil {
			if err == repository.ErrDocDoesNotExist {
				bot.Errorf(chatID, "Категории \"%s\" не существует! %s", categoryTitle, emoji.NoEntry)
				return
			}
			bot.ReportToTheCreator(fmt.Sprintf("[GetCategoriesList] err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			delete(OpsQueue, userID)
			return
		}

		product.CategoryID = category.ID
		OpsQueue[userID].Data = product
		OpsQueue[userID].Step++

		var suppliers []string

		if err := bot.SupplierRepository.DistinctNames(&suppliers); err != nil {
			bot.ReportToTheCreator(fmt.Sprintf("[hookProductAdd] SupplierRepository.DistinctNames | err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			delete(OpsQueue, userID)
			return
		}

		if len(suppliers) == 0 {
			message := "Oops!"
			answer := tgbotapi.NewMessage(chatID, message)
			answer.ReplyMarkup = tgbotapi.NewHideKeyboard(false)
			bot.Bot.Send(answer)

			bot.Errorf(chatID, "Нет поставщиков!")
			delete(OpsQueue, userID)
			return
		}

		message := "Выберите поставщика продукта:"
		answer := tgbotapi.NewMessage(chatID, message)
		answer.ReplyMarkup = keyboards.MarkupByArray(suppliers)
		bot.Bot.Send(answer)

	case getProductSupplier:
		supplierName := update.Message.Text

		var supplier model.Supplier
		if err := bot.SupplierRepository.FindByName(supplierName, &supplier); err != nil {
			if err == repository.ErrDocDoesNotExist {
				bot.Errorf(chatID, "Поставщика \"%s\" не существует! %s", supplierName, emoji.NoEntry)
				return
			}
			bot.ReportToTheCreator(fmt.Sprintf("[GetCategoriesList] err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			delete(OpsQueue, userID)
			return
		}

		product := OpsQueue[userID].Data.(model.Product)
		product.SupplierID = supplier.ID
		OpsQueue[userID].Data = product
		OpsQueue[userID].Step++

		message := "Выберите единицу измерения"
		answer := tgbotapi.NewMessage(chatID, message)
		answer.ReplyMarkup = keyboards.MarkupByArray([]string{
			model.GramUnit.String(),
			model.LiterUnit.String(),
			model.PieceUnit.String(),
		})
		bot.Bot.Send(answer)
	case getProductUnit:
		input := update.Message.Text
		var unit model.Unit
		switch input {
		case model.PieceUnit.String():
			unit = model.PieceUnit
		case model.LiterUnit.String():
			unit = model.LiterUnit
		case model.GramUnit.String():
			unit = model.GramUnit
		default:
			bot.Reply(chatID, "Неверная единица измерения! "+emoji.NoEntry+"\nПопробуй ещё раз")
			return
		}

		product := OpsQueue[userID].Data.(model.Product)
		product.Unit = unit

		if err := bot.ProductRepository.Add(product); err != nil {
			if err == repository.ErrDocAlreadyExists {
				bot.Errorf(chatID,
					"Продукт с таким именем {\"%s\"} уже существует",
					product.Title)
			} else {
				bot.Errorf(chatID,
					"Internal Server Error | write to @pdemian to get some help")
				bot.ReportToTheCreator(
					fmt.Sprintf("[ProductRepository.Add] category: %+v | err: %s", product, err))
			}
			delete(OpsQueue, userID)
			return
		} else {
			delete(OpsQueue, userID)

			var answer tgbotapi.MessageConfig
			message := "Продукт успешно добавлен " + emoji.Check
			answer = tgbotapi.NewMessage(chatID, message)
			answer.ReplyMarkup = tgbotapi.NewHideKeyboard(false)
			bot.Bot.Send(answer)

			answer = tgbotapi.NewMessage(chatID, emoji.House+" *Главное Меню*"+emoji.HouseWithGarden)
			answer.ReplyMarkup = keyboards.MainMenu
			answer.ParseMode = "Markdown"
			bot.Bot.Send(answer)
		}
	}
}
