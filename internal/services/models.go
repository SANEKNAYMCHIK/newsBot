package services

import "time"

type NewsItem struct {
	Title  string
	Link   string
	Date   time.Time
	Source string
}
