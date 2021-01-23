package resource

import (
	"fmt"
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callCategoryGetAllHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	var categories []*model.Category

	if err := bot.CategoryRepository.FindAll(&categories); err != nil {
		bot.ReportToTheCreator(fmt.Sprintf("[CategoryRepository.FindAll] err: %s", err))
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		return
	}

	if len(categories) == 0 {
		bot.Errorf(chatID,
			"В базе данных пока нет категорий")
		return
	}

	var message string

	for i, category := range categories {
		message += fmt.Sprintf("Категория №%d\nНазвание: *%s*\nОписание: \"%s\"\n/remove\\_category\\_%s\n------------------\n",
			i+1,
			category.Title,
			category.Description,
			strings.ReplaceAll(category.ID, "-", "\\_"))
	}

	answer := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		messageID,
		message,
		keyboards.MainMenu)
	answer.ParseMode = "MarkDown"
	bot.Bot.Send(answer)
}
