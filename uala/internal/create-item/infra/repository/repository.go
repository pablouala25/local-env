package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	shareddomain "github.com/Bancar/reauth-bff-aws-lambda/internal/shared/domain"
	sharedmodels "github.com/Bancar/reauth-bff-aws-lambda/internal/shared/models"
)

type Repository struct {
	client *ddb.Client
	table  string
}

func NewRepository(c *ddb.Client, table string) *Repository {
	return &Repository{client: c, table: table}
}

func (r *Repository) CreateItem(ctx context.Context, it *shareddomain.Item) error {
	m := sharedmodels.FromDomain(it)

	_, err := r.client.PutItem(ctx, &ddb.PutItemInput{
		TableName: aws.String(r.table),
		Item: map[string]types.AttributeValue{
			"id":      &types.AttributeValueMemberS{Value: m.ID},
			"message": &types.AttributeValueMemberS{Value: m.Message},
		},
	})
	return err
}
