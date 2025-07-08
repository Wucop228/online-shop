package models

import "time"

type RefreshToken struct {
	ID     int64
	Token  string
	UserID int64
	Expiry time.Duration
}
