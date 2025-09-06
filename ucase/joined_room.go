package ucase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"matrix-news-bot/config"
	"matrix-news-bot/globals"
	"matrix-news-bot/logging"
	"net/http"
)

func JoinedRoom(cfg *config.Config, ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("https://%s/_matrix/client/v3/joined_rooms?access_token=%s", cfg.HomeServerURL, globals.AccessToken)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при запросе joined_rooms: %s", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.GetLogger(ctx)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Не удалось получить список комнат: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Ошибка чтения ответа: %s", err)
	}

	var result struct {
		JoinedRooms []string `json:"joined_rooms"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Ошибка разбора JSON: %s", err)
	}

	return result.JoinedRooms, nil
}
