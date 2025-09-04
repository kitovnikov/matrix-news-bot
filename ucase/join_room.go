package ucase

import (
	"context"
	"fmt"
	"io"
	"matrix-news-bot/config"
	"matrix-news-bot/globals"
	"matrix-news-bot/logging"
	"net/http"
)

func JoinRoom(cfg *config.Config, ctx context.Context, roomID string) error {
	url := fmt.Sprintf("https://%s/_matrix/client/v3/rooms/%s/join?access_token=%s", cfg.HomeServerURL, roomID, globals.AccessToken)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("Ошибка при создании запроса присоединения в комнату %s: %v", roomID, err)

	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Ошибка при отправке запроса присоединения в комнату %s: %v", roomID, err)

	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.GetLogger(ctx)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Принято приглашение и присоединился к комнате %s\n", roomID)
	} else {
		fmt.Printf("Не удалось присоединиться к комнате %s: %s\n", roomID, resp.Status)
	}

	return nil
}
