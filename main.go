package main

import (
	"fmt"
	"os"

	"random/internal/bot"
	"random/internal/config"
	"random/internal/storage"
)

// main initializes configuration, storage, and starts the Telegram bot service.
func main() {
	cfg := config.Load()

	if _, err := storage.Init(cfg); err != nil {
		fmt.Println("redis init warning:", err)
	}

	b, err := bot.New(cfg)
	if err != nil {
		fmt.Println("bot init error:", err)
		os.Exit(1)
	}

	fmt.Println("🚀 Music bot started")
	if err := b.Start(); err != nil {
		fmt.Println("bot runtime error:", err)
		os.Exit(1)
	}
}
