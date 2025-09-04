package main

import (
	"matrix-news-bot/config"
	"matrix-news-bot/storage"
	"strings"
)

func addRSSLinkFromEnv(cfg *config.Config) error {
	links := cfg.RSSLinks
	if links == "" {
		return nil
	}

	rssLinks := strings.Split(links, ",")
	for _, link := range rssLinks {
		cleaned := strings.TrimSpace(link)
		if cleaned == "" {
			continue
		}
		if err := storage.AddRSSLink(cleaned); err != nil {
			return err
		}
	}

	return nil
}
