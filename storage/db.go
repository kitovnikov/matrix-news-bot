package storage

import (
	"context"
	"database/sql"
	"errors"
	"matrix-news-bot/dto"
	"matrix-news-bot/logging"
)

func GetLastToken(nowDateTime string) (string, error) {
	var token string
	row := db.QueryRow(`SELECT token FROM auth_tokens where expired_at > ? ORDER BY id DESC LIMIT 1`, nowDateTime)
	err := row.Scan(&token)
	if err != nil {
		return token, err
	}

	return token, nil
}

func UpdateToken(token string, expiredAt string) error {
	var lastID int
	row := db.QueryRow(`SELECT id FROM auth_tokens ORDER BY id DESC LIMIT 1`)
	err := row.Scan(&lastID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = db.Exec(`INSERT INTO auth_tokens(token, expired_at) VALUES (?, ?)`, token, expiredAt)
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}
	_, err = db.Exec(`UPDATE auth_tokens SET token = ?, expired_at = ? WHERE id = ?`, token, expiredAt, lastID)
	if err != nil {
		return err
	}
	return nil
}

func AddRoom(roomID string) error {
	row := db.QueryRow(`SELECT room_id FROM rooms WHERE room_id=?`, roomID)
	var foundRoomID string
	err := row.Scan(&foundRoomID)
	if errors.Is(err, sql.ErrNoRows) {
		_, err := db.Exec(`INSERT OR IGNORE INTO rooms(room_id) VALUES (?)`, roomID)
		return err
	} else {
		return err
	}
}

func RemoveRoom(roomID string) error {
	_, err := db.Exec(`DELETE FROM rooms WHERE room_id=?`, roomID)
	return err
}

func GetAllRooms(ctx context.Context) ([]string, error) {
	rows, err := db.Query(`SELECT room_id FROM rooms`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logging.GetLogger(ctx)
		}
	}(rows)

	var rooms []string
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, id)
	}
	return rooms, nil
}

func AddBatch(batch string) error {
	_, err := db.Exec(`INSERT INTO batches(last_batch) VALUES (?)`, batch)
	return err
}

func GetLastBatch() (*dto.Batch, error) {
	var batch dto.Batch
	row := db.QueryRow(`SELECT id, last_batch FROM batches ORDER BY id DESC LIMIT 1`)
	err := row.Scan(&batch.ID, &batch.LastBatch)
	if err != nil {
		return nil, err
	}

	return &batch, nil
}

func UpdateBatch(batch dto.Batch) error {
	_, err := db.Exec(`UPDATE batches SET last_batch = ? WHERE id = ?`, batch.LastBatch, batch.ID)
	return err
}

func AddRSSLink(link string) error {
	_, err := db.Exec(`INSERT OR IGNORE INTO rss_links(link) VALUES (?)`, link)
	return err
}

func GetRSSLinks(ctx context.Context) ([]string, error) {
	rows, err := db.Query(`SELECT link FROM rss_links`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logging.GetLogger(ctx)
		}
	}(rows)

	var links []string
	for rows.Next() {
		var link string
		err := rows.Scan(&link)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

func GetLastNewsTime(link string) (string, error) {
	var lastTime sql.NullString
	row := db.QueryRow(`SELECT last_news_time FROM rss_links WHERE link=?`, link)
	err := row.Scan(&lastTime)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if lastTime.Valid {
		return lastTime.String, nil
	}
	return "", err
}

func UpdateLastNewsTime(time string, link string) error {
	_, err := db.Exec(`UPDATE rss_links SET last_news_time = ? WHERE link = ?`, time, link)
	return err
}

func RemoveRSSLink(link string) error {
	_, err := db.Exec(`DELETE FROM rss_links WHERE link = ?`, link)
	return err
}
