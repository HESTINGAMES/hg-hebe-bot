package cmd

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hestingames/hg-hebe-bot/internal/logs"
)

const ChannelRules = "*Normas del Grupo*:\n\n" +
	"❌ Faltas de respeto u ofensas\n" +
	"❌ Temas referentes a Política o Religión\n" +
	"❌ Pornografía\n" +
	"❌ Spam\n\n" +
	"Consulte además las [Reglas del Servicio](https://csgo.hestingames.nat.cu/rules) antes de jugar en nuestros servidores de Counter-Strike: Global Offensive\n\n" +
	"El equipo de administración de HestinGames se reserva el derecho de tomar acciones administrativas ante cualquier comportamiento tóxico, radiactivo o inapropiado. Así como la prohibición de entrada al grupo a usuarios no deseados.\n"

func HandleRules(logger *logs.Logger, hebeBot tgbotapi.BotAPI, update tgbotapi.Update) {
	chatId := update.Message.Chat.ID
	hebeBot.Send(tgbotapi.NewChatAction(chatId, tgbotapi.ChatTyping))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "markdown"
	msg.Text = ChannelRules

	if _, err := hebeBot.Send(msg); err != nil {
		logger.Sugar().Error("Unable to send rules :%s", err)
	}
}
