package main

import (
	"log/slog"
	"math/rand"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/moonlags/markovBot/internal/markov"
)

type server struct {
	bot    *tgbotapi.BotAPI
	logger *slog.Logger
	images []string
	config config
}

func (s *server) run() error {
	chain := markov.NewChain(s.config.prefixLen)
	if err := s.loadGobData(chain); err != nil {
		s.logger.Error(err.Error())
	}

	saveTicker := time.NewTicker(time.Minute * 5)
	defer saveTicker.Stop()

	updates := s.bot.GetUpdatesChan(tgbotapi.NewUpdate(0))
	for update := range updates {
		if len(saveTicker.C) > 0 {
			slog.Info("saving data")
			if err := s.saveGobData(chain); err != nil {
				s.logger.Error(err.Error())
				os.Exit(1)
			}
		}

		if update.Message == nil {
			continue
		}
		slog.Info("new message", "text", update.Message.Text, "user", update.Message.From.UserName)

		if update.FromChat().IsGroup() {
			chain.Add(strings.NewReader(update.Message.Text))
			chain.Add(strings.NewReader(update.Message.Caption))
		}

		if rand.Intn(101) > s.config.chance {
			continue
		}

		text := chain.Generate(rand.Intn(5) + 3)
		s.logger.Info("response", "text", text, "user", update.Message.From.UserName)

		if rand.Intn(101) < s.config.imageChance && len(s.images) > 0 {
			imageID := s.images[rand.Intn(len(s.images))]

			msg := tgbotapi.NewPhoto(update.FromChat().ID, tgbotapi.FileID(imageID))
			msg.Caption = text

			if _, err := s.bot.Send(msg); err != nil {
				return err
			}
		} else {
			msg := tgbotapi.NewMessage(update.FromChat().ID, text)
			if _, err := s.bot.Send(msg); err != nil {
				return err
			}
		}
	}
	return nil
}
