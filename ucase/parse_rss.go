package ucase

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/mmcdole/gofeed"
	"io"
	"math/rand"
	"matrix-news-bot/globals"
	"matrix-news-bot/logging"
	_ "matrix-news-bot/logging"
	"net/http"
	"sort"
	"strconv"
	"time"

	"fmt"
	_ "github.com/mmcdole/gofeed"
	_ "log"
	"matrix-news-bot/config"
	_ "matrix-news-bot/dto"
	"matrix-news-bot/storage"
	_ "time"
)

func ParseRSS(cfg *config.Config, ctx context.Context) {
	checkTime := time.Duration(cfg.CheckTimeMinute)
	time.Sleep(checkTime * time.Minute)

	links, err := storage.GetRSSLinks(ctx)
	if err != nil {
		logging.GetLogger(ctx).Println("Ошибка при получении RSS ссылок из БД ", err)
		return
	}

	for _, link := range links {
		lastNewsTime, err := storage.GetLastNewsTime(link)
		if err != nil {
			logging.GetLogger(ctx).Println("Ошибка получения времени ", err)
			continue
		}

		url := link

		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(url)
		if err != nil {
			logging.GetLogger(ctx).Fatal(err)
		}

		sort.Slice(feed.Items, func(i, j int) bool {
			ti := feed.Items[i].PublishedParsed
			tj := feed.Items[j].PublishedParsed
			if ti == nil || tj == nil {
				return false
			}
			return ti.Before(*tj)
		})

		// Оставляем только последние 6 новостей
		if len(feed.Items) > 6 {
			feed.Items = feed.Items[len(feed.Items)-6:]
		}

		for _, item := range feed.Items {
			t, err := http.ParseTime(item.Published)
			if err != nil {
				logging.GetLogger(ctx).Println("Ошибка парсинга времени новости", err)
				continue
			}

			// Преобразуем в ISO 8601 для хранения в SQLite
			formatted := t.Format("2006-01-02 15:04:05")

			if lastNewsTime == "" || lastNewsTime < formatted {
				err = storage.UpdateLastNewsTime(formatted, link)
				if err != nil {
					logging.GetLogger(ctx).Println("Ошибка при обновлении времени последней новости ", err)
					return
				}
				newsTime, err := storage.GetLastNewsTime(link)
				if err != nil {
					logging.GetLogger(ctx).Println("Ошибка получения времени последней новости ", err)
					return
				}
				lastNewsTime = newsTime

				rooms, err := storage.GetAllRooms(ctx)
				if err != nil {
					logging.GetLogger(ctx).Println("Ошика получении всех комнат ", err)
					return
				}

				for _, roomID := range rooms {
					var text string
					if item.Description != "" {
						text = fmt.Sprintf("%s. %s<br><br><b>%s</b><br><br>%s<br><br><a href='%s'>Подробнее...</a>",
							feed.Title, t.Format("02.01.2006 15:04:05"), item.Title, item.Description, item.Link)
					} else {
						text = fmt.Sprintf("%s. %s<br><br><b>%s</b><br><br><a href='%s'>Подробнее...</a>",
							feed.Title, t.Format("02.01.2006 15:04:05"), item.Title, item.Link)
					}

					logging.GetLogger(ctx).Println("Отправляем сообщение в room " + roomID) // c
					err = sendFormattedMessage(cfg, ctx, roomID, text)
					if err != nil {
						logging.GetLogger(ctx).Println("Ошибка отправки сообщения ", err)
						continue
					}
					time.Sleep(3 * time.Second)
				}

			}
		}
	}

}

func sendFormattedMessage(cfg *config.Config, ctx context.Context, roomID string, text string) error {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	txnID := strconv.Itoa(rng.Intn(10000))

	message := map[string]string{
		"msgtype":        "m.notice",
		"body":           text,
		"format":         "org.matrix.custom.html",
		"formatted_body": text,
	}

	data, _ := json.Marshal(message)

	url := fmt.Sprintf("https://%s/_matrix/client/v3/rooms/%s/send/m.room.message/%s", cfg.HomeServerURL, roomID, txnID)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		logging.GetLogger(ctx).Println("Неудалось отправить запрос для отправки сообщения", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+globals.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logging.GetLogger(ctx).Println("Ошибка при отправке запроса: %v", err)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.GetLogger(ctx)
		}
	}(resp.Body)

	return nil
}
