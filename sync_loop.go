package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"matrix-news-bot/config"
	"matrix-news-bot/dto"
	"matrix-news-bot/globals"
	"matrix-news-bot/logging"
	"matrix-news-bot/storage"
	"matrix-news-bot/ucase"
	"net/http"
	"time"
)

func syncLoop(cfg *config.Config, ctx context.Context) {
	var since string

	batch, err := storage.GetLastBatch()
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logging.GetLogger(ctx).Fatalln("GetLastBatch err", err)
			return
		} else {
			since = ""
		}
	} else {
		since = batch.LastBatch
	}

	token, err := ucase.GetToken(cfg, ctx)
	if err != nil {
		logging.GetLogger(ctx).Fatalln("GetToken err", err)
		return
	}
	globals.AccessToken = token

	joinedRooms, err := ucase.JoinedRoom(cfg, ctx)
	if err != nil {
		return
	}
	for _, room := range joinedRooms {
		err := storage.AddRoom(room)
		if err != nil {
			return
		}
	}

	for {
		time.Sleep(3 * time.Second)
		url := fmt.Sprintf("https://%s/_matrix/client/v3/sync?access_token=%s&timeout=30000",
			cfg.HomeServerURL, globals.AccessToken)
		if since != "" {
			url += "&since=" + since
		}

		resp, err := http.Get(url)
		if err != nil {
			logging.GetLogger(ctx).Fatalln("Ошибка sync: %v", err)
			time.Sleep(3 * time.Second)
		}

		body, _ := io.ReadAll(resp.Body)
		err = resp.Body.Close()
		if err != nil {
			logging.GetLogger(ctx).Println("Ошибка close: %v", err)
			continue
		}

		if resp.StatusCode == http.StatusUnauthorized {
			token, err := ucase.GetToken(cfg, ctx)
			if err != nil {
				time.Sleep(600 * time.Second)
				continue
			}
			globals.AccessToken = token
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("Ошибка sync: %s", resp.Status)
			time.Sleep(3 * time.Second)
			continue
		}

		var syncResp globals.SyncResponse
		if err := json.Unmarshal(body, &syncResp); err != nil {
			logging.GetLogger(ctx).Printf("Ошибка JSON: %v", err)
		}

		processSync(cfg, ctx, syncResp)

		ucase.ParseRSS(cfg, ctx)

		err = ucase.CheckRoom(cfg, ctx)
		if err != nil {
			logging.GetLogger(ctx).Println("CheckRoom err", err)
			return
		}

		since = syncResp.NextBatch

		batch, err = storage.GetLastBatch()
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				logging.GetLogger(ctx).Println("GetLastBatch err", err)
				return
			} else {
				batch = &dto.Batch{}
			}
		}

		if batch.LastBatch == "" {
			err = storage.AddBatch(since)
			if err != nil {
				logging.GetLogger(ctx).Println("AddBatch err", err)
				return
			}
		} else if batch.LastBatch != "" {
			batch = &dto.Batch{
				ID:        batch.ID,
				LastBatch: since,
			}

			err = storage.UpdateBatch(*batch)
			if err != nil {
				logging.GetLogger(ctx).Println("UpdateBatch err", err)
				return
			}
		}
	}

}
