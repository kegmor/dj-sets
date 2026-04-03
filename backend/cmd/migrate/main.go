package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kegmor/dj-sets/backend/internal/database"
)

func handler(ctx context.Context) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return database.RunMigrations(db);
}

func main() {
	lambda.Start(handler)
}