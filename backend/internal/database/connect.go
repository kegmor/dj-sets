package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	_ "github.com/lib/pq"
	
)

type Secrets struct {
	Host		string `json:"host"`
	Port		int    `json:"port"`
	Username	string `json:"username"`
	Password 	string `json:"password"`
	DBName 		string `json:"dbname"`
}

func Connect() (*sql.DB, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config %w", err)
	}

	// Create Secrets Manager Client
	svc := secretsmanager.NewFromConfig(cfg)
	
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String("rds-credentials"),
		VersionStage: aws.String("AWSCURRENT"),
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
	return sql.Open("postgres", connStr)
}