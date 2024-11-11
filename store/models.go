package store

import "time"

type User struct {
	ID          int       `json:"id" storm:"id,increment"`
	Username       string    `json:"username" storm:"unique"`
	HashedPassword string    `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
}

type URLMapping struct {
	ID          int       `json:"id" storm:"id,increment"`
	UserId      string    `json:"user_id" storm:"index"`
	ShortUrl    string    `json:"short_url" storm:"index"`
	OriginalUrl string    `json:"original_url" storm:"index"`
	CreatedAt   time.Time `json:"created_at" storm:"index"`
}

type Session struct {
	ID           string    `json:"id" storm:"id,increment"`
	Username     string    `json:"username" storm:"index"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}
