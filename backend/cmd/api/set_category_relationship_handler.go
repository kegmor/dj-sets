package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type CreateSCRelationship struct {
	Category string `json:"category"`
}

func extractSetIdFromPath(path string) (uuid.UUID, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return uuid.UUID{}, fmt.Errorf("invalid path")
	}
	return uuid.Parse(parts[2])
}

func handleSetCategories(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	setID, err := extractSetIdFromPath(request.Path)
	if err != nil {
		return respondError(400, "failed to extract set id", err)
	}
	switch request.HTTPMethod{
	case "GET":
		cat, err := set_categories.GetAllCategoriesForSet(ctx, setID)
		if err != nil {
			return respondError(400, "could not retrieve all cateogries for a set", err)
		}
		return respondJSON(200, cat)
	case "POST":
		var cscr CreateSCRelationship
		err := json.Unmarshal([]byte(request.Body), &cscr)
		if err != nil {
			return respondError(400, "failed to get unmarshal to set categories", err)
		}
		if cscr.Category == "" {
			return respond(400, "category is required"), nil
		}		
		err = set_categories.AddCategoryToSet(ctx, setID, cscr.Category)
		if err != nil {
			return respondError(500, "failed to add category to set", err)
		}
		return respond(200, "category added to set"), nil
	
	case "DELETE":
		parts := strings.Split(request.Path, "/")
		if len(parts) < 5 {
			return respond(400, "category name is required in path"), nil
		}
		categoryName := parts[4]
		err := set_categories.RemoveCategoryFromSet(ctx, setID, categoryName)
		if err != nil {
			return respondError(400, "failed to remove category from set", err)
		}
		return respond(200, "category removed from set"), nil

	default:
		return respond(405, "method not allowed"), nil
	}
	return respond(405, "failed to handle request"), nil
}