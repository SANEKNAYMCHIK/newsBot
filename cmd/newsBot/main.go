package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/SANEKNAYMCHIK/newsBot/internal/app/parser"
	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	token := os.Getenv("TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Printf("Some errors with tgbot: %s", err)
		panic(err)
	}
	fmt.Println(bot)
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
