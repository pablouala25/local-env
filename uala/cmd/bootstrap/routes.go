package bootstrap

import (
	"context"
	"net/http"

	"github.com/Bancar/uala-auth-team-go-dependencies/libs/lambda"

	dynamodb "github.com/Bancar/reauth-bff-aws-lambda/pkg/aws/dynamodb"

	config "github.com/Bancar/reauth-bff-aws-lambda/cmd/config"
	createhandler "github.com/Bancar/reauth-bff-aws-lambda/internal/create-item/infra/handler"
	createrepository "github.com/Bancar/reauth-bff-aws-lambda/internal/create-item/infra/repository"
	createusecase "github.com/Bancar/reauth-bff-aws-lambda/internal/create-item/usecase"
	getitemhandler "github.com/Bancar/reauth-bff-aws-lambda/internal/get-item/infra/handler"
	getitemrepository "github.com/Bancar/reauth-bff-aws-lambda/internal/get-item/infra/repository"
	getitemusecase "github.com/Bancar/reauth-bff-aws-lambda/internal/get-item/usecase"
)

func routes(ctx context.Context, env *config.Config) lambda.Routes {

	// Cliente común de DynamoDB
	dct := dynamodb.NewDynamoClient(ctx, env.DynamoTableName, env.AWSRegion)

	// create-item
	createRepo := createrepository.NewRepository(dct, env.DynamoTableName)
	createUC := createusecase.NewUseCase(createRepo)
	createHandler := createhandler.NewHandler(createUC)

	// get-item
	getItemRepo := getitemrepository.NewRepository(dct, env.DynamoTableName)
	getItemUC := getitemusecase.NewUseCase(getItemRepo)
	getItemHandler := getitemhandler.NewHandler(getItemUC)

	return lambda.Routes{
		{
			Path: "/api/v1",
			Endpoints: lambda.Endpoints{
				{
					Name:    "create-item",
					Method:  http.MethodPost,
					Path:    "/items",
					Handler: http.HandlerFunc(createHandler.CreateItem),
					Enabled: true,
				},
				{
					Name:    "get-item",
					Method:  http.MethodGet,
					Path:    "/items/{id}",
					Handler: http.HandlerFunc(getItemHandler.GetItem),
					Enabled: true,
				},
			},
		},
	}
}
