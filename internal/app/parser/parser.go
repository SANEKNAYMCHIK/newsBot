package parser

import (
	"fmt"
	"sync"

	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
	"github.com/mmcdole/gofeed"
)

func Parse(url string, wg *sync.WaitGroup, ch chan<- services.NewsItem) {
	defer wg.Done()
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		fmt.Printf("Error of parsing this url: %s with this error: %v", url, err)
	}
	for _, item := range feed.Items {
		ch <- *services.NewNewsItem(item.Title, item.Link, *item.PublishedParsed, item.Description)
	}
}
