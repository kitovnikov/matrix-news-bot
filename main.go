package main

import (
	"context"
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
		logging.GetLogger(ctx).Fatalln("Ошибка при запуске БД", err)
		return
	}

	err = addRSSLinkFromEnv(cfg, ctx)
	if err != nil {
		logging.GetLogger(ctx).Fatalln("Ошибка при загрузке RSS ссылок", err)
		return
	}

	syncLoop(cfg, ctx)
}
