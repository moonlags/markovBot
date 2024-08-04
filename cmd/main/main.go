package main

import (
	"flag"
	"log/slog"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type config struct {
	prefixLen   int
	chance      int
	imageChance int
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	var cfg config

	flag.IntVar(&cfg.prefixLen, "prefix", 1, "prefix length in words")
	flag.IntVar(&cfg.chance, "chance", 15, "chance of answering to a message")
	flag.IntVar(&cfg.imageChance, "image", 10, "chance of getting an image (1/chance) * (1/image)")

	flag.Parse()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	server := server{
		logger: logger,
		config: cfg,
		bot:    bot,
	}

	if err := server.run(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
