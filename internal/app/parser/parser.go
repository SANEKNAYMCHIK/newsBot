package parser

import (
	"sync"

	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
)

func Parse(url string, wg *sync.WaitGroup, ch chan<- services.NewsItem) {
	// fp := gofeed.NewParser()
	// feed, _ := fp.ParseURL(url)
	// for _, item := range feed.Items {
	// 	fmt.Println(item.Title)
	// 	fmt.Println(item.Link)
	// 	fmt.Println(item.PublishedParsed)
	// 	fmt.Println(item.Description)
	// }
}
