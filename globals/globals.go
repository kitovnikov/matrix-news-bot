package globals

import "time"

var AccessToken string

type SyncResponse struct {
	NextBatch string `json:"next_batch"`
	Rooms     struct {
		Invite map[string]interface{} `json:"invite"`
		Join   map[string]interface{} `json:"join"`
	} `json:"rooms"`
}

var LastCheckRooms = time.Now().Unix()
