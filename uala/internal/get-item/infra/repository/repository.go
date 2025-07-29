package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	shareaderrors "github.com/Bancar/reauth-bff-aws-lambda/internal/shared/errors"

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

func (r *Repository) GetItem(ctx context.Context, id string) (*shareddomain.Item, error) {
	out, err := r.client.GetItem(ctx, &ddb.GetItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, shareaderrors.ErrNotFound
	}

	it := sharedmodels.Item{
		ID:      getString(out.Item["id"]),
		Message: getString(out.Item["message"]),
	}
	return it.ToDomain(), nil
}

func getString(av types.AttributeValue) string {
	if s, ok := av.(*types.AttributeValueMemberS); ok {
		return s.Value
	}
	return ""
}
