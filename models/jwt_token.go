package models

import "time"

type TokenInfo struct {
	UserID  int
	Email   string
	IsAdmin bool
	Expired time.Time
}
