package resource

import (
	"fmt"
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) commandProductShow(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	if len(update.Message.Text) < 14 {
		bot.Errorf(chatID, "Wrong type of command!")
		return
	}
	productID := strings.ReplaceAll(update.Message.Text[len(update.Message.Text)-36:], "_", "-")

	var product model.Product
	if err := bot.ProductRepository.FindByID(productID, &product); err != nil {
		if err == repository.ErrDocDoesNotExist {
			bot.Errorf(chatID,
				"No product with such ID: \"%s\"", productID)
			return
		}
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		bot.ReportToTheCreator(
			fmt.Sprintf("[commandProductShow] ProductRepository.FindByID | productID: {%s} | err: %s", productID, err))
		return
	}

	var category model.Category
	if err := bot.CategoryRepository.FindByID(product.CategoryID, &category); err != nil {
		if err == repository.ErrDocDoesNotExist {
			bot.Errorf(chatID, "Категории c ID \"%s\" не существует! %s", product.CategoryID, emoji.NoEntry)
			return
		}
		bot.ReportToTheCreator(fmt.Sprintf("[commandProductShow] CategoryRepository.FindByID err: %s", err))
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		return
	}

	var supplier model.Supplier
	if err := bot.SupplierRepository.FindByID(product.SupplierID, &supplier); err != nil {
		if err == repository.ErrDocDoesNotExist {
			bot.Errorf(chatID, "Поставщика c ID \"%s\" не существует! %s", product.SupplierID, emoji.NoEntry)
			return
		}
		bot.ReportToTheCreator(fmt.Sprintf("[commandProductShow] SupplierRepository.FindByID err: %s", err))
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		return
	}

	message := fmt.Sprintf("Продукт: *%s*\n"+
		"Цена закупки: *%.2f*\n"+
		"Цена продажи: *%.2f*\n"+
		"Количество на складе: *%d*\n"+
		"Единица измерения *%s*\n"+
		"Категория: *%s*\n"+
		"Поставщик: *%s*\n"+
		"Описание: *%s*\n"+
		"/remove\\_product\\_%s\n",
		product.Title,
		product.PurchasingPrice,
		product.BidPrice,
		product.Quantity,
		product.Unit,
		category.Title,
		supplier.Name,
		product.Description,
		strings.ReplaceAll(product.ID, "-", "\\_"))
	answer := tgbotapi.NewMessage(chatID, message)
	answer.ParseMode = "MarkDown"
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Bot.Send(answer)
}
