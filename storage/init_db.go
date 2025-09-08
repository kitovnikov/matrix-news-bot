package storage

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"matrix-news-bot/logging"
)

var db *sql.DB

func InitDB(ctx context.Context) error {
	var err error
	db, err = sql.Open("sqlite3", "./bot.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS rooms (
            id INTEGER PRIMARY KEY,
            room_id TEXT NOT NULL
        )
    `)

	if err != nil {
		logging.GetLogger(ctx).Fatal(err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS batches (
            id INTEGER PRIMARY KEY,
            last_batch TEXT NOT NULL
        )
    `)

	if err != nil {
		logging.GetLogger(ctx).Fatal(err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS rss_links (
            id INTEGER PRIMARY KEY,
            link TEXT NOT NULL UNIQUE,
            last_news_time TEXT
        )
    `)

	if err != nil {
		logging.GetLogger(ctx).Fatal(err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS auth_tokens (
            id INTEGER PRIMARY KEY,
            token TEXT NOT NULL UNIQUE,
            expired_at TEXT NOT NULL
        )
    `)

	if err != nil {
		logging.GetLogger(ctx).Fatal(err)
	}

	return nil
}
