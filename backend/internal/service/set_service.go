package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	lambdaService "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/google/uuid"
	"github.com/kegmor/dj-sets/backend/internal/repository"
)

type SetService struct {
	db 				*repository.Queries
	lambdaClient 	*lambdaService.Client
	ytLambda 		string
}

type YouTubeLambdaRequest struct {
	VideoID string `json:"video_id"`
}

type YouTubeLambdaResponse struct {
	Title        string `json:"title"`
	ChannelTitle string `json:"channel_title"`
	Error        string `json:"error,omitempty"`
}

func NewSetService(db *repository.Queries, lambdaClient *lambdaService.Client, ytLambda string) *SetService {
	return &SetService{
		db:             db,
		lambdaClient:   lambdaClient,
		ytLambda: 		ytLambda,
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

func (s *SetService) fetchYouTubeDetails(videoID string) (*YouTubeLambdaResponse, error) {
	payload, err := json.Marshal(YouTubeLambdaRequest{VideoID: videoID})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal youtube request: %w", err)
	}
 
	result, err := s.lambdaClient.Invoke(context.Background(), &lambdaService.InvokeInput{
		FunctionName: aws.String(s.ytLambda),
		Payload:      payload,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to invoke youtube lambda: %w", err)
	}
 
	var response YouTubeLambdaResponse
	err = json.Unmarshal(result.Payload, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal youtube response: %w", err)
	}
 
	if response.Error != "" {
		return nil, fmt.Errorf("youtube lambda error: %s", response.Error)
	}
 
	return &response, nil
}

func (s *SetService) CreateDjSet(webAddress string, djName string) (*repository.Set, error) {
	
	videoId, err := extractVideoId(webAddress)
	if err != nil {
		return &repository.Set{}, fmt.Errorf("failed to get video id %w", err)
	}
	
	gvd, err := s.fetchYouTubeDetails(videoId)
	if err != nil {
		return &repository.Set{}, fmt.Errorf("Unable to retrieve video data from youtube %w", err)
	}
	
	createdSet, err := s.db.CreateSet(context.Background(), repository.CreateSetParams{
		ID: 		 uuid.New(),
		VideoID: 	 videoId,
		Title: 		 gvd.Title,
		DjName: 	 djName,
		ChannelName: gvd.ChannelTitle,
		Url: 		 webAddress,
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
