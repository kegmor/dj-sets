package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/kegmor/dj-sets/backend/internal/repository"
	"github.com/kegmor/dj-sets/backend/internal/youtube"
)

type SetService struct {
	db *repository.Queries
	yt *youtube.YouTube
}

func NewSetService(db *repository.Queries, yt *youtube.YouTube) *SetService {
	return &SetService{
		db: db, 
		yt: yt,
	}
}

func extractVideoId(webAddress string) (string, error) {
	u, err := url.Parse(webAddress)
	if err != nil {
		return "", fmt.Errorf("unable to parse url %w", err)
	}
	if u.Host == "youtu.be" {
		return strings.TrimPrefix(u.Path, "/"), nil
	}
	if id, found := strings.CutPrefix(u.Path, "/embed/"); found {
		return id, nil
	}
	videoId := u.Query().Get("v")
	if videoId == "" {
		return "", fmt.Errorf("no video id found in url")
	}
	return videoId, nil
}

func (s *SetService) CreateDjSet(webAddress string, djName string) (*repository.Set, error) {
	
	videoId, err := extractVideoId(webAddress)
	if err != nil {
		return &repository.Set{}, fmt.Errorf("failed to get video id %w", err)
	}
	
	gvd, err := s.yt.GetVideoDetails(videoId)
	if err != nil {
		return &repository.Set{}, fmt.Errorf("Unable to retrieve video data from youtube %w", err)
	}
	
	createdSet, err := s.db.CreateSet(context.Background(), repository.CreateSetParams{
		ID: uuid.New(),
		VideoID: videoId,
		Title: gvd.Title,
		DjName: djName,
		ChannelName: sql.NullString{String: gvd.ChannelTitle, Valid: true},
		Url: webAddress,
	})
	if err != nil {
		return &repository.Set{}, fmt.Errorf("failed to create dj set %w", err)
	}
	
	return &createdSet, nil
}

func (s *SetService) GetAllDjSets() ([]repository.Set, error) {
	djSets, err := s.db.GetAllSets(context.Background())
	if err != nil {
		return []repository.Set{}, fmt.Errorf("failed to get all dj sets %w", err)
	}
	return djSets, nil
}

func (s *SetService) GetDjSetById(id uuid.UUID) (repository.Set, error) {
	djSet, err := s.db.GetSetById(context.Background(), id)
	if err != nil {
		return repository.Set{}, fmt.Errorf("failed to get dj set by id %w", err)
	}
	return djSet, nil
}

func (s *SetService) DeleteDjSetById (id uuid.UUID) (repository.Set, error) {
	djSet, err := s.db.DeleteSetById(context.Background(), id)
	if err != nil {
		return repository.Set{}, fmt.Errorf("failed to delete dj set by id %w", err)
	}
	return djSet, nil
}
