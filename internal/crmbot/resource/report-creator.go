package resource

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func (bot *CrmBotService) ReportToTheCreator(report string) {
	bot.Bot.Send(tgbotapi.NewMessage(bot.Cfg.CreatorID, report))
}
