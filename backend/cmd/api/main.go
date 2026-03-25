package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"	
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/kegmor/dj-sets/backend/internal/database"
	_ "github.com/lib/pq"
)

type Secrets struct {
	Host		string `json:"host"`
	Port		int    `json:"port"`
	Username	string `json:"username"`
	Password 	string `json:"password"`
	DBName 		string `json:"dbname"`
}

type YouTube struct {
	YoutubeAPIKey	string `json:"api_key"`
}

var db *sql.DB
var youtubeKey string

func init() {
	secretName := "rds-credentials"
	region := os.Getenv("AWS_REGION")
	youtubeAPIKey := "youtube-api-key"

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

	ytKey := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(youtubeAPIKey),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal(err.Error())
	}

	yt, err := svc.GetSecretValue(context.TODO(), ytKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	var ytKeyString string = *yt.SecretString
	var ytak YouTube
	err = json.Unmarshal([]byte(ytKeyString), &ytak)
	if err != nil {
		log.Fatal(err)
	}
	youtubeKey = ytak.YoutubeAPIKey

	var secretString string = *result.SecretString

	var secret Secrets
	err = json.Unmarshal([]byte(secretString), &secret)
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
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body: "hello world",	
	}, nil
}

func main() {
	lambda.Start(handler)
}
