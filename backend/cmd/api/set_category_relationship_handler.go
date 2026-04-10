package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/kegmor/dj-sets/backend/internal/repository"
)

type CreateSCRelationship struct {
	SetID string `json:"set_id"`
	Category string `json:"category"`
}

func setCategoryRelationship(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyRequest, error) {
	switch request.HTTPMethod{
	case "POST":
		

	}
	return request, fmt.Errorf("placeholder")
}