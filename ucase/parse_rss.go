package ucase

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/mmcdole/gofeed"
	"io"
	"log"
	"math/rand"
	"matrix-news-bot/globals"
	"matrix-news-bot/logging"
	_ "matrix-news-bot/logging"
	"net/http"
	"strconv"
	"time"

	"fmt"
	_ "github.com/mmcdole/gofeed"
	_ "log"
	"matrix-news-bot/config"
	_ "matrix-news-bot/dto"
	_ "matrix-news-bot/logging"
	"matrix-news-bot/storage"
	_ "time"
)

func ParseRSS(cfg *config.Config, ctx context.Context) {
	checkTime := time.Duration(cfg.CheckTimeMinute)
	time.Sleep(checkTime * time.Minute)

	links, err := storage.GetRSSLinks(ctx)
	if err != nil {
		fmt.Println("Ошибка при получении RSS ссылок из БД ", err)
		return
	}

	for _, link := range links {
		fmt.Println(link)
		lastNewsTime, err := storage.GetLastNewsTime(link)
		if err != nil {
			fmt.Println("Ошибка получения времени ", err)
			return
		}

		url := link

		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(url)
		if err != nil {
			logging.GetLogger(ctx).Fatal(err)
		}

		for _, item := range feed.Items {
			t, err := time.Parse(time.RFC1123Z, item.Published)
			if err != nil {
				log.Fatal(err)
			}

			// Преобразуем в ISO 8601 для хранения в SQLite
			formatted := t.Format("2006-01-02 15:04:05")

			if lastNewsTime == "" || lastNewsTime < formatted {
				err := storage.UpdateLastNewsTime(formatted, link)
				if err != nil {
					fmt.Println("Ошибка при обновлении времени последней новости ", err)

					return
				}
				newsTime, err := storage.GetLastNewsTime(link)
				if err != nil {
					fmt.Println("Ошибка получения времени последней новости ", err)
					return
				}
				lastNewsTime = newsTime

				rooms, err := storage.GetAllRooms(ctx)
				if err != nil {
					fmt.Println("Ошика получении всех комнат ", err)
					return
				}

				for _, roomID := range rooms {
					var text string
					if item.Description != "" {
						text = fmt.Sprintf("%s. %s<br><br><b>%s</b><br><br><b>%s</b><br><br>Ссылка: %s",
							feed.Title, formatted, item.Title, item.Description, item.Link)
					} else {
						text = fmt.Sprintf("%s. %s<br><br><b>%s</b><br><br>Ссылка: %s",
							feed.Title, formatted, item.Title, item.Link)
					}

					err := sendFormattedMessage(cfg, ctx, roomID, text)
					if err != nil {
						fmt.Println("Ошибка отправки сообщения ", err)
						return
					}
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

	fmt.Println("response Status:", resp.Status)

	return nil
}

//	links, err := storage.GetRSSLinks(ctx)
//	if err != nil {
//		fmt.Println("Ошибка при получении RSS ссылок из БД ", err)
//		return
//	}
//	fmt.Println("Links:", links)
//	fmt.Println("Начало парсинга")
//	var link string
//	for _, link = range links {
//		fmt.Println(link)
//		lastNewsTime, err := storage.GetLastNewsTime(link)
//		if err != nil {
//			fmt.Println("Ошибка получения времени ", err)
//			return
//		}
//		fmt.Println("СТАРОЕ ВРЕСЯ", lastNewsTime)
//		fmt.Println(link)
//		var oldTime time.Time
//		if lastNewsTime != "" {
//			locMSK, _ := time.LoadLocation("Europe/Moscow")
//			oldTime, err = time.ParseInLocation("2006-01-02 15:04:05", lastNewsTime, locMSK)
//			if err != nil {
//				fmt.Println("Ошибка форматирования старого времени", err)
//				return
//			}
//			fmt.Println("ФОРМАТИРОВАННОЕ СТАРОЕ ВРЕМЯ", oldTime)
//
//		}
//		fmt.Println("Обрабатываем ссылки: ", link)
//
//		url := link
//
//		fp := gofeed.NewParser()
//		feed, err := fp.ParseURL(url)
//		if err != nil {
//			logging.GetLogger(ctx).Fatal(err)
//		}
//
//		fmt.Println(feed.Title)
//		for _, item := range feed.Items {
//			t, err := time.Parse(time.RFC1123Z, item.Published)
//			if err != nil {
//				log.Fatal(err)
//			}
//
//			// Преобразуем в ISO 8601 для хранения в SQLite
//			formatted := t.Format("2006-01-02 15:04:05")
//
//			fmt.Println(2)
//			fmt.Println(lastNewsTime < formatted, formatted, "LAST_NEWS_TIME: ", lastNewsTime)
//			if lastNewsTime < formatted {
//				time.Sleep(3 * time.Second)
//				err = storage.UpdateLastNewsTime(formatted, link)
//				if err != nil {
//					fmt.Println("Ошибка при обновлении времени последней новости", err)
//					return
//				}
//
//				rooms, err := storage.GetAllRooms(ctx)
//				if err != nil {
//					return
//				}
//
//				fmt.Printf("Канал :%s\nЗаголовок: %s\nОписание: %s\nСсылка: %s\nДата: %s",
//					feed.Title, item.Title, item.Description, item.Link, formatted)
//				for _, roomID := range rooms {
//					var text string
//					_ = roomID
//					_ = text
//					desc := fmt.Sprintf("Описание: %s", item.Description)
//					if desc != "" {
//						text = fmt.Sprintf("Канал: %s\nЗаголовок: %s\n %s\nСсылка: %s\nДата: %s",
//							feed.Title, item.Title, desc, item.Link, formatted)
//					} else {
//						text = fmt.Sprintf("Канал: %s\nЗаголовок: %s\nСсылка: %s\nДата: %s",
//							feed.Title, item.Title, item.Link, formatted)
//					}
//
//					//err := SendMessage(cfg, ctx, roomID, text)
//					//if err != nil {
//					//	return
//					//}
//				}
//				fmt.Println(5)
//
//			}
//		}
//	}
//}
