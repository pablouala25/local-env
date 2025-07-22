package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
)

var (
	dbClient  *dynamodb.Client
	tableName string
)

// Item representa un registro en DynamoDB
type Item struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func init() {
	// Intenta cargar .env (opcional, si usas godotenv)
	_ = godotenv.Load()

	// Lee el endpoint de DynamoDB Local, si existe
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")

	// Carga la configuración de AWS
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(os.Getenv("AWS_REGION")),
	}
	if endpoint != "" {
		// Sobrescribe el endpoint para DynamoDB
		opts = append(opts, config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
				if service == dynamodb.ServiceID {
					return aws.Endpoint{
						URL:           endpoint,
						SigningRegion: os.Getenv("AWS_REGION"),
					}, nil
				}
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			}),
		))
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}
	dbClient = dynamodb.NewFromConfig(cfg)
	tableName = os.Getenv("TABLE_NAME")
	if tableName == "" {
		panic("TABLE_NAME environment variable is required")
	}
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parsea el body JSON
	var payload map[string]string
	if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid request"}, nil
	}

	// Crea un nuevo item
	item := Item{
		ID:      payload["id"],
		Message: payload["message"],
	}
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Marshal error"}, nil
	}

	// Escribe en DynamoDB
	_, err = dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "DynamoDB error"}, nil
	}

	body, _ := json.Marshal(item)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}
