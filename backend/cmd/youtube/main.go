package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/kegmor/dj-sets/backend/internal/youtube"
)

type YouTubeSecret struct {
	APIKey string `json:"api_key"`
}

type YouTubeRequest struct {
	VideoID string `json:"video_id"`
}

type YouTubeResponse struct {
	Title        string `json:"title"`
	ChannelTitle string `json:"channel_title"`
	Error        string `json:"error,omitempty"`
}

var ytClient *youtube.YouTube

func init() {
	region := os.Getenv("AWS_REGION")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	svc := secretsmanager.NewFromConfig(cfg)
	result, err := svc.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String("youtube-api-key"),
		VersionStage: aws.String("AWSCURRENT"),
	})
	if err != nil {
		log.Fatal(err)
	}

	var secret YouTubeSecret
	err = json.Unmarshal([]byte(*result.SecretString), &secret)
	if err != nil {
		log.Fatal(err)
	}

	ytClient = youtube.NewYouTube(secret.APIKey, &http.Client{Timeout: 10 * time.Second})
}

func handler(ctx context.Context, request YouTubeRequest) (YouTubeResponse, error) {
	if request.VideoID == "" {
		return YouTubeResponse{Error: "video_id is required"}, fmt.Errorf("video_id is required")
	}

	details, err := ytClient.GetVideoDetails(request.VideoID)
	if err != nil {
		return YouTubeResponse{Error: err.Error()}, err
	}

	return YouTubeResponse{
		Title:        details.Title,
		ChannelTitle: details.ChannelTitle,
	}, nil
}

func main() {
	lambda.Start(handler)
}