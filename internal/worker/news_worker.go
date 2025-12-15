package worker

import (
	"context"
	"log"
	"time"

	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
)

type NewsWorker struct {
	rssService *services.RssService
	interval   time.Duration
	isRunning  bool
	cancel     context.CancelFunc
}

func NewNewsWorker(rssService *services.RssService, interval time.Duration) *NewsWorker {
	return &NewsWorker{
		rssService: rssService,
		interval:   interval,
		isRunning:  false,
	}
}

func (w *NewsWorker) Start(ctx context.Context) {
	if w.isRunning {
		log.Println("NewsWorker is already running")
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	w.cancel = cancel
	w.isRunning = true

	log.Printf("Starting NewsWorker with interval %v", w.interval)

	go w.runTask(ctx)

	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("NewsWorker stopping...")
				w.isRunning = false
				return
			case <-ticker.C:
				go w.runTask(ctx)
			}
		}
	}()
}

func (w *NewsWorker) Stop() {
	if w.cancel != nil {
		w.cancel()
		w.isRunning = false
	}
}

func (w *NewsWorker) runTask(ctx context.Context) {
	log.Println("Starting news fetch task...")
	startTime := time.Now()

	saved, err := w.rssService.FetchAndSaveNews(ctx)
	if err != nil {
		log.Printf("Error in news fetch task: %v", err)
		return
	}

	duration := time.Since(startTime)
	log.Printf("News fetch task completed. Saved %d new items in %v", saved, duration)
}

func (w *NewsWorker) RunOnce(ctx context.Context) (int, error) {
	return w.rssService.FetchAndSaveNews(ctx)
}
