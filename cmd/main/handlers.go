package main

import (
	"log/slog"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *server) handlePhoto(update tgbotapi.Update) {
	slog.Info("photo", "caption", update.Message.Caption, "firstname", update.Message.From.FirstName)

	if !update.Message.Chat.IsGroup() {
		return
	}

	s.images = append(s.images, update.Message.Photo[0].FileID)
	s.chain.Add(strings.NewReader(update.Message.Caption))
}

func (s *server) handleText(update tgbotapi.Update) {
	slog.Info("message", "text", update.Message.Text, "firstname", update.Message.From.FirstName)

	if !update.Message.Chat.IsGroup() {
		return
	}

	s.chain.Add(strings.NewReader(update.Message.Text))
}
