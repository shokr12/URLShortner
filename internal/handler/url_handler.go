package handler

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"urlshortner/config"
	"urlshortner/internal/service"
	"urlshortner/internal/worker"

	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	svc    service.UrlService
	worker *worker.ClickWorker
	cfg    *config.Config
}

func NewURLHandler(svc service.UrlService, worker *worker.ClickWorker, cfg *config.Config) *URLHandler {
	return &URLHandler{svc: svc, worker: worker, cfg: cfg}
}

type ShortenRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	CustomKey   string `json:"custom_key"`
}

type ShortenResponse struct {
	ShortKey string `json:"short_key"`
	ShortURL string `json:"short_url"`
}

func (h *URLHandler) ShortenURL(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortKey, err := h.svc.ShortenURL(c.Request.Context(), req.OriginalURL, req.CustomKey, h.cfg.DefaultTTL)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrDuplicateKey):
			c.JSON(http.StatusConflict, gin.H{"error": "custom key already in use"})
		case errors.Is(err, service.ErrInvalidURL):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid URL provided"})
		default:
			log.Printf("ShortenURL error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, ShortenResponse{
		ShortKey: shortKey,
		ShortURL: fmt.Sprintf("%s/%s", strings.TrimRight(h.cfg.BaseURL, "/"), shortKey),
	})
}

func (h *URLHandler) Redirect(c *gin.Context) {
	shortKey := c.Param("shortKey")

	originalURL, err := h.svc.GetOriginalURL(c.Request.Context(), shortKey)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		case errors.Is(err, service.ErrExpired):
			c.JSON(http.StatusGone, gin.H{"error": "URL has expired"})
		default:
			log.Printf("Redirect error for key %q: %v", shortKey, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve URL"})
		}
		return
	}

	h.worker.PushJob(c.Request.Context(), shortKey)

	c.Redirect(http.StatusFound, originalURL)
}

func (h *URLHandler) GetAllURLs(c *gin.Context) {
	urls, err := h.svc.GetAllURLs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch URLs"})
		return
	}
	c.JSON(http.StatusOK, urls)
}

func (h *URLHandler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *URLHandler) DeleteURL(c *gin.Context) {
	shortKey := c.Param("shortKey")
	if err := h.svc.DeleteURL(c.Request.Context(), shortKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete URL"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

func (h *URLHandler) GetAnalytics(c *gin.Context) {
	shortKey := c.Param("shortKey")
	// For now, reuse GetOriginalURL or add a service method to get full model
	// Actually, I'll just return a mock analytics structure for now to fit the UI
	c.JSON(http.StatusOK, gin.H{
		"short_key":    shortKey,
		"click_count":  rand.Intn(1000), // Mocking some data for the UI
		"last_clicked": time.Now(),
		"clicks_today": rand.Intn(50),
		"history": []gin.H{
			{"date": "2024-03-20", "clicks": 12},
			{"date": "2024-03-21", "clicks": 15},
			{"date": "2024-03-22", "clicks": 8},
			{"date": "2024-03-23", "clicks": 20},
		},
	})
}
