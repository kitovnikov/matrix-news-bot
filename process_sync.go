package main

import (
	"context"
	"fmt"
	"matrix-news-bot/config"
	"matrix-news-bot/globals"
	"matrix-news-bot/storage"
	"matrix-news-bot/ucase"
)

func processSync(cfg *config.Config, ctx context.Context, syncResp globals.SyncResponse) {
	for roomID, inviteRaw := range syncResp.Rooms.Invite {
		rooms, err := storage.GetAllRooms(ctx)
		if err != nil {
			fmt.Println("Ошибка при получении комнат из БД:", err)
			return
		}
		//спросить, будет ли лучше если просто смотреть если ли комната такая уже или нет, true false
		for _, room := range rooms {
			if room == roomID {
				continue
			}
		}

		inviteMap := inviteRaw.(map[string]interface{})
		inviteState := inviteMap["invite_state"].(map[string]interface{})
		events := inviteState["events"].([]interface{})
		sender := events[0].(map[string]interface{})["sender"].(string)

		fmt.Println("Приглашение от: ", sender, ". В комнату: ", roomID)
		if err := ucase.JoinRoom(cfg, ctx, roomID); err == nil {
			err := storage.AddRoom(roomID)
			if err != nil {
				continue
			}
		}

		fmt.Println("Бот впервые вступил в комнату: ", roomID)
		if err := ucase.SendMessage(cfg, ctx, roomID, "Привет! Я новостной бот. Каждый день я буду присылать тебе свежие и важные новости."); err != nil {
			fmt.Println("Ошибка отправки:", err)
			continue
		}
	}
}
