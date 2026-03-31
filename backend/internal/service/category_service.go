package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kegmor/dj-sets/backend/internal/repository"
)

type CatService struct {
	db 	*repository.Queries
}

func NewCatService (db *repository.Queries) *CatService {
	return &CatService{
		db: db,
	}
}

func (c *CatService) CreateCategory(ctx context.Context, category string) (*repository.Category, error) {
	createdCategory, err := c.db.CreateCategory(ctx, repository.CreateCategoryParams{
		ID: uuid.New(),
		Name: category,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create category %v", err)
	}
	return &createdCategory, nil
}

func (c *CatService) GetAllMusicCategories(ctx context.Context) ([]repository.Category, error) {
	categories, err := c.db.GetAllCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all music categories %w", err)
	}
	return categories, nil
}

func (c *CatService) DeleteCategory(ctx context.Context, name string) (*repository.Category, error) {
	deletedCategory, err := c.db.DeleteCategoryByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to delete category by name %w", err)
	}
	return &deletedCategory, nil
}

func (c *CatService) GetCategory(ctx context.Context, name string) (*repository.Category, error) {
	retrievedCategory, err := c.db.GetCategoryByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by name %w", err)
	}
	return &retrievedCategory, nil
}

