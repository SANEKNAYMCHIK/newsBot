package services

import "time"

type NewsItem struct {
	Title       string
	Link        string
	Date        time.Time
	Description string
}

func NewNewsItem() *NewsItem {
	return &NewsItem{}
}
