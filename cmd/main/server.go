package main

import (
	"log/slog"
	"math/rand"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/moonlags/markovBot/internal/markov"
)

type server struct {
	chain  *markov.Chain
	bot    *tgbotapi.BotAPI
	logger *slog.Logger
	images []string
	config config
}

func (s *server) run() error {
	if err := s.loadGobData(); err != nil {
		s.logger.Error("Can not load gob data", "err", err)
	}

	saveTicker := time.Tick(time.Minute * 5)

	updates := s.bot.GetUpdatesChan(tgbotapi.NewUpdate(0))
	for update := range updates {
		if len(saveTicker) > 0 {
			<-saveTicker

			if err := s.saveGobData(); err != nil {
				s.logger.Error("Can not save gob data", "err", err)
				os.Exit(1)
			}
		}

		if update.Message == nil {
			continue
		}

		if len(update.Message.Photo) > 0 {
			s.handlePhoto(update)
		}
		s.handleText(update)

		if rand.Intn(101) > s.config.chance {
			continue
		}

		text := s.chain.Generate(rand.Intn(5) + 3)
		s.logger.Info("response", "text", text)

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
