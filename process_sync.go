package main

import (
	"context"
	"matrix-news-bot/config"
	"matrix-news-bot/globals"
	"matrix-news-bot/logging"
	"matrix-news-bot/storage"
	"matrix-news-bot/ucase"
)

func processSync(cfg *config.Config, ctx context.Context, syncResp globals.SyncResponse) {
	for roomID, inviteRaw := range syncResp.Rooms.Invite {
		rooms, err := storage.GetAllRooms(ctx)
		if err != nil {
			logging.GetLogger(ctx).Println("Ошибка при получении комнат из БД:", err)
			return
		}

		for _, room := range rooms {
			if room == roomID {
				continue
			}
		}

		inviteMap := inviteRaw.(map[string]interface{})
		inviteState := inviteMap["invite_state"].(map[string]interface{})
		events := inviteState["events"].([]interface{})
		sender := events[0].(map[string]interface{})["sender"].(string)

		logging.GetLogger(ctx).Println("Приглашение от: ", sender, ". В комнату: ", roomID)
		if err := ucase.JoinRoom(cfg, ctx, roomID); err == nil {
			logging.GetLogger(ctx).Println("Ошибка при присоединении к комнате:", err)
			continue
		}

		logging.GetLogger(ctx).Println("Бот впервые вступил в комнату: ", roomID)
		if err := ucase.SendMessage(cfg, ctx, roomID, "Привет! Я новостной бот. Каждый день я буду присылать тебе свежие и важные новости."); err != nil {
			logging.GetLogger(ctx).Println("Ошибка отправки сообщения:", err)
			continue
		}
	}
}
