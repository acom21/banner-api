package service

import (
	"context"
	"sync"
	"time"

	"github.com/acom21/banner-api/internal/models"
	"github.com/acom21/banner-api/internal/repository"
)

// Service defines the interface for banner business logic
type Service interface {
	RegisterClick(ctx context.Context, bannerID int) error
	GetStats(ctx context.Context, bannerID int, from, to time.Time) ([]*models.ClickStat, error)
}

// service implements Service interface
type service struct {
	repo  repository.Repository
	cache *bannerCache
}

type bannerCache struct {
	data map[int]time.Time
	mu   sync.RWMutex
	ttl  time.Duration
}

// New creates a new banner service
func New(repo repository.Repository) Service {
	return &service{
		repo: repo,
		cache: &bannerCache{
			data: make(map[int]time.Time),
			ttl:  5 * time.Minute,
		},
	}
}

// validateBanner checks if banner exists using cache
func (s *service) validateBanner(ctx context.Context, bannerID int) error {
	if s.cache.exists(bannerID) {
		return nil
	}

	_, err := s.repo.GetBanner(ctx, bannerID)
	if err == nil {
		s.cache.set(bannerID)
	}
	return err
}

// RegisterClick registers a click for a banner
func (s *service) RegisterClick(ctx context.Context, bannerID int) error {
	if err := s.validateBanner(ctx, bannerID); err != nil {
		return err
	}
	return s.repo.IncrementClick(ctx, bannerID)
}

// GetStats retrieves click statistics for a banner
func (s *service) GetStats(ctx context.Context, bannerID int, from, to time.Time) ([]*models.ClickStat, error) {
	if err := s.validateBanner(ctx, bannerID); err != nil {
		return nil, err
	}
	return s.repo.GetStats(ctx, bannerID, from, to)
}

// exists checks if banner is in cache and not expired
func (c *bannerCache) exists(bannerID int) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	lastChecked, ok := c.data[bannerID]
	return ok && time.Since(lastChecked) < c.ttl
}

// set adds banner to cache with cleanup
func (c *bannerCache) set(bannerID int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[bannerID] = time.Now()

	// Cleanup old entries if cache grows too large
	if len(c.data) > 1000 {
		cutoff := time.Now().Add(-c.ttl)
		for id, lastChecked := range c.data {
			if lastChecked.Before(cutoff) {
				delete(c.data, id)
			}
		}
	}
}
