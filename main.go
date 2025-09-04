package main

import (
	"context"
	"fmt"
	"matrix-news-bot/config"
	"matrix-news-bot/logging"
	"matrix-news-bot/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.MustLoad()

	logger := logging.NewLogger()
	ctx = logging.ContextWithLogger(ctx, logger)

	logging.GetLogger(ctx).Infoln("Запуск бота")

	err := storage.InitDB(ctx)
	if err != nil {
		fmt.Println("Ошибка при запуске БД", err)
		return
	}

	err = addRSSLinkFromEnv(cfg)
	if err != nil {
		fmt.Println("Ошибка при загрузке RSS ссылок", err)
		return
	}

	syncLoop(cfg, ctx)
}
