package main

import (
	"sync"

	"github.com/SANEKNAYMCHIK/newsBot/internal/app/parser"
	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
)

func main() {
	sources := []string{
		"https://habr.com/ru/rss/articles/",
		"https://russian.rt.com/rss",
		"https://lenta.ru/rss",
		"https://rss.nytimes.com/services/xml/rss/nyt/World.xml",
		"https://research.swtch.com/feed.atom",
	}
	// parser.Parse("https://habr.com/ru/rss/articles/")

	var wg sync.WaitGroup
	ch := make(chan services.NewsItem, 100)
	for _, url := range sources {
		wg.Add(1)
		go parser.Parse(url, &wg, ch)
	}
}
