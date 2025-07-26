package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/acom21/banner-api/internal/models"
	"github.com/acom21/banner-api/internal/service"
	"github.com/acom21/banner-api/pkg/logger"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// Handler defines the interface for banner operations
type Handler interface {
	RegisterClick(ctx *fasthttp.RequestCtx, bannerID int)
	GetStats(ctx *fasthttp.RequestCtx, bannerID int)
}

// BannerHandler handles HTTP requests for banner operations
type BannerHandler struct {
	service service.Service
}

// New creates a new banner handler
func New(svc service.Service) Handler {
	return &BannerHandler{
		service: svc,
	}
}

// RegisterClick registers a click for a banner
// @Summary Register banner click
// @Description Increments click counter for specified banner
// @Tags banners
// @Param bannerID path int true "Banner ID"
// @Success 200 "Click registered successfully"
// @Failure 400 "Invalid banner ID"
// @Failure 404 "Banner not found"
// @Failure 500 "Internal server error"
// @Router /counter/{bannerID} [get]
func (h *BannerHandler) RegisterClick(ctx *fasthttp.RequestCtx, bannerID int) {
	reqCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.service.RegisterClick(reqCtx, bannerID); err != nil {
		logger.Log.Warn("Banner not found", zap.Int("bannerID", bannerID), zap.Error(err))
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

// GetStats retrieves click statistics for a banner
// @Summary Get banner statistics
// @Description Returns click statistics for specified banner within time range
// @Tags banners
// @Param bannerID path int true "Banner ID"
// @Param request body models.StatsRequest true "Time range"
// @Success 200 {object} models.StatsResponse "Statistics data"
// @Failure 400 "Invalid request"
// @Failure 404 "Banner not found"
// @Failure 500 "Internal server error"
// @Router /stats/{bannerID} [post]
func (h *BannerHandler) GetStats(ctx *fasthttp.RequestCtx, bannerID int) {
	var req models.StatsRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if req.From == "" || req.To == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	from, err := time.Parse("2006-01-02T15:04:05", req.From)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	to, err := time.Parse("2006-01-02T15:04:05", req.To)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if to.Before(from) || to.Equal(from) {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	reqCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stats, err := h.service.GetStats(reqCtx, bannerID, from, to)
	if err != nil {
		logger.Log.Warn("Banner not found", zap.Int("bannerID", bannerID), zap.Error(err))
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	response := &models.StatsResponse{
		Stats: make([]models.StatEntry, len(stats)),
	}

	for i, stat := range stats {
		response.Stats[i] = models.StatEntry{
			Timestamp: stat.Timestamp.Format("2006-01-02T15:04:05"),
			Value:     stat.Count,
		}
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(response)
}
