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
	// 1. Carga solo la región en config (sin resolver endpoints aquí)
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		return "", fmt.Errorf("config load: %w", err)
	}

	// 2. Crea el cliente e **añade** el endpoint local con BaseEndpoint
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("DYNAMODB_ENDPOINT"))
	}) // :contentReference[oaicite:0]{index=0}

	// 3. Nombre de tabla desde env
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		return "", fmt.Errorf("please set DYNAMODB_TABLE_NAME env var")
	}

	// 4. Inserta un ítem
	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"id":      &types.AttributeValueMemberS{Value: "1"},
			"message": &types.AttributeValueMemberS{Value: "¡Hola desde DynamoDB!!!"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("put item: %w", err)
	}

	// 5. Recupera ese mismo ítem
	getOut, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: "1"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("get item: %w", err)
	}

	// 6. Extrae atributos y los devuelve
	idAttr := getOut.Item["id"].(*types.AttributeValueMemberS)
	msgAttr := getOut.Item["message"].(*types.AttributeValueMemberS)

	return fmt.Sprintf("Item recuperado → id=%s, message=%s", idAttr.Value, msgAttr.Value), nil
}

func main() {
	lambda.Start(handler)
}
