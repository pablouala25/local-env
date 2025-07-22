package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func handler(ctx context.Context) (string, error) {
	// 1. Carga configuración apuntando al endpoint local
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           os.Getenv("DYNAMODB_ENDPOINT"),
					SigningRegion: os.Getenv("AWS_REGION"),
				}, nil
			},
		)),
	)
	if err != nil {
		return "", fmt.Errorf("config load: %w", err)
	}

	// 2. Crea cliente DynamoDB
	client := dynamodb.NewFromConfig(cfg)

	// Nombre de la tabla, pásalo en la variable de entorno DYNAMODB_TABLE_NAME
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		return "", fmt.Errorf("please set DYNAMODB_TABLE_NAME env var")
	}

	// 3. Pone un ítem en la tabla
	item := map[string]types.AttributeValue{
		"id":      &types.AttributeValueMemberS{Value: "1"},
		"message": &types.AttributeValueMemberS{Value: "¡Hola desde DynamoDB!"},
	}
	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		return "", fmt.Errorf("put item: %w", err)
	}

	// 4. Recupera ese mismo ítem
	getOut, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: "1"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("get item: %w", err)
	}

	// 5. Extrae atributos y los devuelve
	idAttr, _ := getOut.Item["id"].(*types.AttributeValueMemberS)
	msgAttr, _ := getOut.Item["message"].(*types.AttributeValueMemberS)

	return fmt.Sprintf("Item recuperado → id=%s, message=%s", idAttr.Value, msgAttr.Value), nil
}

func main() {
	lambda.Start(handler)
}
