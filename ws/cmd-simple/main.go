package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) (string, error) {
	fmt.Println("¡Hola desde el handler Lambda!")
	return "ok2!!!", nil
}

func main() {
	lambda.Start(handler)
}
