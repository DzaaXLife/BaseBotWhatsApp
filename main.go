package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"whatsapp-bot/bot"
	"whatsapp-bot/config"
	"whatsapp-bot/logger"
)

func main() {
	log := logger.New()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config", "error", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b, err := bot.New(ctx, cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize bot", "error", err)
	}

	if err := b.Start(); err != nil {
		log.Fatal("Failed to start bot", "error", err)
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down bot...")
	b.Stop()
	log.Info("Bot stopped gracefully")
}
