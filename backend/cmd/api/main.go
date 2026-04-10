package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	lambdaService "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/kegmor/dj-sets/backend/internal/database"
	"github.com/kegmor/dj-sets/backend/internal/repository"
	"github.com/kegmor/dj-sets/backend/internal/service"
)



var corsHeaders = map[string]string{
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Headers": "Content-Type,x-api-key",
	"Access-Control-Allow-Methods": "GET,POST,DELETE,OPTIONS",
}

var set *service.SetService
var category *service.CatService
var set_categories *service.SetCategoryService

func init() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	
	ytLambda := os.Getenv("YOUTUBE_LAMBDA")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	lambdaClient := lambdaService.NewFromConfig(cfg)
	set = service.NewSetService(repository.New(db), lambdaClient, ytLambda)
	category = service.NewCatService(repository.New(db))
	set_categories = service.NewSetCategoryService(repository.New(db))
}

func respond(statusCode int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers:    corsHeaders,
	}
}
 
func respondJSON(statusCode int, data interface{}) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return respond(500, fmt.Sprintf("failed to marshal response: %v", err)), nil
	}
	return respond(statusCode, string(body)), nil
}
 
func respondError(statusCode int, message string, err error) (events.APIGatewayProxyResponse, error) {
	return respond(statusCode, fmt.Sprintf("%s: %v", message, err)), nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    if request.HTTPMethod == "OPTIONS" {
		return respond(200, ""), nil
	}
	
	if strings.HasPrefix(request.Path, "/sets") && strings.Contains(request.Path, "/categories") {
    	return handleSetCategories(ctx, request)
	} else if strings.HasPrefix(request.Path, "/sets") {
        return handleSets(ctx, request)
    } else if strings.HasPrefix(request.Path, "/categories") {
        return handleCategories(ctx, request)
    }
    return respond(404, "not found"), nil
}

func main() {
	lambda.Start(handler)
}
