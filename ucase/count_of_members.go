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

type JoinedMembersResponse struct {
	Joined map[string]interface{} `json:"joined"`
}

func CountOfMembers(cfg *config.Config, ctx context.Context, roomID string) (count int, err error) {
	url := fmt.Sprintf("https://%s/_matrix/client/v3/rooms/%s/joined_members?access_token=%s", cfg.HomeServerURL, roomID, globals.AccessToken)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("ошибка запроса участников: %v", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.GetLogger(ctx)
		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)

	var members JoinedMembersResponse
	if err := json.Unmarshal(body, &members); err != nil {
		return 0, fmt.Errorf("ошибка разбора JSON: %v", err)
	}

	return len(members.Joined), nil
}
