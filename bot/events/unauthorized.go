package events

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hestingames/hg-hebe-bot/internal/logs"
)

const (
	unauthorizedMessage = "*Chat no autorizado*\n\n" +
		"🙅🏻‍♀️ Hebe no está autorizada a funcionar en este Chat\n" +
		"💁🏻‍♀️ Si cree que eso puede ser un error, contace a [HestinGames](https://t.me/hestingames)"
)

func HandleUnauthorizedMessage(logger *logs.Logger, hebeBot tgbotapi.BotAPI, update tgbotapi.Update) {
	chatId := update.Message.Chat.ID

	// Log attemp ?

	msg := tgbotapi.NewMessage(chatId, "")
	msg.ParseMode = "markdown"
	msg.Text = unauthorizedMessage

	if _, err := hebeBot.Send(msg); err != nil {
		logger.Sugar().Error("Unable to send the unauthorized message :%s", err)
	}

}
