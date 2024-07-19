package main

import (
	"github.com/codingconcepts/env"
)

type Env struct {
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN" required:"true"`
	TelegramChatID   string `env:"TELEGRAM_CHAT_ID" required:"true"`
	PublicWebhookURL string `env:"PUBLIC_WEBHOOK_URL" required:"true"`
}

func LoadEnv() (*Env, error) {
	var ev Env
	if err := env.Set(&ev); err != nil {
		return nil, err
	}
	return &ev, nil
}
