// package bootstrap

// import (
// 	"context"
// 	"log/slog"
// 	"os"
// 	"strings"

// 	"github.com/Bancar/uala-auth-team-go-dependencies/libs/lambda"
// )

// func SetupLambda() *lambda.Lambda {
// 	ctx := context.Background()

// 	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
// 	logger.Info("bootstrap started...")

// 	region := getenv("AWS_REGION", "us-east-1")
// 	table := os.Getenv("DYNAMODB_TABLE_NAME")
// 	if strings.TrimSpace(table) == "" {
// 		panic("DYNAMODB_TABLE_NAME not set")
// 	}

// 	cli := newDynamoClient(ctx, region, os.Getenv("DYNAMODB_ENDPOINT"))
// 	a := &app{db: cli, table: table}

// 	routes := createRoutes(a)

// 	settings := lambda.Settings{
// 		Metrics:      lambda.MetricsOptions{Namespace: "Demo/Items"},
// 		Tracing:      lambda.TracingOptions{Enabled: false},
// 		ErrorHandler: apierrors.NewAPIErrorHandlerV2(logger),
// 	}

// 	return lambda.MustSetupProxy(routes, settings)
// }

package bootstrap

import (
	"context"
	"log/slog"
	"os"

	"github.com/Bancar/reauth-bff-aws-lambda/pkg/logging"
	"github.com/Bancar/uala-auth-team-go-dependencies/libs/errors"

	"github.com/Bancar/uala-auth-team-go-dependencies/libs/lambda"

	"github.com/Bancar/reauth-bff-aws-lambda/cmd/config"
)

const metricsNamespace = "ReauthBFF"

func SetupLambda() *lambda.Lambda {
	ctx := context.Background()
	env := config.LoadConfig()

	logger := lambda.Logger(os.Stdout, logging.ContextAttrs())
	logger.InfoContext(ctx, "bootstrap started...", slog.Any("env_config", env))

	lambdaSettings := lambda.Settings{
		Metrics:      lambda.MetricsOptions{Namespace: metricsNamespace},
		Tracing:      lambda.TracingOptions{Enabled: false},
		ErrorHandler: errors.NewAPIErrorHandlerV2(logger),
	}

	routes := routes(ctx, env)

	return lambda.MustSetupProxy(routes, lambdaSettings)
}
