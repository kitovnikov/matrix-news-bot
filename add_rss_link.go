package main

import (
	"context"
	"errors"
	"matrix-news-bot/config"
	"matrix-news-bot/storage"
	"strings"
)

var (
	ErrEmptyLinks = errors.New("Empty Links")
)

func addRSSLinkFromEnv(cfg *config.Config, ctx context.Context) error {
	links := cfg.RSSLinks
	if links == "" {
		return ErrEmptyLinks
	}

	dbLinks, err := storage.GetRSSLinks(ctx)
	if err != nil {
		return err
	}

	configLinks := strings.Split(links, ",")
	for _, configLink := range configLinks {
		cleanedConfig := strings.TrimSpace(configLink)
		if cleanedConfig == "" {
			continue
		}

		existsLink := false

		for _, dbLink := range dbLinks {
			if cleanedConfig == dbLink {
				existsLink = true
			}
		}

		if !existsLink {
			err = storage.AddRSSLink(cleanedConfig)
			if err != nil {
				return err
			}
		}

		if err := storage.AddRSSLink(cleanedConfig); err != nil {
			return err
		}
	}

	for _, dbLink := range dbLinks {
		existsDbLink := false
		for _, configLink := range configLinks {
			if configLink == dbLink {
				existsDbLink = true
			}
		}
		if !existsDbLink {
			err = storage.RemoveRSSLink(dbLink)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
