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
	config config
}

func (s *server) run() error {
	if err := s.loadGobData(); err != nil {
		slog.Warn("Can not load gob data", "err", err)
	}

	saveTicker := time.Tick(time.Minute * 5)

	updates := s.bot.GetUpdatesChan(tgbotapi.NewUpdate(0))
	for update := range updates {
		if len(saveTicker) > 0 {
			<-saveTicker

			if err := s.saveGobData(); err != nil {
				slog.Error("Can not save gob data", "err", err)
				os.Exit(1)
			}
		}

		if update.Message == nil {
			continue
		}
		s.handleText(update)

		if rand.Intn(101) > s.config.chance {
			continue
		}

		text := s.chain.Generate(rand.Intn(10) + 3)
		slog.Info("response", "text", text)

		msg := tgbotapi.NewMessage(update.FromChat().ID, text)
		if _, err := s.bot.Send(msg); err != nil {
			return err
		}
	}
	return nil
}
