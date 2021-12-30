package events

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/hestingames/hg-hebe-bot/bot/actions"
	"github.com/hestingames/hg-hebe-bot/internal/logs"
)

func HandleNewChatMembers(logger *logs.Logger, hebeBot tgbotapi.BotAPI, update tgbotapi.Update) {
	for i := range update.Message.NewChatMembers {
		if !update.Message.NewChatMembers[i].IsBot {
			actions.SayHello(logger, hebeBot, update.Message.Chat.ID, update.Message.NewChatMembers[i])
		}
	}
}
