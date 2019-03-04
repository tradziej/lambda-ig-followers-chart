package models

import "time"

// DataLog stores log
type DataLog struct {
	Time      time.Time `json:"time"`
	Followers string    `json:"followers"`
	Username  string    `json:"username"`
}
