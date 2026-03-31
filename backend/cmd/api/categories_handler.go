package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	
)

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func handleCategories(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod{
	case "GET":		
		data, err := category.GetAllMusicCategories(ctx)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:		fmt.Sprintf("failed to get categories %v", err),
			}, nil
		}
		res, err := json.Marshal(data)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:		fmt.Sprintf("failed to marshal response %v", err),
			}, nil
		}		
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body: 		string(res),
			Headers: 	map[string]string{"Content-Type": "application/json"},
		}, nil						
	case "POST":		
		var cat CreateCategoryRequest
		err := json.Unmarshal([]byte(request.Body), &cat)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       fmt.Sprintf("invalid request body %v", err),
			}, nil
		}
		if cat.Name == "" {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "category name is required",
			}, nil
		}
		result, err := category.CreateCategory(ctx, cat.Name)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:		fmt.Sprintf("unable to create category %v", err),
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
	case "DELETE":
		name := request.PathParameters["name"]
		if name == "" {
			parts := strings.Split(request.Path, "/")
			name = parts[len(parts) - 1]
		}			
		result, err := category.DeleteCategory(ctx, name)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:		fmt.Sprintf("unable to delete category %v", err),
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
}