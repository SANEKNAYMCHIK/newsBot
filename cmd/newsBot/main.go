package main

import (
	"fmt"
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
	var wg sync.WaitGroup
	// channel for keeping news
	ch := make(chan services.NewsItem, 100)
	for _, url := range sources {
		wg.Add(1)
		go parser.Parse(url, &wg, ch)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	var news []services.NewsItem
	for item := range ch {
		news = append(news, item)
		fmt.Println(item.Link)
	}
}
