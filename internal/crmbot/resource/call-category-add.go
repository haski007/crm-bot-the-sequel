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

func (bot *CrmBotService) callCategoryAddHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	userID := update.CallbackQuery.From.ID

	// ---> adding to queue
	OpsQueue[userID] = &Operation{
		Name: OperationType_CategoryAdd,
		Step: 0,
	}

	//bot.Reply(chatID, fmt.Sprintln(userID))
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
	case 0:
		// TODO: add validate if name is not empty
		OpsQueue[userID].Data = model.Category{
			ID:          uuid.New().String(),
			Title:       update.Message.Text,
			Description: "",
		}
		OpsQueue[userID].Step++
		bot.Reply(chatID, "Введите описание для категории:")
	case 1:
		category := OpsQueue[userID].Data.(model.Category)
		category.Description = update.Message.Text

		var message string
		if err := bot.CategoryRepository.Add(category); err != nil {
			if err == repository.ErrDocAlreadyExists {
				message = "Категория с таким именем уже существует! " + emoji.NoEntry
			} else {
				message = "Internal Server Error | write to @pdemian to get some help"
				bot.ReportToTheCreator(
					fmt.Sprintf("[CategoryRepository.Add] category: %+v | err: %s", category, err))
			}
			delete(OpsQueue, userID)
		} else {
			message = "Категория успешно добавлена " + emoji.Check
		}
		answer := tgbotapi.NewMessage(chatID, message)
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Bot.Send(answer)
	}
}
