package worker

import (
	"context"
	"log"
	"urlshortner/internal/repo"
)

type ClickWorker struct {
	repo    *repo.UrlRepo
	JobChan chan string
}

func NewClickWorker(repo *repo.UrlRepo) *ClickWorker {
	return &ClickWorker{
		repo:    repo,
		JobChan: make(chan string, 100), // Buffer size of 100
	}
}

func (w *ClickWorker) Start(ctx context.Context) {
	log.Println("Starting Click Worker...")
	for {
		select {
		case shortKey := <-w.JobChan:
			err := w.repo.IncrementClickCount(shortKey, context.Background())
			if err != nil {
				log.Printf("Worker error incrementing click count for %s: %v\n", shortKey, err)
			}
		case <-ctx.Done():
			log.Println("Stopping Click Worker...")
			return
		}
	}
}

func (w *ClickWorker) PushJob(ctx context.Context, shortKey string) {
	select {
	case w.JobChan <- shortKey:
	default:
		log.Printf("Worker queue full, dropping job for %s\n", shortKey)
	}
}
