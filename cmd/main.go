package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"urlshortner/config"
	"urlshortner/internal/handler"
	"urlshortner/internal/repo"
	"urlshortner/internal/service"
	"urlshortner/internal/worker"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}
	// 1. Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Initialize Connections
	pgConn, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pgConn.Close(context.Background())

	rdbClient, err := config.InitRedis(cfg)
	if err != nil {
		log.Printf("Warning: failed to connect to redis: %v", err)
	}
	if rdbClient != nil {
		defer rdbClient.Close()
	}

	// 3. Initialize Repositories
	pgRepo := repo.NewUrlRepo(pgConn, rdbClient)
	cacheRepo := repo.NewCacheRepo(rdbClient)

	// 4. Initialize Worker
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	clickWorker := worker.NewClickWorker(pgRepo)
	go clickWorker.Start(workerCtx)

	// 5. Initialize Service & Handler
	urlService := service.NewUrlService(cacheRepo, *pgRepo)
	urlHandler := handler.NewURLHandler(urlService, clickWorker, cfg)

	// 6. Setup Router
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Simple CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.POST("/shorten", urlHandler.ShortenURL)
	r.GET("/:shortKey", urlHandler.Redirect)
	r.GET("/links", urlHandler.GetAllURLs)
	r.GET("/stats", urlHandler.GetStats)
	r.DELETE("/links/:shortKey", urlHandler.DeleteURL)
	r.GET("/analytics/:shortKey", urlHandler.GetAnalytics)

	// 7. Start server with graceful shutdown
	srv := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("starting server on %s", cfg.ServerAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// 8. Block until shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	log.Println("server stopped cleanly")
}
