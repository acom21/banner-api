package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/acom21/banner-api/internal/models"
	"github.com/acom21/banner-api/pkg/database"
)

// Repository defines the interface for banner data operations
type Repository interface {
	GetBanner(ctx context.Context, id int) (*models.Banner, error)
	IncrementClick(ctx context.Context, bannerID int) error
	GetStats(ctx context.Context, bannerID int, from, to time.Time) ([]*models.ClickStat, error)
}

// repository implements Repository interface
type repository struct {
	db *database.DB
}

// New creates a new banner repository
func New(db *database.DB) Repository {
	return &repository{db: db}
}

// GetBanner retrieves a banner by ID
func (r *repository) GetBanner(ctx context.Context, id int) (*models.Banner, error) {
	query := "SELECT id, name FROM banners WHERE id = $1"

	var banner models.Banner
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&banner.ID, &banner.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get banner %d: %w", id, err)
	}

	return &banner, nil
}

// IncrementClick registers a click for a banner with minute-level aggregation
func (r *repository) IncrementClick(ctx context.Context, bannerID int) error {

	timestamp := time.Now().Truncate(time.Minute)

	query := `
		INSERT INTO click_stats (banner_id, timestamp, count)
		VALUES ($1, $2, 1)
		ON CONFLICT (banner_id, timestamp)
		DO UPDATE SET count = click_stats.count + 1`

	_, err := r.db.Pool.Exec(ctx, query, bannerID, timestamp)
	if err != nil {
		return fmt.Errorf("failed to increment click for banner %d: %w", bannerID, err)
	}

	return nil
}

// GetStats retrieves click statistics for a banner within a time range
func (r *repository) GetStats(ctx context.Context, bannerID int, from, to time.Time) ([]*models.ClickStat, error) {

	fromMinute := from.Truncate(time.Minute)
	toMinute := to.Truncate(time.Minute)

	query := `
		SELECT banner_id, timestamp, count
		FROM click_stats
		WHERE banner_id = $1 AND timestamp >= $2 AND timestamp <= $3
		ORDER BY timestamp ASC`

	rows, err := r.db.Pool.Query(ctx, query, bannerID, fromMinute, toMinute)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats for banner %d: %w", bannerID, err)
	}
	defer rows.Close()

	var stats []*models.ClickStat
	for rows.Next() {
		stat := &models.ClickStat{}
		if err := rows.Scan(&stat.BannerID, &stat.Timestamp, &stat.Count); err != nil {
			return nil, fmt.Errorf("failed to scan click stat: %w", err)
		}
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return stats, nil
}
