package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Carga configuración AWS SDK v2 y el cliente DynamoDB :contentReference[oaicite:0]{index=0}
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return serverError(err)
	}
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if ep := os.Getenv("DYNAMODB_ENDPOINT"); ep != "" {
			o.BaseEndpoint = aws.String(ep)
		}
	})

	table := os.Getenv("DYNAMODB_TABLE_NAME")
	if table == "" {
		return clientError(http.StatusBadRequest, "DYNAMODB_TABLE_NAME not set")
	}

	switch req.HTTPMethod {
	case http.MethodPost:
		// POST /item → guarda id y message desde JSON en el body :contentReference[oaicite:1]{index=1}
		var payload struct {
			ID      string `json:"id"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
			return clientError(http.StatusBadRequest, "invalid JSON")
		}
		_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(table),
			Item: map[string]types.AttributeValue{
				"id":      &types.AttributeValueMemberS{Value: payload.ID},
				"message": &types.AttributeValueMemberS{Value: payload.Message},
			},
		})
		if err != nil {
			return serverError(err)
		}
		body, _ := json.Marshal(map[string]string{"status": "created", "id": payload.ID})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusCreated,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(body),
		}, nil

	case http.MethodGet:
		// GET /item/{id} → lee path parameter “id” y recupera ítem de Dynamo :contentReference[oaicite:2]{index=2}
		id := req.PathParameters["id"]
		if id == "" {
			return clientError(http.StatusBadRequest, "path parameter id required")
		}
		out, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(table),
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: id}},
		})
		if err != nil {
			return serverError(err)
		}
		if out.Item == nil {
			return clientError(http.StatusNotFound, "item not found")
		}
		msg := out.Item["message"].(*types.AttributeValueMemberS).Value
		body, _ := json.Marshal(map[string]string{"id": id, "message": msg})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(body),
		}, nil

	default:
		return clientError(http.StatusMethodNotAllowed, "method not allowed")
	}
}

func main() {
	lambda.Start(handler) // Arranca el handler como Lambda Proxy Integration :contentReference[oaicite:3]{index=3}
}

// clientError construye una respuesta 4xx
func clientError(status int, msg string) (events.APIGatewayProxyResponse, error) {
	body, _ := json.Marshal(map[string]string{"error": msg})
	return events.APIGatewayProxyResponse{StatusCode: status, Headers: map[string]string{"Content-Type": "application/json"}, Body: string(body)}, nil
}

// serverError construye una respuesta 500
func serverError(err error) (events.APIGatewayProxyResponse, error) {
	fmt.Fprintln(os.Stderr, err)
	body, _ := json.Marshal(map[string]string{"error": "internal server error"})
	return events.APIGatewayProxyResponse{StatusCode: 500, Headers: map[string]string{"Content-Type": "application/json"}, Body: string(body)}, nil
}
