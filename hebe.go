package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	csgoapi "github.com/hestingames/hg-hebe-bot/api"
	"github.com/hestingames/hg-hebe-bot/internal/environment"
)

var (
	statusMessages map[int64]int
)

func StartHebeBot() {
	statusMessages = make(map[int64]int)

	logger.Info("Initializing bot...")
	hebeBot, err := tgbotapi.NewBotAPI(AppConfig.BotToken)
	if err != nil {
		logger.Panic("Unable to initialize telegram bot")
	}

	if environment.IsLocal() {
		hebeBot.Debug = true
	}

	logger.Sugar().Infof("Authorized on account: %s", hebeBot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := hebeBot.GetUpdatesChan(u)

	for update := range updates {
		// Ignore any non-Message updates
		if update.Message == nil {
			continue
		}

		// Ignore any non-command Messages
		if !update.Message.IsCommand() {
			continue
		}

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "csgo":
			HandleStatus(*hebeBot, update)
		}
	}
}

const StatsRetriveErrorMessage = "El servicio se encuentra : *ONLINE*\n" +
	"Ha ocurrido un error al obtener las estadísticas 😅\n"

func HandleStatus(hebeBot tgbotapi.BotAPI, update tgbotapi.Update) {
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
