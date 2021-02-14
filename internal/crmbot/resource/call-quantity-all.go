package resource

import (
	"fmt"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callQuantityAllHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	userID := update.CallbackQuery.From.ID

	message := "Выберите категорию продукта:"

	var categories []string

	if err := bot.CategoryRepository.DistinctCategories(&categories); err != nil {
		bot.ReportToTheCreator(fmt.Sprintf("[callQuantitySet] DistinctCategories | err: %s", err))
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		return
	}

	if len(categories) == 0 {
		message := "Oops!"
		answer := tgbotapi.NewMessage(chatID, message)
		answer.ReplyMarkup = tgbotapi.NewHideKeyboard(false)
		bot.Bot.Send(answer)

		bot.Errorf(chatID, "Нет категорий!")
		return
	}

	OpsQueue[userID] = &Operation{
		Name: OperationType_QuantityAll,
		Step: getProductCategory,
		Data: nil,
	}

	answer := tgbotapi.NewMessage(chatID, message)
	answer.ReplyMarkup = keyboards.MarkupByArray(categories)
	bot.Bot.Send(answer)
}

func (bot *CrmBotService) hookQuantityAll(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	op := OpsQueue[userID]

	switch op.Step {
	case getProductCategory:
		categoryTitle := update.Message.Text

		var category model.Category
		if err := bot.CategoryRepository.FindByTitle(categoryTitle, &category); err != nil {
			if err == repository.ErrDocDoesNotExist {
				bot.Errorf(chatID, "Категории \"%s\" не существует! %s", categoryTitle, emoji.NoEntry)
				return
			}
			bot.ReportToTheCreator(fmt.Sprintf("[hookQuantitySet] CategoryRepository.FindByTitle err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			delete(OpsQueue, userID)
			return
		}

		var products []*model.Product
		if err := bot.ProductRepository.FindAllByCategoryID(category.ID, &products); err != nil {
			bot.ReportToTheCreator(fmt.Sprintf("[hookQuantitySet] ProductRepository.FindAllByCategoryID | err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			return
		}

		if len(products) == 0 {
			bot.Errorf(chatID,
				"В базе данных пока нет продуктов в этой категории")
			delete(OpsQueue, userID)
			return
		}
		delete(OpsQueue, userID)

		var message string
		for i, product := range products {
			message += fmt.Sprintf("%d) *%s* (%s) | *%d* | %.2f/%.2f\n",
				i+1,
				product.Title,
				categoryTitle,
				product.Quantity,
				product.PurchasingPrice,
				product.BidPrice,
			)
		}

		var answer tgbotapi.MessageConfig
		answer = tgbotapi.NewMessage(chatID, message)
		answer.ParseMode = "Markdown"
		answer.ReplyMarkup = tgbotapi.NewHideKeyboard(false)
		bot.Bot.Send(answer)

		answer = tgbotapi.NewMessage(chatID, emoji.House+" *Главное Меню*"+emoji.HouseWithGarden)
		answer.ReplyMarkup = keyboards.MainMenu
		answer.ParseMode = "Markdown"
		bot.Bot.Send(answer)
	}
}
