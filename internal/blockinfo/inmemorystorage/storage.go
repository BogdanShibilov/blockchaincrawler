package inmemorystorage

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/entity"
	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/repository"
)

type InMemoryStorage struct {
	repo repository.BlockRepo
	data []*entity.Block
	mu   *sync.RWMutex
	l    *zap.SugaredLogger
}

func NewInMem(repo repository.BlockRepo, l *zap.SugaredLogger) *InMemoryStorage {
	return &InMemoryStorage{
		repo: repo,
		data: make([]*entity.Block, 0),
		mu:   &sync.RWMutex{},
		l:    l,
	}
}

func (s *InMemoryStorage) Run(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	s.sync()
	ctx := context.TODO()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.sync()
		}
	}
}

func (s *InMemoryStorage) sync() {
	s.mu.Lock()
	defer s.mu.Unlock()

	blocks, err := s.repo.GetLastNBlocks(context.Background(), 20)
	if err != nil {
		s.l.Errorf("failed to sync inmem and db: %v", err)
	}

	s.data = blocks
}

func (s *InMemoryStorage) GetData() []*entity.Block {
	s.mu.RLock()
	defer s.mu.RUnlock()

	dataCopy := make([]*entity.Block, len(s.data))
	copy(dataCopy, s.data)
	return dataCopy
}
