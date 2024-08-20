package main

import (
	"log/slog"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *server) handleText(update tgbotapi.Update) {
	slog.Info("message", "text", update.Message.Text, "firstname", update.Message.From.FirstName)

	if !update.Message.Chat.IsGroup() || strings.HasPrefix(update.Message.Text, "/") {
		return
	}

	s.chain.Add(strings.NewReader(update.Message.Text))
}
