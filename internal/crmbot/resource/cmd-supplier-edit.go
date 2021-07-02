package resource

import (
	"fmt"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	"github.com/Haski007/crm-bot-the-sequel/pkg/validate"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	getSupplierFieldToEdit Step = iota
	getSupplierValueToEdit
)

func (bot *CrmBotService) commandSupplierEditHandler(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	if len(update.Message.Text) < 16 {
		bot.Errorf(chatID, "Wrong type of command!")
		return
	}
	supplierID := strings.ReplaceAll(update.Message.Text[len(update.Message.Text)-36:], "_", "-")

	OpsQueue[userID] = &Operation{
		Name: OperationType_SupplierEdit,
		Step: 0,
		Data: model.SupplierEdit{
			ID:    supplierID,
			Field: "",
		},
	}
	message := "Что нужно изменить " + emoji.QuestionMark
	answer := tgbotapi.NewMessage(chatID, message)
	answer.ReplyMarkup = keyboards.MarkupByArray([]string{
		model.SupplierEditName.String(),
		model.SupplierEditPhone.String(),
		model.SupplierEditDescription.String(),
	})
	answer.ParseMode = config.MarkdownParseMode
	bot.Bot.Send(answer)
}

func (bot *CrmBotService) hookSupplierEdit(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	op := OpsQueue[userID]
	switch op.Step {
	case getSupplierFieldToEdit:
		field, err := model.NewSupplierEditField(update.Message.Text)
		if err != nil {
			if err == model.ErrNoSuchSupplierEditField {
				bot.Reply(chatID, "Такого поля нет! "+emoji.NoEntry+"\nПопробуй ещё раз")
				return
			}
			bot.ReportToTheCreator(fmt.Sprintf("[hookSupplierEdit] NewSupplierEditField | err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			return
		}

		supEdit := OpsQueue[userID].Data.(model.SupplierEdit)
		supEdit.Field = field
		OpsQueue[userID].Data = supEdit
		OpsQueue[userID].Step++

		var answer tgbotapi.MessageConfig
		message := "Введите новое значение: "
		answer = tgbotapi.NewMessage(chatID, message)
		answer.ReplyMarkup = tgbotapi.NewHideKeyboard(false)
		bot.Bot.Send(answer)

	case getSupplierValueToEdit:
		supEdit := OpsQueue[userID].Data.(model.SupplierEdit)
		value := update.Message.Text

		if supEdit.Field == model.SupplierEditPhone {
			value = validate.PhoneNumber(value)
			if value == "" {
				bot.Reply(chatID, "Неверный формат номера телефона "+emoji.NoEntry+"\nПопробуй ещё раз")
				return
			}
		}

		if err := bot.SupplierRepository.UpdateField(supEdit.ID,
			supEdit.Field.BsonField(),
			value,
		); err != nil {
			if err == repository.ErrDocDoesNotExist {
				bot.Errorf(chatID, "Такого поставщика не существует")
				return
			}
			bot.ReportToTheCreator(fmt.Sprintf("[hookSupplierEdit] NewSupplierEditField | err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			delete(OpsQueue, userID)
			return
		}

		delete(OpsQueue, userID)
		message := "Поставщик успешно изменён " + emoji.Check
		answer := tgbotapi.NewMessage(chatID, message)
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Bot.Send(answer)
	}
}
