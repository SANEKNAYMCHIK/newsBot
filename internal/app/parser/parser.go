package parser

import (
	"fmt"
	"sync"
	"time"

	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
	"github.com/mmcdole/gofeed"
)

const AMOUNT_NEWS int = 3

func Parse(url string, wg *sync.WaitGroup, ch chan<- services.NewsItem) {
	defer wg.Done()
	time.Sleep(5 * time.Second)
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		fmt.Printf("Error of parsing this url: %s with this error: %v", url, err)
	}
	for idx, item := range feed.Items {
		if idx == AMOUNT_NEWS {
			return
		}
		ch <- *services.NewNewsItem(item.Title, item.Link, *item.PublishedParsed, item.Description)
	}
}
