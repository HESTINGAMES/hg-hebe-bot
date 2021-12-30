package cmd

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	csgoapi "github.com/hestingames/hg-hebe-bot/api"
	"github.com/hestingames/hg-hebe-bot/internal/logs"
)

var (
	statusMessages map[int64]int
)

const StatsRetriveErrorMessage = "El servicio se encuentra : *ONLINE*\n" +
	"Ha ocurrido un error al obtener las estadísticas 😅\n"

func HandleStatus(logger *logs.Logger, hebeBot tgbotapi.BotAPI, update tgbotapi.Update) {
	// FIXME
	if len(statusMessages) == 0 {
		statusMessages = make(map[int64]int)
	}

	chatId := update.Message.Chat.ID

	msg := tgbotapi.NewMessage(chatId, "")
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "markdown" // html, markdown

	hebeBot.Send(tgbotapi.NewChatAction(chatId, tgbotapi.ChatTyping))

	msg.Text = "🕹 [HestinGames](http://hestingames.nat.cu)\n" +
		"🎮 *Counter-Strike: Global Offensive*\n" +
		"🛒 [Mercado](https://csgo.hestingames.nat.cu)\n\n"

	if playingNow, err := csgoapi.GetPlayingNow(); err != nil {
		msg.Text += StatsRetriveErrorMessage
	} else {
		msg.Text += fmt.Sprintf("📊 Estadísticas del Servicio 📊\n"+
			"🔫 Playing Now: %d\n\n", playingNow)

		if queueStatus, err := csgoapi.GetMatchakingQueueStatus(); err != nil {
			msg.Text += StatsRetriveErrorMessage
		} else {
			serverStatus := csgoapi.ParseServerStatus(queueStatus)

			// HACK : Sometimes the API fucks up and retrieves invalid queue info
			if playingNow != int(serverStatus.PlayingNow) {
				msg.Text += fmt.Sprintf("\n"+
					"📯 Matchmaking Casual 📯\n"+
					"Sigma: %d\n"+
					"Delta: %d\n"+
					"Dust II: %d\n"+
					"Hostages : %d\n",
					serverStatus.SigmaPlaying,
					serverStatus.DeltaPlaying,
					serverStatus.DustIIPlaying,
					serverStatus.HostagesPlaying)
			} else {
				msg.Text += StatsRetriveErrorMessage
			}
		}
	}

	// Remove older status message
	if messageId, ok := statusMessages[chatId]; ok {
		hebeBot.Send(tgbotapi.NewDeleteMessage(chatId, messageId))
	}

	// Send current status message
	if rsp, err := hebeBot.Send(msg); err == nil {
		statusMessages[chatId] = rsp.MessageID
	}
}
