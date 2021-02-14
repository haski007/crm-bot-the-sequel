package resource

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	getProductFieldToEdit Step = iota
	getProductValueToEdit
)

func (bot *CrmBotService) commandProductEditHandler(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	if len(update.Message.Text) < 16 {
		bot.Errorf(chatID, "Wrong type of command!")
		return
	}
	productID := strings.ReplaceAll(update.Message.Text[len(update.Message.Text)-36:], "_", "-")

	OpsQueue[userID] = &Operation{
		Name: OperationType_ProductEdit,
		Step: 0,
		Data: model.ProductEdit{
			ID:    productID,
			Field: "",
		},
	}
	message := "Что нужно изменить " + emoji.QuestionMark
	answer := tgbotapi.NewMessage(chatID, message)
	answer.ReplyMarkup = keyboards.MarkupByArray([]string{
		model.ProductEditTitle.String(),
		model.ProductEditPurPrice.String(),
		model.ProductEditBidPrice.String(),
		model.ProductEditUnit.String(),
		model.ProductEditCategory.String(),
		model.ProductEditSupplier.String(),
		model.ProductEditDescription.String(),
	})
	answer.ParseMode = "MarkDown"
	bot.Bot.Send(answer)
}

func (bot *CrmBotService) hookProductEdit(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	op := OpsQueue[userID]
	switch op.Step {
	case getProductFieldToEdit:
		field, err := model.NewProductEditField(update.Message.Text)
		if err != nil {
			if err == model.ErrNoSuchProductEditField {
				bot.Reply(chatID, "Такого поля нет! "+emoji.NoEntry+"\nПопробуй ещё раз")
				return
			}
			bot.ReportToTheCreator(fmt.Sprintf("[hookProductEdit] NewProductEditField | err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			delete(OpsQueue, userID)
			return
		}

		categoryEdit := OpsQueue[userID].Data.(model.ProductEdit)
		categoryEdit.Field = field
		OpsQueue[userID].Data = categoryEdit
		OpsQueue[userID].Step++

		switch field {
		case model.ProductEditUnit:
			message := "Выберите новую единицу измерения"
			answer := tgbotapi.NewMessage(chatID, message)
			answer.ReplyMarkup = keyboards.MarkupByArray([]string{
				model.GramUnit.String(),
				model.LiterUnit.String(),
				model.PieceUnit.String(),
			})
			bot.Bot.Send(answer)
		case model.ProductEditCategory:
			answer := tgbotapi.NewMessage(chatID, "Выберите новую категорию")

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
		case model.ProductEditSupplier:
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

		default:
			var answer tgbotapi.MessageConfig
			message := "Введите новое значение: "
			answer = tgbotapi.NewMessage(chatID, message)
			answer.ReplyMarkup = tgbotapi.NewHideKeyboard(false)
			bot.Bot.Send(answer)
		}

	case getProductValueToEdit:
		categoryEdit := OpsQueue[userID].Data.(model.ProductEdit)
		var value interface{}

		switch categoryEdit.Field {
		case model.ProductEditPurPrice:
			fallthrough
		case model.ProductEditBidPrice:
			price, err := strconv.ParseFloat(update.Message.Text, 64)
			if err != nil {
				bot.Reply(chatID, "Неверный тип данных! "+emoji.NoEntry+"\n*Попробуйте ещё раз!*")
				return
			}

			value = price
		case model.ProductEditUnit:
			unit, err := model.NewUnit(update.Message.Text)
			if err != nil {
				bot.Reply(chatID, "Неверная единица измерения! "+emoji.NoEntry+"\n*Попробуй ещё раз*")
				return
			}
			value = unit.String()
		case model.ProductEditCategory:
			categoryTitle := update.Message.Text

			var category model.Category
			if err := bot.CategoryRepository.FindByTitle(categoryTitle, &category); err != nil {
				if err == repository.ErrDocDoesNotExist {
					bot.Reply(chatID, fmt.Sprintf("Категории \"%s\" не существует! %s\n*Попробуй ещё раз*\n", categoryTitle, emoji.NoEntry))
					return
				}
				bot.ReportToTheCreator(fmt.Sprintf("[hookProductEdit] GetCategoriesList | err: %s", err))
				bot.Errorf(chatID,
					"Internal Server Error | write to @pdemian to get some help")
				delete(OpsQueue, userID)
				return
			}
			value = category.ID

		case model.ProductEditSupplier:
			supplierName := update.Message.Text

			var supplier model.Supplier
			if err := bot.SupplierRepository.FindByName(supplierName, &supplier); err != nil {
				if err == repository.ErrDocDoesNotExist {
					bot.Reply(chatID, fmt.Sprintf("Поставщика \"%s\" не существует! %s\n*Попробуй ещё раз*\n", supplierName, emoji.NoEntry))
					return
				}
				bot.ReportToTheCreator(fmt.Sprintf("[GetCategoriesList] err: %s", err))
				bot.Errorf(chatID,
					"Internal Server Error | write to @pdemian to get some help")
				delete(OpsQueue, userID)
				return
			}
			value = supplier.ID

		default:
			value = update.Message.Text
		}

		if err := bot.ProductRepository.UpdateField(
			categoryEdit.ID,
			categoryEdit.Field.BsonField(),
			value,
		); err != nil {
			if err == repository.ErrDocDoesNotExist {
				bot.Errorf(chatID, "Такой категории не существует")
				return
			}
			bot.ReportToTheCreator(fmt.Sprintf("[hookProductEdit] NewProductEditField | err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			delete(OpsQueue, userID)
			return
		}

		delete(OpsQueue, userID)

		var answer tgbotapi.MessageConfig
		message := "Изменяем..." + emoji.FaceThinking
		answer = tgbotapi.NewMessage(chatID, message)
		answer.ReplyMarkup = tgbotapi.NewHideKeyboard(false)
		bot.Bot.Send(answer)

		message = "Продукт успешно изменён " + emoji.Check
		answer = tgbotapi.NewMessage(chatID, message)
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Bot.Send(answer)
	}
}
