package services

import "time"

type NewsItem struct {
	Title       string
	Link        string
	Date        time.Time
	Description string
	Website     string
}

func NewNewsItem(title string, link string, date time.Time, descr string) *NewsItem {
	var site string =
	return &NewsItem{
		Title:       title,
		Link:        link,
		Date:        date,
		Description: descr,
		Website:     site,
	}
}
