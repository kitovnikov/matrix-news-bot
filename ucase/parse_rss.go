package ucase

import (
	"context"
	"github.com/mmcdole/gofeed"
	"log"
	"matrix-news-bot/logging"
	_ "matrix-news-bot/logging"
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
	//checkTime := time.Duration(cfg.CheckTimeMinute)
	//time.Sleep(checkTime * time.Minute)

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

			fmt.Println(lastNewsTime < formatted, formatted, "LAST_NEWS_TIME: ", lastNewsTime)
			if lastNewsTime == "" || lastNewsTime < formatted {
				err := storage.UpdateLastNewsTime(formatted, link)
				if err != nil {
					return
				}
				newsTime, err := storage.GetLastNewsTime(link)
				if err != nil {
					return
				}
				lastNewsTime = newsTime
			}

			if lastNewsTime < formatted {
				time.Sleep(3 * time.Second)
				err = storage.UpdateLastNewsTime(formatted, link)
				if err != nil {
					fmt.Println("Ошибка при обновлении времени последней новости", err)
					return
				}

				rooms, err := storage.GetAllRooms(ctx)
				if err != nil {
					return
				}

				fmt.Printf("Канал :%s\nЗаголовок: %s\nОписание: %s\nСсылка: %s\nДата: %s",
					feed.Title, item.Title, item.Description, item.Link, formatted)
				for _, roomID := range rooms {
					var text string
					_ = roomID
					_ = text
					desc := fmt.Sprintf("Описание: %s", item.Description)
					if desc != "" {
						text = fmt.Sprintf("Канал: %s\nЗаголовок: %s\n %s\nСсылка: %s\nДата: %s",
							feed.Title, item.Title, desc, item.Link, formatted)
					} else {
						text = fmt.Sprintf("Канал: %s\nЗаголовок: %s\nСсылка: %s\nДата: %s",
							feed.Title, item.Title, item.Link, formatted)
					}

					err := SendMessage(cfg, ctx, roomID, text)
					if err != nil {
						return
					}
				}
			}
		}

	}
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
