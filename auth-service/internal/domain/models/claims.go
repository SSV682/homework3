package models

import "time"

type Claims struct {
	ID     string
	Expire time.Time
}
