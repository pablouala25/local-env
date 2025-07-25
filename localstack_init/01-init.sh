#!/usr/bin/env bash
set -euo pipefail

# -----------------------------------------------------------------------------
# 01-init.sh
# Bootstraps CRUD Lambda, SQS queue and API Gateway in LocalStack.
# Idempotent: puede ejecutarse varias veces sin erroes.
# -----------------------------------------------------------------------------

echo "🚀 Bootstrapping CRUD API + SQS + Lambda in LocalStack..."

# Validación de variables del entorno
: "${AWS_REGION:?AWS_REGION not set}"
: "${DYNAMODB_TABLE_NAME:?DYNAMODB_TABLE_NAME not set}"
: "${DYNAMODB_ENDPOINT:?DYNAMODB_ENDPOINT not set}"
: "${LAMBDA_CRUD:?LAMBDA_CRUD not set}"
: "${SQS_QUEUE_NAME:?SQS_QUEUE_NAME not set}"
: "${API_NAME:?API_NAME not set}"
: "${RESOURCE_NAME:?RESOURCE_NAME not set}"

REGION="$AWS_REGION"
TABLE="$DYNAMODB_TABLE_NAME"
ENDPOINT="$DYNAMODB_ENDPOINT"
FUNC_CRUD="$LAMBDA_CRUD"
QUEUE="$SQS_QUEUE_NAME"
API="$API_NAME"
RESOURCE="$RESOURCE_NAME"
ZIP_PATH="/opt/code/bootstrap.zip"  # Debe montarse en docker-compose

# 1) Crear o actualizar Lambda CRUD
if ! awslocal lambda get-function --function-name "$FUNC_CRUD" >/dev/null 2>&1; then
  echo "- Creating Lambda function '$FUNC_CRUD'"
  awslocal lambda create-function \
    --function-name "$FUNC_CRUD" \
    --runtime provided.al2 \
    --handler bootstrap \
    --role arn:aws:iam::000000000000:role/lambda-role \
    --zip-file fileb://"$ZIP_PATH" \
    --environment "Variables={AWS_REGION=$REGION,DYNAMODB_ENDPOINT=$ENDPOINT,DYNAMODB_TABLE_NAME=$TABLE}"
else
  echo "- Lambda '$FUNC_CRUD' already exists, skipping"
fi

# 2) Crear o conseguir SQS queue
if ! QUEUE_URL=$(awslocal sqs get-queue-url --queue-name "$QUEUE" --region "$REGION" 2>/dev/null); then
  echo "- Creating SQS queue '$QUEUE'"
  QUEUE_URL=$(awslocal sqs create-queue --queue-name "$QUEUE" --region "$REGION" --query 'QueueUrl' --output text)
else
  echo "- SQS queue '$QUEUE' exists"
fi
QUEUE_ARN=$(awslocal sqs get-queue-attributes \
  --queue-url "$QUEUE_URL" \
  --attribute-names QueueArn \
  --region "$REGION" \
  --query 'Attributes.QueueArn' --output text)

# 3) Crear o conseguir API Gateway
API_ID=$(awslocal apigateway get-rest-apis --query "items[?name=='$API'].id" --output text)
if [ -z "$API_ID" ]; then
  echo "- Creating API Gateway '$API'"
  API_ID=$(awslocal apigateway create-rest-api --name "$API" --region "$REGION" --query 'id' --output text)
else
  echo "- API Gateway '$API' exists (ID: $API_ID)"
fi

# 4) Obtener root resource
ROOT_ID=$(awslocal apigateway get-resources --rest-api-id "$API_ID" --region "$REGION" --query 'items[0].id' --output text)

# 5) Crear recursos /<RESOURCE> y /<RESOURCE>/{id}
ITEMS_ID=$(awslocal apigateway create-resource --rest-api-id "$API_ID" --parent-id "$ROOT_ID" --path-part "$RESOURCE" --region "$REGION" --query 'id' --output text || echo "")
if [ -z "$ITEMS_ID" ]; then
  ITEMS_ID=$(awslocal apigateway get-resources --rest-api-id "$API_ID" --region "$REGION" --query "items[?path=='/$RESOURCE'].id" --output text)
  echo "- Resource '/$RESOURCE' exists"
else
  echo "- Created '/$RESOURCE' (ID: $ITEMS_ID)"
fi

ITEM_ID_ID=$(awslocal apigateway create-resource --rest-api-id "$API_ID" --parent-id "$ITEMS_ID" --path-part "{id}" --region "$REGION" --query 'id' --output text || echo "")
if [ -z "$ITEM_ID_ID" ]; then
  ITEM_ID_ID=$(awslocal apigateway get-resources --rest-api-id "$API_ID" --region "$REGION" --query "items[?path=='/$RESOURCE/{id}'].id" --output text)
  echo "- Resource '/$RESOURCE/{id}' exists"
else
  echo "- Created '/$RESOURCE/{id}' (ID: $ITEM_ID_ID)"
fi

# 6) Función auxiliar para configurar integración proxy
add_proxy() {
  local res_id=$1 method=$2 func=$3
  echo "   • Configuring [$method] → Lambda '$func'"
  awslocal apigateway put-method --rest-api-id "$API_ID" --resource-id "$res_id" --http-method "$method" --authorization-type NONE --region "$REGION"
  awslocal apigateway put-integration --rest-api-id "$API_ID" --resource-id "$res_id" --http-method "$method" --type AWS_PROXY --integration-http-method POST --uri "arn:aws:apigateway:$REGION:lambda:path/2015-03-31/functions/arn:aws:lambda:$REGION:000000000000:function:$func/invocations" --region "$REGION"
  awslocal lambda add-permission --function-name "$func" --statement-id "apigw-${func}-${method}-${res_id}" --action lambda:InvokeFunction --principal apigateway.amazonaws.com --source-arn "arn:aws:execute-api:$REGION:000000000000:$API_ID/*/*/$RESOURCE" --region "$REGION" 2>/dev/null || true
}

echo "- Configuring CRUD methods"
add_proxy "$ITEMS_ID" POST   "$FUNC_CRUD"
add_proxy "$ITEM_ID_ID" GET    "$FUNC_CRUD"
add_proxy "$ITEM_ID_ID" PUT    "$FUNC_CRUD"
add_proxy "$ITEM_ID_ID" DELETE "$FUNC_CRUD"

# 7) Deploy
echo "- Deploying to stage 'dev'"
awslocal apigateway create-deployment --rest-api-id "$API_ID" --stage-name dev --region "$REGION" >/dev/null

echo "✅ Initialization complete!"
echo "  • API URL: http://localhost:4566/restapis/$API_ID/dev/_user_request_/$RESOURCE"
