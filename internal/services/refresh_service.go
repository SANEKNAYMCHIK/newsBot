package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/SANEKNAYMCHIK/newsBot/internal/repositories"
	"github.com/google/uuid"
)

type RefreshRequest struct {
	ID        string
	UserID    int64
	Timestamp time.Time
	Status    string
	Result    int
}

type RefreshService struct {
	rssService       *RssService
	subscriptionRepo repositories.SubscriptionRepository

	mu              sync.RWMutex
	userLastRequest map[int64]time.Time
	minRequestGap   time.Duration

	requestQueue chan *RefreshRequest
	maxQueueSize int
	workers      int

	requests sync.Map
}

func NewRefreshService(
	rssService *RssService,
	subscriptionRepo repositories.SubscriptionRepository,
	workers int,
	maxQueueSize int,
	minRequestGap time.Duration,
) *RefreshService {
	return &RefreshService{
		rssService:       rssService,
		subscriptionRepo: subscriptionRepo,
		userLastRequest:  make(map[int64]time.Time),
		minRequestGap:    minRequestGap,
		requestQueue:     make(chan *RefreshRequest, maxQueueSize),
		workers:          workers,
		maxQueueSize:     maxQueueSize,
	}
}

func (s *RefreshService) RequestRefresh(ctx context.Context, userID int64) (*RefreshRequest, error) {
	if !s.canRequestRefresh(userID) {
		return nil, fmt.Errorf("пожалуйста, подождите перед следующим обновлением")
	}
	if len(s.requestQueue) >= s.maxQueueSize {
		return nil, fmt.Errorf("очередь обновлений переполнена, попробуйте позже")
	}

	req := &RefreshRequest{
		ID:        uuid.New().String(),
		UserID:    userID,
		Timestamp: time.Now(),
		Status:    "pending",
	}

	s.requests.Store(req.ID, req)

	s.mu.Lock()
	s.userLastRequest[userID] = time.Now()
	s.mu.Unlock()

	select {
	case s.requestQueue <- req:
		req.Status = "queued"
		s.requests.Store(req.ID, req)
		log.Printf("Запрос на обновление добавлен в очередь: %s для пользователя %d", req.ID, userID)
		return req, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return nil, fmt.Errorf("очередь переполнена")
	}
}

func (s *RefreshService) canRequestRefresh(userID int64) bool {
	s.mu.RLock()
	lastRequest, exists := s.userLastRequest[userID]
	s.mu.RUnlock()
	if !exists {
		return true
	}
	return time.Since(lastRequest) >= s.minRequestGap
}

func (s *RefreshService) Start(ctx context.Context) {
	log.Printf("Starting RefreshService with %d workers", s.workers)
	for i := 0; i < s.workers; i++ {
		go s.processQueue(ctx)
	}
	go s.cleanupOldRequests(ctx)
}

func (s *RefreshService) processQueue(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-s.requestQueue:
			s.processRequest(ctx, req)
		}
	}
}

func (s *RefreshService) processRequest(ctx context.Context, req *RefreshRequest) {
	req.Status = "processing"
	s.requests.Store(req.ID, req)

	saved, err := s.rssService.FetchForUser(ctx, req.UserID)
	if err != nil {
		req.Status = "failed"
		log.Printf("Failed to refresh news for user %d: %v", req.UserID, err)
	} else {
		req.Status = "completed"
		req.Result = saved
		log.Printf("Completed refresh for user %d: saved %d items", req.UserID, saved)
	}

	s.requests.Store(req.ID, req)
}

func (s *RefreshService) GetRequestStatus(requestID string) (*RefreshRequest, bool) {
	if val, ok := s.requests.Load(requestID); ok {
		return val.(*RefreshRequest), true
	}
	return nil, false
}

func (s *RefreshService) cleanupOldRequests(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.requests.Range(func(key, value interface{}) bool {
				req := value.(*RefreshRequest)
				if time.Since(req.Timestamp) > 24*time.Hour {
					s.requests.Delete(key)
				}
				return true
			})

			s.mu.Lock()
			for userID, lastRequest := range s.userLastRequest {
				if time.Since(lastRequest) > 24*time.Hour {
					delete(s.userLastRequest, userID)
				}
			}
			s.mu.Unlock()
		}
	}
}
