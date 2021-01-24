package resource

import (
	"fmt"
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callProductGetAllHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	userID := update.CallbackQuery.From.ID

	message := "Выберите категорию продукта:"

	var categories []string

	if err := bot.CategoryRepository.DistinctCategories(&categories); err != nil {
		bot.ReportToTheCreator(fmt.Sprintf("[callProductGetAllHandler] DistinctCategories | err: %s", err))
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		return
	}

	if len(categories) == 0 {
		message := "Oops!"
		answer := tgbotapi.NewMessage(chatID, message)
		answer.ReplyMarkup = tgbotapi.NewHideKeyboard(false)
		bot.Bot.Send(answer)

		bot.Errorf(chatID, "Нет поставщиков!")
		return
	}

	OpsQueue[userID] = &Operation{
		Name: OperationType_ProductGetByCategory,
		Step: 4,
		Data: nil,
	}

	answer := tgbotapi.NewMessage(chatID, message)
	answer.ReplyMarkup = keyboards.MarkupByArray(categories)
	bot.Bot.Send(answer)
}

func (bot *CrmBotService) hookProductGetByCategory(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	categoryTitle := update.Message.Text

	op := OpsQueue[userID]
	switch op.Step {
	case getProductCategory.Int():

		var category model.Category
		if err := bot.CategoryRepository.FindByTitle(categoryTitle, &category); err != nil {
			if err == repository.ErrDocDoesNotExist {
				bot.Errorf(chatID, "Категории \"%s\" не существует! %s", categoryTitle, emoji.NoEntry)
				return
			}
			bot.ReportToTheCreator(fmt.Sprintf("[hookProductGetByCategory] CategoryRepository.FindByTitle err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			delete(OpsQueue, userID)
			return
		}

		var products []*model.Product
		if err := bot.ProductRepository.FindAllByCategoryID(category.ID, &products); err != nil {
			bot.ReportToTheCreator(fmt.Sprintf("[hookProductGetByCategory] ProductRepository.FindAllByCategoryID | err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			return
		}

		if len(products) == 0 {
			bot.Errorf(chatID,
				"В базе данных пока нет продуктов в этой категрии")
			delete(OpsQueue, userID)
			return
		}
		delete(OpsQueue, userID)

		message := fmt.Sprintf(emoji.Eye+" Продукты категории: *%s*"+emoji.Eye+"\n", category.Title)

		for i, product := range products {
			message += fmt.Sprintf("Продукт №%d\nНазвание: *%s*\n/show\\_product\\_%s\n------------------\n",
				i+1,
				product.Title,
				strings.ReplaceAll(product.ID, "-", "\\_"))
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
