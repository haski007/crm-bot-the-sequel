package resource

import (
	"fmt"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"

	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/google/uuid"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	getCategoryTitle       step = iota
	getCategoryDescription step = iota
)

func (bot *CrmBotService) callCategoryAddHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	userID := update.CallbackQuery.From.ID

	// ---> adding to queue
	OpsQueue[userID] = &Operation{
		Name: OperationType_CategoryAdd,
		Step: 0,
	}

	delete := tgbotapi.NewDeleteMessage(chatID, messageID)
	bot.Bot.Send(delete)

	message := "Введите название категории:"
	bot.Reply(chatID, message)
}

func (bot *CrmBotService) hookCategoryAdd(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	op := OpsQueue[userID]

	switch op.Step {
	case getCategoryTitle.Int():
		OpsQueue[userID].Data = model.Category{
			ID:          uuid.New().String(),
			Title:       update.Message.Text,
			Description: "",
		}
		OpsQueue[userID].Step++
		bot.Reply(chatID, "Введите описание для категории:")
	case getCategoryDescription.Int():
		category := OpsQueue[userID].Data.(model.Category)
		category.Description = update.Message.Text

		var message string
		if err := bot.CategoryRepository.Add(category); err != nil {
			if err == repository.ErrDocAlreadyExists {
				bot.Errorf(chatID,
					"Категория с таким именем {\"%s\"} уже существует",
					category.Title)
			} else {
				bot.Errorf(chatID,
					"Internal Server Error | write to @pdemian to get some help")
				bot.ReportToTheCreator(
					fmt.Sprintf("[CategoryRepository.Add] category: %+v | err: %s", category, err))
			}
			delete(OpsQueue, userID)
			return
		} else {
			message = "Категория успешно добавлена " + emoji.Check
			answer := tgbotapi.NewMessage(chatID, message)
			answer.ReplyMarkup = keyboards.MainMenu
			bot.Bot.Send(answer)
		}

	}
}
