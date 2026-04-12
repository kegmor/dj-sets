package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

func handleSearch(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod{
	case "GET":
		query := request.QueryStringParameters["q"]
		if query == ""{
			return respond(400, "could not get query paramater"), nil
		}
		result, err := search.Search(ctx, query)
		if err != nil {
			return respondError(400, "failed to retrieve query results",err)
		}
		return respondJSON(200, result)
	default:
    	return respond(405, "method not allowed"), nil
	}
}