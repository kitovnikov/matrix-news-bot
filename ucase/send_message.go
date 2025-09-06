package ucase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"matrix-news-bot/config"
	"matrix-news-bot/globals"
	"matrix-news-bot/logging"
	"net/http"
	"strconv"
	"time"
)

func SendMessage(cfg *config.Config, ctx context.Context, roomID string, text string) error {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	txnID := strconv.Itoa(rng.Intn(10000))

	message := map[string]string{
		"msgtype": "m.notice",
		"body":    text,
	}

	data, _ := json.Marshal(message)

	url := fmt.Sprintf("https://%s/_matrix/client/v3/rooms/%s/send/m.room.message/%s", cfg.HomeServerURL, roomID, txnID)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("Неудалось отправить запрос для отправки сообщения", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+globals.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Ошибка при отправке запроса: %v", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.GetLogger(ctx)
		}
	}(resp.Body)

	return nil
}
