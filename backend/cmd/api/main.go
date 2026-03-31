package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	lambdaService "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/kegmor/dj-sets/backend/internal/database"
	"github.com/kegmor/dj-sets/backend/internal/repository"
	"github.com/kegmor/dj-sets/backend/internal/service"
	_ "github.com/lib/pq"
)

type Secrets struct {
	Host		string `json:"host"`
	Port		int    `json:"port"`
	Username	string `json:"username"`
	Password 	string `json:"password"`
	DBName 		string `json:"dbname"`
}

var db *sql.DB
var set *service.SetService
var category *service.CatService

func init() {
	secretName := "rds-credentials"
	region := os.Getenv("AWS_REGION")
	ytLambda := os.Getenv("YOUTUBE_LAMBDA")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager Client
	svc := secretsmanager.NewFromConfig(cfg)
	
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal(err.Error())
	}

	var secret Secrets
	err = json.Unmarshal([]byte(*result.SecretString), &secret)
	if err != nil {
		log.Fatal(err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", secret.Host, secret.Port, secret.Username, secret.Password, secret.DBName)

	var openError error
	db, openError = sql.Open("postgres", connStr)
	if openError != nil {
		log.Fatal(openError)
	}
	
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v",err)
	}

	err = database.RunMigrations(db)
	if err != nil {
		log.Fatal(err)
	}

	lambdaClient := lambdaService.NewFromConfig(cfg)
	set = service.NewSetService(repository.New(db), lambdaClient, ytLambda)
	category = service.NewCatService(repository.New(db))
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    if strings.HasPrefix(request.Path, "/sets") {
        return handleSets(ctx, request)
    } else if strings.HasPrefix(request.Path, "/categories") {
        return handleCategories(ctx, request)
    }
    return events.APIGatewayProxyResponse{StatusCode: 404}, nil
}

func main() {
	lambda.Start(handler)
}
