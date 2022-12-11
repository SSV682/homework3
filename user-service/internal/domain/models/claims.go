package models

import "time"

type Claims struct {
	ID     int64
	Expire time.Time
}
