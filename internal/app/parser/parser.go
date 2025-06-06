package parser

import (
	"log"
	"sync"

	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
	"github.com/mmcdole/gofeed"
)

const AMOUNT_NEWS int = 3

func Parse(url string, wg *sync.WaitGroup, ch chan<- services.NewsItem) {
	defer wg.Done()
	// time.Sleep(5 * time.Second)
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		log.Printf("Error of parsing this url: %s with this error: %v", url, err)
	}
	for idx, item := range feed.Items {
		if idx == AMOUNT_NEWS {
			return
		}
		ch <- *services.NewNewsItem(item.Title, item.Link, *item.PublishedParsed, item.Description)
	}
}

func ParseAllNews(sources []string) map[string][]services.NewsItem {
	var wg sync.WaitGroup
	ch := make(chan services.NewsItem, 100)
	for _, url := range sources {
		wg.Add(1)
		go Parse(url, &wg, ch)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	newsVals := make(map[string][]services.NewsItem)
	for item := range ch {
		newsVals[item.Website] = append(newsVals[item.Website], item)
	}
	return newsVals
}
