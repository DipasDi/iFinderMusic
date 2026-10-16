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

	// Подключение к Redis (OpenServer, Docker, локально)
	if _, err := storage.Init(cfg); err != nil {
		fmt.Println("Предупреждение по Redis:", err)
	}

	b, err := bot.New(cfg)
	if err != nil {
		fmt.Println("Ошибка инициализации бота:", err)
		os.Exit(1)
	}

	fmt.Println("🚀 Музыкальный бот успешно запущен!")
	if err := b.Start(); err != nil {
		fmt.Println("Ошибка работы бота:", err)
		os.Exit(1)
	}
}
