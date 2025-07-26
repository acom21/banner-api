package http

import (
	"strconv"
	"strings"

	"github.com/acom21/banner-api/internal/handlers"

	"github.com/valyala/fasthttp"
)

// Router handles HTTP routing
type Router struct {
	handler handlers.Handler
}

// New creates a new router
func New(handler handlers.Handler) *Router {
	return &Router{handler: handler}
}

// Handler returns the fasthttp request handler
func (r *Router) Handler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		method := string(ctx.Method())

		if method == "GET" && strings.HasPrefix(path, "/counter/") {
			bannerID, valid := r.extractBannerID(path, "/counter/")
			if !valid {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			r.handler.RegisterClick(ctx, bannerID)
			return
		}

		if method == "POST" && strings.HasPrefix(path, "/stats/") {
			bannerID, valid := r.extractBannerID(path, "/stats/")
			if !valid {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			r.handler.GetStats(ctx, bannerID)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
}

// extractBannerID extracts and validates banner ID from path
func (r *Router) extractBannerID(path, prefix string) (int, bool) {
	bannerIDStr := strings.TrimPrefix(path, prefix)
	if bannerIDStr == "" {
		return 0, false
	}

	bannerID, err := strconv.Atoi(bannerIDStr)
	if err != nil || bannerID <= 0 {
		return 0, false
	}

	return bannerID, true
}
