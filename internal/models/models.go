package models

import "time"

type User struct {
	ID           int    `json:"id" db:"id"`
	TgChatID     *int64 `json:"tg_chat_id,omitempty" db:"tg_chat_id"`
	TgUsername   string `json:"tg_username,omitempty" db:"tg_username"`
	TgFirstName  string `json:"tg_first_name,omitempty" db:"tg_first_name"`
	Email        string `json:"email,omitempty" db:"email"`
	PasswordHash string `json:"password_hash,omitempty" db:"password_hash"`
	Role         string `json:"role" db:"role"`
}

type Category struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Source struct {
	ID         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	URL        string `json:"url" db:"url"`
	CategoryID *int   `json:"category_id,omitempty" db:"category_id"`
	IsActive   bool   `json:"is_active" db:"is_active"`
}

type NewsItem struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	URL         string    `json:"url" db:"url"`
	PublishedAt time.Time `json:"published_at" db:"published_at"`
	SourceID    int       `json:"source_id" db:"source_id"`
	GUID        string    `json:"guid" db:"guid"`
}

type UserSource struct {
	UserID   int `json:"user_id" db:"user_id"`
	SourceID int `json:"source_id" db:"source_id"`
}

type SentNews struct {
	NewsID int `json:"news_id" db:"news_id"`
	UserID int `json:"user_id" db:"user_id"`
}
