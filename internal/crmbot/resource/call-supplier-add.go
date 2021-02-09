package resource

import (
	"fmt"
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/pkg/validate"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
)

const (
	getSupplierName Step = iota
	getSupplierDescription
	getSupplierPhone
)

func (bot *CrmBotService) callSupplierAddHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	userID := update.CallbackQuery.From.ID

	// ---> adding to queue
	OpsQueue[userID] = &Operation{
		Name: OperationType_SupplierAdd,
		Step: 0,
	}

	delete := tgbotapi.NewDeleteMessage(chatID, messageID)
	bot.Bot.Send(delete)

	message := "Введите ФИО поставщика:"
	bot.Reply(chatID, message)
}

func (bot *CrmBotService) hookSupplierAdd(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	op := OpsQueue[userID]

	switch op.Step {
	case getSupplierName:
		OpsQueue[userID].Data = model.Supplier{
			ID:          uuid.New().String(),
			Name:        strings.TrimSpace(update.Message.Text),
			Description: "",
			Phone:       "",
		}
		OpsQueue[userID].Step++
		bot.Reply(chatID, "Введите описание поставщика:")
	case getSupplierDescription:
		supplier := OpsQueue[userID].Data.(model.Supplier)
		supplier.Description = strings.TrimSpace(update.Message.Text)
		OpsQueue[userID].Data = supplier

		OpsQueue[userID].Step++
		bot.Reply(chatID, "Введите номер телефона поставщика:")
	case getSupplierPhone:
		input := strings.TrimSpace(update.Message.Text)

		phone := validate.PhoneNumber(input)
		if phone == "" {
			bot.Reply(chatID, fmt.Sprintf("Неверный формат номера \"%s\"", input))
			return
		}

		supplier := OpsQueue[userID].Data.(model.Supplier)

		supplier.Phone = phone

		var message string
		if err := bot.SupplierRepository.Add(supplier); err != nil {
			if err == repository.ErrDocAlreadyExists {
				bot.Errorf(chatID,
					"Поставщик с таким именем {\"%s\"} уже существует",
					supplier.Name)
			} else {
				bot.Errorf(chatID,
					"Internal Server Error | write to @pdemian to get some help")
				bot.ReportToTheCreator(
					fmt.Sprintf("[SupplierRepository.Add] supplier: %+v | err: %s", supplier, err))
			}
			delete(OpsQueue, userID)
			return
		} else {
			delete(OpsQueue, userID)
			message = "Поставщик успешно добавлен " + emoji.Check
			answer := tgbotapi.NewMessage(chatID, message)
			answer.ReplyMarkup = keyboards.MainMenu
			bot.Bot.Send(answer)
		}
	}
}
