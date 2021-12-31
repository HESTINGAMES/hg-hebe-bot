package hebe

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hestingames/hg-hebe-bot/bot/cmd"
	"github.com/hestingames/hg-hebe-bot/bot/events"
	"github.com/hestingames/hg-hebe-bot/config"
	"github.com/hestingames/hg-hebe-bot/internal/environment"
	"github.com/hestingames/hg-hebe-bot/internal/logs"
)

const (
	csgoGroup = 1001456543257
)

var (
	logger *logs.Logger
)

func Initialize(log *logs.Logger) {
	logger = log
}

func StartBot() {
	logger.Info("Initializing bot...")
	hebeBot, err := tgbotapi.NewBotAPI(config.AppConfig.BotToken)
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

		// Ignore message from other chats unless we are in local development environment
		if update.Message.Chat.ID != csgoGroup && !environment.IsLocal() {
			events.HandleUnauthorizedMessage(logger, *hebeBot, update)
			continue
		}

		if len(update.Message.NewChatMembers) != 0 {
			events.HandleNewChatMembers(logger, *hebeBot, update)
			continue
		}

		// Ignore any non-command Messages
		if !update.Message.IsCommand() {
			continue
		}

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "rules":
			cmd.HandleRules(logger, *hebeBot, update)
		case "csgo":
			cmd.HandleStatus(logger, *hebeBot, update)
		}
	}
}
