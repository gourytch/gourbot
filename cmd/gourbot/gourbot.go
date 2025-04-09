package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"gourbot/internal/config"
	"gourbot/internal/logger"
	"gourbot/internal/tgbot"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Инициализация логгера
	logger := logger.InitLogger(cfg)

	// Инициализация Telegram-бота
	tgBot, err := tgbot.NewTgBot(cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize TgBot: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Запуск бота
	if err := tgBot.Start(ctx); err != nil {
		logger.Fatalf("TgBot stopped with error: %v", err)
	} else {
		logger.Info("TgBot finished without error")
	}
	logger.Info("That's all, folks!")
}
