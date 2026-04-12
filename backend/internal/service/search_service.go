package service

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/kegmor/dj-sets/backend/internal/repository"
)

type SetSearchService struct {
	db *repository.Queries
}

func NewSetSearchService(db *repository.Queries) *SetSearchService {
	return &SetSearchService{
		db: db,
	}
}

func (sss *SetSearchService) Search(ctx context.Context, query string) ([]repository.Set, error){
	ch := make(chan []repository.Set, 4)
	var wg sync.WaitGroup
	wg.Add(4)
	go func(){
		defer wg.Done()
		sets, err := sss.db.GetSetsByDjName(ctx, query)
		if err == nil {
			ch <- sets
		}
	}()
	go func(){
		defer wg.Done()
		sets, err := sss.db.GetSetsByChannelName(ctx, query)
		if err == nil {
			ch <- sets
		}
	}()
	go func(){
		defer wg.Done()
		sets, err := sss.db.GetSetsByTitle(ctx, query)
		if err == nil {
			ch <- sets
		}
	}()
	go func(){
		defer wg.Done()
		sets, err := sss.db.GetSetsByCategory(ctx, query)
		if err == nil {
			ch <- sets
		}
	}()
	go func(){
		wg.Wait()
		close(ch)
	}()

	var allSets []repository.Set
	for sets := range ch {
		allSets = append(allSets, sets...)
	}
	seen := make(map[uuid.UUID]bool)
	var uniqueSets []repository.Set
	for _, set := range allSets {
		if !seen[set.ID] {
			seen[set.ID] = true
			uniqueSets = append(uniqueSets, set)
		}
	}	
	return uniqueSets, nil
}