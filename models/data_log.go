package models

import "time"

// DataLog stores log
type DataLog struct {
	Date      time.Time
	Followers string
}
