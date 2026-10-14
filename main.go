package main

import (
	"fmt"
	"os"

	"random/internal/bot"
	"random/internal/config"
	"random/internal/storage"
)

func main() {
	cfg := config.Load()

	// Connect to Redis (OpenServer, Docker, or local instance)
	if _, err := storage.Init(cfg); err != nil {
		fmt.Println("Redis warning:", err)
	}

	b, err := bot.New(cfg)
	if err != nil {
		fmt.Println("Bot initialization error:", err)
		os.Exit(1)
	}

	fmt.Println("🚀 Music bot successfully started!")
	if err := b.Start(); err != nil {
		fmt.Println("Bot execution error:", err)
		os.Exit(1)
	}
}
