package main

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func handleCategories(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod{
	case "GET":		
		result, err := category.GetAllMusicCategories(ctx)
		if err != nil {
			return respondError(400, "failed to get categories", err)
		}
		return respondJSON(200, result)						
	case "POST":		
		var cat CreateCategoryRequest
		err := json.Unmarshal([]byte(request.Body), &cat)
		if err != nil {
			return respondError(400, "invalid request body", err)
		}
		if cat.Name == "" {
			return respond(400, "category name is required"), nil
		}
		result, err := category.CreateCategory(ctx, cat.Name)
		if err != nil {
			return respondError(400, "unable to create category", err)
		}
		return respondJSON(200, result)
	case "DELETE":
		name := request.PathParameters["name"]
		if name == "" {
			parts := strings.Split(request.Path, "/")
			name = parts[len(parts) - 1]
		}			
		result, err := category.DeleteCategory(ctx, name)
		if err != nil {
			return respondError(400, "unable to delete category", err)
		}
		return respondJSON(200, result)			
	default:
		return respond(405, "default failed categories message"), nil
	}
}