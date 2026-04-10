package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kegmor/dj-sets/backend/internal/repository"
)

type SetCategoryService struct {
	db *repository.Queries
}

func NewSetCategoryService(db *repository.Queries) *SetCategoryService {
	return &SetCategoryService{
		db: db,
	}
}

func (sc *SetCategoryService) AddCategoryToSet(ctx context.Context, set_id uuid.UUID, categoryName string) error {
	cat, err := sc.db.GetCategoryByName(ctx, categoryName)
	if err != nil {
		return fmt.Errorf("could not retrive category id %w", err)
	}
	err = sc.db.AddCategoryToSet(ctx, repository.AddCategoryToSetParams{
		SetID: set_id,
		CategoryID: cat.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create category %v", err)
	}
	return nil
}

func (sc *SetCategoryService) GetAllCategoriesForSet(ctx context.Context, setId uuid.UUID) ([]repository.Category, error) {
	cat, err := sc.db.GetAllCategoriesForSet(ctx, setId)
	if err != nil {
		return []repository.Category{}, fmt.Errorf("could not retrieve category id %w", err)
	}
	
	return cat, nil
}

func (sc *SetCategoryService) RemoveCategoryFromSet(ctx context.Context, set_id uuid.UUID, categoryName string) error {
	cat, err := sc.db.GetCategoryByName(ctx, categoryName)
	if err != nil {
		return fmt.Errorf("could not retrieve category id %w", err)
	}
	err = sc.db.RemoveCategoryFromSet(ctx, repository.RemoveCategoryFromSetParams{
		SetID: set_id,
		CategoryID: cat.ID,
	})
	if err != nil {
		return fmt.Errorf("could not remove category from set %w", err)
	}
	return nil
}