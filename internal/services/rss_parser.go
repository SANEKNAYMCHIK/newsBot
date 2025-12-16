package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
)

type RssItem struct {
	Title       string
	Link        string
	Date        time.Time
	Description string
	GUID        string
}

type RssParser struct {
	parser     *gofeed.Parser
	maxWorkers int
	maxRetries int
}

func NewRssParser(maxWorkers int) *RssParser {
	if maxWorkers <= 0 {
		maxWorkers = 10
	}
	return &RssParser{
		parser:     gofeed.NewParser(),
		maxWorkers: maxWorkers,
		maxRetries: 3,
	}
}

func (p *RssParser) ParseURL(URL string) ([]RssItem, error) {
	log.Printf("Started parse source: %s", URL)
	for i := range p.maxRetries {
		feed, err := p.parser.ParseURL(URL)
		if err == nil {
			var items []RssItem
			for _, item := range feed.Items {
				rssItem := p.convertToRssItem(item, URL)
				items = append(items, rssItem)
			}
			return items, nil
		}
		log.Printf("Can't parse on attempt %d", i+1)
		time.Sleep(time.Second)
	}
	return nil, fmt.Errorf("failed to parse URL: %s after %d attempts", URL, p.maxRetries)
}

func (p *RssParser) convertToRssItem(item *gofeed.Item, URL string) RssItem {
	publishTime := time.Now()
	if item.PublishedParsed != nil {
		publishTime = *item.PublishedParsed
	}
	guid := p.generateGUID(item, URL)
	return RssItem{
		Title:       strings.TrimSpace(item.Title),
		Link:        strings.TrimSpace(item.Link),
		Date:        publishTime,
		Description: strings.TrimSpace(item.Description),
		GUID:        guid,
	}
}

func (p *RssParser) generateGUID(item *gofeed.Item, URL string) string {
	if item.GUID != "" {
		return item.GUID
	}
	valueForCode := fmt.Sprintf("%s|%s|%v", URL, item.Link, item.Published)
	hash := sha256.Sum256([]byte(valueForCode))
	return hex.EncodeToString(hash[:])
}

func (p *RssParser) ParseURLsWithPool(urls []string) (map[string][]RssItem, error) {
	if len(urls) == 0 {
		return make(map[string][]RssItem), nil
	}
	tasks := make(chan string, len(urls))
	res := make(chan struct {
		url   string
		items []RssItem
		err   error
	}, len(urls))

	var wg sync.WaitGroup
	for i := 0; i < p.maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range tasks {
				items, err := p.ParseURL(url)
				res <- struct {
					url   string
					items []RssItem
					err   error
				}{url, items, err}
			}
		}()
	}
	for _, url := range urls {
		tasks <- url
	}
	close(tasks)
	go func() {
		wg.Wait()
		close(res)
	}()

	resItems := make(map[string][]RssItem)
	var hasErrors bool
	for val := range res {
		if val.err != nil {
			log.Printf("Error parsing\nURL:%s\nError:%v", val.url, val.err)
			hasErrors = true
			continue
		}
		resItems[val.url] = val.items
	}
	if hasErrors {
		return resItems, fmt.Errorf("some sources failed to parse")
	}
	return resItems, nil
}
