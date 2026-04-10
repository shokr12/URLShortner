package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
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
