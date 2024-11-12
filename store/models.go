package store

import "time"

type User struct {
	ID             string    `json:"id" storm:"id"`
	Username       string    `json:"username" storm:"unique"`
	HashedPassword string    `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
}

type URLMapping struct {
	ID          string    `json:"id" storm:"id"`
	UserId      string    `json:"user_id" storm:"index"`
	ShortUrl    string    `json:"short_url" storm:"unique"`
	OriginalUrl string    `json:"original_url" storm:"index"`
	CreatedAt   time.Time `json:"created_at" storm:"index"`
}

type ClickStat struct {
	ID           string    `json:"id" storm:"id"`
	UrlMappingId string    `json:"url_mapping_id" storm:"index"`
	ClickTime    time.Time `json:"click_time" storm:"index"`
	IpAddress    string    `json:"ip_address" storm:"index"`
	Referer      string    `json:"referer" storm:"index"`
	UserAgent    string    `json:"user_agent" storm:"index"`
}

type Session struct {
	ID           string    `json:"id" storm:"id"`
	Username     string    `json:"username" storm:"index"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}
