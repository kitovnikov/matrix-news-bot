package ucase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"matrix-news-bot/config"
	"matrix-news-bot/logging"
	"matrix-news-bot/storage"
	"net/http"
	"time"
)

func GetToken(cfg *config.Config, ctx context.Context) (string, error) {
	tokenInDB, err := storage.GetLastToken(time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		logging.GetLogger(ctx).Println("Ошибка получения токена из БД", err)
	} else {
		return tokenInDB, nil
	}

	url := fmt.Sprintf("https://%s/_matrix/client/v3/login", cfg.HomeServerURL)

	payload := map[string]string{
		"type":     "m.login.password",
		"user":     cfg.Login,
		"password": cfg.Password,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("Ошибка при создании запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.GetLogger(ctx).Println("Ошибка закрытия get token", err)
			time.Sleep(300 * time.Second)
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("не удалось получить токен: %s\n%s", resp.Status, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access_token not found")
	}

	err = storage.UpdateToken(token, time.Now().Add(24*time.Hour).Format("2006-01-02 15:04:05"))
	if err != nil {
		logging.GetLogger(ctx).Println("Ошибка обновления токена в БД", err)
	}

	return token, nil
}
