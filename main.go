package main

import (
	"context"

	"github.com/hestingames/hg-hebe-bot/api"
	"github.com/hestingames/hg-hebe-bot/internal/logs"
	"go.uber.org/zap"
)

// Check security issues
// gosec ./...

var (
	logger *logs.Logger
	ctx    context.Context
)

func main() {
	ctx = context.Background()
	logger = logs.FromContext(ctx)

	logger.Info("HebeBot - HestinGames")

	// Load Config
	logFn := func(key string, err error, msg string) {
		logger.Error(msg, zap.Error(err), zap.String("key", key))
	}
	loadConfig(logFn)

	// Initialize CSGO api client
	api.InitializeCsgoApi(AppConfig.ApiBaseUrl)

	// Initialize Telegram bot
	StartHebeBot()
}
