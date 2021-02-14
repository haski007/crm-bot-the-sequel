package resource

import (
	"fmt"
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	getCategoryFieldToEdit Step = iota
	getCategoryValueToEdit
)

func (bot *CrmBotService) commandCategoryEditHandler(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	if len(update.Message.Text) < 16 {
		bot.Errorf(chatID, "Wrong type of command!")
		return
	}
	categoryID := strings.ReplaceAll(update.Message.Text[len(update.Message.Text)-36:], "_", "-")

	OpsQueue[userID] = &Operation{
		Name: OperationType_CategoryEdit,
		Step: 0,
		Data: model.CategoryEdit{
			ID:    categoryID,
			Field: "",
		},
	}
	message := "Что нужно изменить " + emoji.QuestionMark
	answer := tgbotapi.NewMessage(chatID, message)
	answer.ReplyMarkup = keyboards.MarkupByArray([]string{
		model.CategoryEditTitle.String(),
		model.CategoryEditDescription.String(),
	})
	answer.ParseMode = "MarkDown"
	bot.Bot.Send(answer)
}

func (bot *CrmBotService) hookCategoryEdit(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	op := OpsQueue[userID]
	switch op.Step {
	case getCategoryFieldToEdit:
		field, err := model.NewCategoryEditField(update.Message.Text)
		if err != nil {
			if err == model.ErrNoSuchCategoryEditField {
				bot.Reply(chatID, "Такого поля нет! "+emoji.NoEntry+"\nПопробуй ещё раз")
				return
			}
			bot.ReportToTheCreator(fmt.Sprintf("[hookCategoryEdit] NewCategoryEditField | err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			delete(OpsQueue, userID)
			return
		}

		categoryEdit := OpsQueue[userID].Data.(model.CategoryEdit)
		categoryEdit.Field = field
		OpsQueue[userID].Data = categoryEdit
		OpsQueue[userID].Step++

		var answer tgbotapi.MessageConfig
		message := "Введите новое значение: "
		answer = tgbotapi.NewMessage(chatID, message)
		answer.ReplyMarkup = tgbotapi.NewHideKeyboard(false)
		bot.Bot.Send(answer)

	case getCategoryValueToEdit:
		categoryEdit := OpsQueue[userID].Data.(model.CategoryEdit)
		value := update.Message.Text

		if err := bot.CategoryRepository.UpdateField(
			categoryEdit.ID,
			categoryEdit.Field.BsonField(),
			value,
		); err != nil {
			if err == repository.ErrDocDoesNotExist {
				bot.Errorf(chatID, "Такой категории не существует")
				return
			}
			bot.ReportToTheCreator(fmt.Sprintf("[hookCategoryEdit] NewCategoryEditField | err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			delete(OpsQueue, userID)
			return
		}

		delete(OpsQueue, userID)
		message := "Категория успешно изменёна " + emoji.Check
		answer := tgbotapi.NewMessage(chatID, message)
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Bot.Send(answer)
	}
}
