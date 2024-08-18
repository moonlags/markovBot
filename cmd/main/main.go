package main

import (
	"flag"
	"log/slog"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"

	"github.com/moonlags/markovBot/internal/markov"
)

type config struct {
	prefixLen int
	chance    int
}

func init() {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stdout, nil)))

	if err := godotenv.Load(); err != nil {
		slog.Error("Can not load env variables", "err", err)
		os.Exit(1)
	}
}

func main() {
	var cfg config

	flag.IntVar(&cfg.prefixLen, "prefix", 2, "prefix length in words")
	flag.IntVar(&cfg.chance, "chance", 15, "chance of answering to a message")

	flag.Parse()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		slog.Error("Can not create bot api", "err", err)
		os.Exit(1)
	}

	server := server{
		chain:  markov.NewChain(cfg.prefixLen),
		config: cfg,
		bot:    bot,
	}

	if err := server.run(); err != nil {
		slog.Error("Can not run the server", "err", err)
		os.Exit(1)
	}
}
