package migrate

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kegmor/dj-sets/backend/internal/database"
)

func handler(ctx context.Context, connStr string) error {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v",err)
	}

	return database.RunMigrations(db);
}

func main() {
	lambda.Start(handler)
}