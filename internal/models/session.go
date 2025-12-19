package models

type Session struct {
	ID     string `redis:"id"`
	UserID string `redis:"user_id"`
}
