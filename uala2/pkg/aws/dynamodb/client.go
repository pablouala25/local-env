package dynamodb

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewDynamoClient(ctx context.Context, tableName string, region string) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		panic(err)
	}

	// 👇 Soporte para LocalStack
	if ep := os.Getenv("DYNAMODB_ENDPOINT"); ep != "" {
		return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(ep)
		})
	}

	if tableName == "" {
		panic(errors.New("dynamodb table name not set"))
	}

	return dynamodb.NewFromConfig(cfg)
}
