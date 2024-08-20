package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/moonlags/markovBot/internal/markov"
	"github.com/moonlags/runware-go"
)

type server struct {
	chain   *markov.Chain
	bot     *tgbotapi.BotAPI
	runware *runware.Client
	config  config
}

func (s *server) run() error {
	if err := s.loadGobData(); err != nil {
		slog.Warn("Can not load gob data", "err", err)
	}

	saveTicker := time.Tick(time.Minute * 5)

	updates := s.bot.GetUpdatesChan(tgbotapi.NewUpdate(-1))
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

		var promptStart string
		if _, err := fmt.Sscanf(update.Message.Text, "/image %s", &promptStart); err == nil {
			promptIndex := strings.Index(update.Message.Text, promptStart)
			prompt := update.Message.Text[promptIndex:]

			slog.Info("generating image", "prompt", prompt)

			url, err := s.runware.TextToImage(runware.TextToImageArgs{
				Model:          "runware:100@1",
				PositivePrompt: prompt,
			})
			if err != nil {
				slog.Error("Can not generate image", "err", err)
				continue
			}

			msg := tgbotapi.NewPhoto(update.FromChat().ID, tgbotapi.FileURL(url))
			msg.Caption = s.chain.Generate(rand.Intn(10) + 3)

			if _, err := s.bot.Send(msg); err != nil {
				slog.Error("Can not send message", "err", err)
				continue
			}

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
			slog.Error("Can not send message", "err", err)
			continue
		}
	}
	return nil
}
