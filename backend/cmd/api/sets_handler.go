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
				return respondError(400, "failed to get dj sets", err)
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
				return respondError(400, "invalid set id", err)
			}
			data, err := set.GetDjSetById(ctx, parsedID)
			if err != nil {
				return respondError(400, "unable to get set by id", err)
			}
			result = data			
		}
		return respondJSON(200, result)			
	case "POST":
		data, err := extractDjAndUrl(request.Body)
		if err != nil {
			return respondError(400, "failed to extract url and djName ", err)
		}
		if data.DjName == "" || data.URL == "" {
			return respond(400, "url and dj name are required"), nil
		}

		result, err := set.CreateDjSet(ctx, data.URL, data.DjName)
		if err != nil {
			return respondError(500, "failed to create dj set", err)
		}

		return respondJSON(200, result)
	case "DELETE":
		if !strings.HasPrefix(request.Path, "/sets") {
			return respond(404, "missing delete sets prefix"), nil
		}
		id := request.PathParameters["id"]
		if id == "" {
			parts := strings.Split(request.Path, "/")
			id = parts[len(parts) - 1]
		}			
		parsedID, err := uuid.Parse(id)
		if err != nil {
			return respondError(400,"invalid set id", err) 
		}
		result, err := set.DeleteDjSetById(ctx, parsedID)
		if err != nil {
			return respondError(400,"unable to delete set by id", err)
		}

		return respondJSON(200, result)			
	default:
		return respond(405, "default failed delete set operation message"), nil
	}
}