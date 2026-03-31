package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/kegmor/dj-sets/backend/internal/repository"
)

type CreateSetRequest struct {
	URL		string `json:"url"`
	DjName	string `json:"dj_name"`
}

func extractDjAndUrl(body string) (*CreateSetRequest, error) {
	var req CreateSetRequest
	err := json.Unmarshal([]byte(body), &req)
	if err != nil {
		return nil, fmt.Errorf("invalid request body %w", err)
	}
	return &req, nil
}

func handleSets(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod{
	case "GET":
		result := []repository.Set{}
		if request.Path == "/sets" {
			data, err := set.GetAllDjSets(ctx)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:		fmt.Sprintf("failed to get dj sets %v", err),
				}, nil
			}
			result = data
		} else if strings.HasPrefix(request.Path, "/sets") { 
			id := request.PathParameters["id"]
			if id == "" {
				parts := strings.Split(request.Path, "/")
				id = parts[len(parts) - 1]
			}
			parsedID, err := uuid.Parse(id)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:		fmt.Sprintf("invalid set id %v", err),
				}, nil 
			}
			data, err := set.GetDjSetById(ctx, parsedID)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:		fmt.Sprintf("unable to get set by id %v", err),
				}, nil
			}
			result = data			
		}
		body, err := json.Marshal(result)
		if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 500,
					Body:		fmt.Sprintf("failed to marshal response %v", err),
				}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body: 		string(body),
			Headers: 	map[string]string{"Content-Type": "application/json"},
		}, nil				
	case "POST":
		if request.Path == "/sets" {
			data, err := extractDjAndUrl(request.Body)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:		fmt.Sprintf("failed to extract url and djName %v", err),
				}, nil
			}
			if data.DjName == "" || data.URL == "" {
				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:		fmt.Sprintf("url and dj name are required %v", err),
				}, nil
			}

			result, err := set.CreateDjSet(ctx, data.URL, data.DjName)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 500,
					Body:		fmt.Sprintf("failed to create dj set %v", err),
				}, nil
			}

			body, err := json.Marshal(result)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 500,
					Body:		fmt.Sprintf("failed to marshal response %v", err),
				}, nil
			}

			return events.APIGatewayProxyResponse{
				StatusCode: 201,
				Body: 		string(body),
				Headers: 	map[string]string{"Content-Type": "application/json"},
			}, nil
		}
	case "DELETE":
		if !strings.HasPrefix(request.Path, "/sets") {
			return events.APIGatewayProxyResponse{StatusCode: 404}, nil
		}
		id := request.PathParameters["id"]
		if id == "" {
			parts := strings.Split(request.Path, "/")
			id = parts[len(parts) - 1]
		}			
		parsedID, err := uuid.Parse(id)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:		fmt.Sprintf("invalid set id %v", err),
			}, nil 
		}
		result, err := set.DeleteDjSetById(ctx, parsedID)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:		fmt.Sprintf("unable to delete set by id %v", err),
			}, nil
		}
		body, err := json.Marshal(result)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:		fmt.Sprintf("failed to marshal response %v", err),
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body: 		string(body),
			Headers: 	map[string]string{"Content-Type": "application/json"},
		}, nil			
	default:
		return events.APIGatewayProxyResponse{StatusCode: 405}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 404}, nil
}