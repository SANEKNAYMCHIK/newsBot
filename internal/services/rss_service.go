package services

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/internal/repositories"
)

type RssService struct {
	sourceRepo repositories.SourceRepository
	newsRepo   repositories.NewsRepository
	parser     *RssParser
}

func NewRssService(
	sourceRepo repositories.SourceRepository,
	newsRepo repositories.NewsRepository,
	parser *RssParser,
) *RssService {
	return &RssService{
		sourceRepo: sourceRepo,
		newsRepo:   newsRepo,
		parser:     parser,
	}
}

func (s *RssService) FetchAndSaveNews(ctx context.Context) (int, error) {
	sources, err := s.sourceRepo.GetActive(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get active source: %w", err)
	}
	if len(sources) == 0 {
		log.Println("No active sources found")
		return 0, nil
	}
	log.Println(len(sources))
	var urls []string
	for _, source := range sources {
		urls = append(urls, source.URL)
	}
	parsedResults, err := s.parser.ParseURLsWithPool(urls)
	if err != nil {
		log.Printf("Parser completed with errors: %w", err)
		// Часть источников не спарсилась, но все равно продолжаем работу дальше
	}
	var totalSaved int
	var mu sync.Mutex

	var wg sync.WaitGroup
	maxSavers := 5

	saveTask := make(chan struct {
		source models.Source
		items  []RssItem
	}, len(sources))

	for i := 0; i < maxSavers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range saveTask {
				saved := s.saveSourceNews(ctx, task.source, task.items)
				mu.Lock()
				totalSaved += saved
				log.Printf("Successfully saved %d new news items from %s", saved, task.source.Name)
				mu.Unlock()
			}
		}()
	}
	for _, source := range sources {
		if items, ok := parsedResults[source.URL]; ok && len(items) > 0 {
			saveTask <- struct {
				source models.Source
				items  []RssItem
			}{source, items}
		}
	}
	close(saveTask)
	wg.Wait()
	return totalSaved, nil
}

func (s *RssService) saveSourceNews(ctx context.Context, source models.Source, items []RssItem) int {
	var saved int

	for _, item := range items {
		content := item.Description
		newsItem := &models.NewsItem{
			Title:       item.Title,
			Content:     &content,
			URL:         item.Link,
			PublishedAt: item.Date,
			SourceID:    source.ID,
			GUID:        item.GUID,
		}

		exists, err := s.newsRepo.ExistsByGUID(ctx, int(source.ID), item.GUID)
		if err != nil {
			log.Printf("Error checking existence for GUID %s: %v", item.GUID, err)
			continue
		}

		if exists {
			continue
		}

		if err := s.newsRepo.Create(ctx, newsItem); err != nil {
			log.Printf("Failed to save news '%s': %v", item.Title, err)
			continue
		}
		saved++
	}
	return saved
}

func (s *RssService) FetchForUser(ctx context.Context, userID int64) (int, error) {
	sources, err := s.sourceRepo.GetActiveForUser(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user subscriptions: %w", err)
	}

	if len(sources) == 0 {
		return 0, nil
	}

	var saved int
	for _, source := range sources {
		items, err := s.parser.ParseURL(source.URL)
		if err != nil {
			log.Printf("Failed to parse source %s for user %d: %v", source.Name, userID, err)
			continue
		}

		count := s.saveSourceNews(ctx, source, items)
		saved += count
	}

	return saved, nil
}
