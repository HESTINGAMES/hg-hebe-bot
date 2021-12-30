package actions

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/hestingames/hg-hebe-bot/internal/logs"
)

const (
	welcomeMessage = "*Hola* %s\n" +
		"🙋🏻‍♀️ Bienvenid@ al grupo oficial de HestinGames para el servicio de 🔫 Counter-Strike: Global Offensive\n\n" +
		"🤝 Por favor respete las reglas (/rules) del grupo."
)

func SayHello(logger *logs.Logger, hebeBot tgbotapi.BotAPI, chatId int64, user tgbotapi.User) {
	hebeBot.Send(tgbotapi.NewChatAction(chatId, tgbotapi.ChatTyping))

	// Use user username in welcome message if have one
	name := user.FirstName
	if len(user.UserName) > 0 {
		name = fmt.Sprintf("@%s", user.UserName)
	}

	msg := tgbotapi.NewMessage(chatId, "")
	msg.ParseMode = "markdown"
	msg.Text = fmt.Sprintf(welcomeMessage, name)

	if _, err := hebeBot.Send(msg); err != nil {
		logger.Sugar().Error("Unable to send the welcome message :%s", err)
	}
}
