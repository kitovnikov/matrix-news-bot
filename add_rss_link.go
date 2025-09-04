package main

import (
	"context"
	"matrix-news-bot/config"
	"matrix-news-bot/storage"
	"strings"
)

func addRSSLinkFromEnv(cfg *config.Config, ctx context.Context) error {
	links := cfg.RSSLinks
	if links == "" {
		return nil
	}

	allRssLinks, err := storage.GetRSSLinks(ctx)
	if err != nil {
		return err
	}

	rssLinks := strings.Split(links, ",")
	for _, link := range rssLinks {
		cleaned := strings.TrimSpace(link)
		if cleaned == "" {
			continue
		}
		for _, linkRss := range allRssLinks {
			if linkRss != cleaned {
				err := storage.RemoveRSSLink(linkRss)
				if err != nil {
					return err
				}
			}
		}
		if err := storage.AddRSSLink(cleaned); err != nil {
			return err
		}
	}
	return nil
}
