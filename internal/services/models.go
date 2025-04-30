package services

import "time"

type NewsItem struct {
	Title       string
	Link        string
	Date        time.Time
	Description string
}

func NewNewsItem(title string, link string, date time.Time, descr string) *NewsItem {
	return &NewsItem{
		Title:       title,
		Link:        link,
		Date:        date,
		Description: descr,
	}
}
