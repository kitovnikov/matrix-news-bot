package ucase

import (
	"context"
	"fmt"
	"matrix-news-bot/config"
	"matrix-news-bot/globals"
	"time"
)

func CheckRoom(cfg *config.Config, ctx context.Context) error {

	timeNow := time.Now().Unix()

	if timeNow-globals.LastCheckRooms >= 300 {
		rooms, err := JoinedRoom(cfg, ctx)
		if err != nil {
			return err
		}

		fmt.Println("Получен список комнат")

		for _, roomID := range rooms {
			members, err := CountOfMembers(cfg, ctx, roomID)
			if err != nil {
				return err
			}

			if members <= 1 {
				fmt.Printf("Выходим")
				err := LeaveRoom(cfg, ctx, roomID)
				if err != nil {
					return err
				}
			}
			//спросить, добавлять ли это время чека в бд или просто в памяти приложения пусть будет
			globals.LastCheckRooms = timeNow
		}
	}
	return nil
}
