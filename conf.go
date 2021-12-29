package main

import (
	"github.com/hestingames/hg-hebe-bot/internal/distconf"
)

const configPath = "config.json"

type config struct {
	BotToken   string // Telegram HTTP Api bot token
	ApiBaseUrl string // CSGOGC api base url
}

var AppConfig *config

// Configuration will be read from top to bottom of the readers list.
func loadConfig(log distconf.Logger) {
	jconf := distconf.JSONConfig{}
	jconf.RefreshFile(configPath)
	readers := []distconf.Reader{&jconf} // &distconf.Env{}
	d := &distconf.Distconf{Logger: log, Readers: readers}

	AppConfig = &config{
		BotToken:   d.Str("BotToken", "invalid:token").Get(),
		ApiBaseUrl: d.Str("ApiBaseUrl", "http://127.0.0.1/").Get(),
	}
}
