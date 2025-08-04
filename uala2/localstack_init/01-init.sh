#!/usr/bin/env bash
set -euo pipefail

echo "🚀 Bootstrapping CRUD API + SQS + Lambda in LocalStack..."

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
ZIP_PATH="/opt/code/bootstrap.zip"

# Nuevo: base path (api/v1)
TRIM() { sed -e 's#^/##' -e 's#/$##'; }
BASE_PATH_RAW="${BASE_PATH:-api/v1}"
BASE_PATH_TRIMMED="$(echo "$BASE_PATH_RAW" | TRIM)"   # ej: api/v1

###############################################
# 0) DynamoDB - ensure table exists at $ENDPOINT
###############################################
echo "- Ensuring DynamoDB table '$TABLE' at endpoint: $ENDPOINT"
if aws dynamodb describe-table \
      --table-name "$TABLE" \
      --endpoint-url "$ENDPOINT" \
      --region "$REGION" >/dev/null 2>&1; then
  echo "  • Table '$TABLE' already exists"
else
  echo "  • Creating table '$TABLE' (PAY_PER_REQUEST)"
  aws dynamodb create-table \
    --table-name "$TABLE" \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url "$ENDPOINT" \
    --region "$REGION" >/dev/null
  echo "  • Waiting for table to be ACTIVE…"
  aws dynamodb wait table-exists \
    --table-name "$TABLE" \
    --endpoint-url "$ENDPOINT" \
    --region "$REGION"
  echo "  • Table '$TABLE' is ready"
fi
###############################################

# 1) Lambda
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

# 2) SQS
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

# 3) API Gateway
API_ID=$(awslocal apigateway get-rest-apis --query "items[?name=='$API'].id" --output text)
if [ -z "$API_ID" ] || [ "$API_ID" = "None" ]; then
  echo "- Creating API Gateway '$API'"
  API_ID=$(awslocal apigateway create-rest-api --name "$API" --region "$REGION" --query 'id' --output text)
else
  echo "- API Gateway '$API' exists (ID: $API_ID)"
fi

ROOT_ID=$(awslocal apigateway get-resources --rest-api-id "$API_ID" --region "$REGION" --query 'items[0].id' --output text)

# 3.a) Crear recursos anidados del base path: /api/v1
PARENT_ID="$ROOT_ID"
CURRENT_PATH=""
IFS='/' read -r -a SEGMENTS <<< "$BASE_PATH_TRIMMED"
for seg in "${SEGMENTS[@]}"; do
  [ -z "$seg" ] && continue
  CANDIDATE_ID=$(awslocal apigateway create-resource \
      --rest-api-id "$API_ID" \
      --parent-id "$PARENT_ID" \
      --path-part "$seg" \
      --region "$REGION" \
      --query 'id' --output text 2>/dev/null || true)

  CURRENT_PATH="${CURRENT_PATH}/${seg}"
  if [ -z "$CANDIDATE_ID" ] || [ "$CANDIDATE_ID" = "None" ]; then
    CANDIDATE_ID=$(awslocal apigateway get-resources \
      --rest-api-id "$API_ID" \
      --region "$REGION" \
      --query "items[?path=='$CURRENT_PATH'].id" \
      --output text)
    echo "- Resource '$CURRENT_PATH' exists"
  else
    echo "- Created '$CURRENT_PATH' (ID: $CANDIDATE_ID)"
  fi
  PARENT_ID="$CANDIDATE_ID"
done
BASE_ID="$PARENT_ID"   # ID de /api/v1

# 3.b) Crear /items y /items/{id} bajo /api/v1
ITEMS_ID=$(awslocal apigateway create-resource --rest-api-id "$API_ID" --parent-id "$BASE_ID" --path-part "$RESOURCE" --region "$REGION" --query 'id' --output text 2>/dev/null || true)
if [ -z "$ITEMS_ID" ] || [ "$ITEMS_ID" = "None" ]; then
  ITEMS_ID=$(awslocal apigateway get-resources --rest-api-id "$API_ID" --region "$REGION" --query "items[?path=='/${BASE_PATH_TRIMMED}/$RESOURCE'].id" --output text)
  echo "- Resource '/${BASE_PATH_TRIMMED}/$RESOURCE' exists"
else
  echo "- Created '/${BASE_PATH_TRIMMED}/$RESOURCE' (ID: $ITEMS_ID)"
fi

ITEM_ID_ID=$(awslocal apigateway create-resource --rest-api-id "$API_ID" --parent-id "$ITEMS_ID" --path-part "{id}" --region "$REGION" --query 'id' --output text 2>/dev/null || true)
if [ -z "$ITEM_ID_ID" ] || [ "$ITEM_ID_ID" = "None" ]; then
  ITEM_ID_ID=$(awslocal apigateway get-resources --rest-api-id "$API_ID" --region "$REGION" --query "items[?path=='/${BASE_PATH_TRIMMED}/$RESOURCE/{id}'].id" --output text)
  echo "- Resource '/${BASE_PATH_TRIMMED}/$RESOURCE/{id}' exists"
else
  echo "- Created '/${BASE_PATH_TRIMMED}/$RESOURCE/{id}' (ID: $ITEM_ID_ID)"
fi

# 4) Integraciones + permisos
add_proxy() {
  local res_id=$1 method=$2 func=$3 res_path=$4
  echo "   • Configuring [$method] → Lambda '$func' on /$res_path"
  awslocal apigateway put-method --rest-api-id "$API_ID" --resource-id "$res_id" --http-method "$method" --authorization-type NONE --region "$REGION"
  awslocal apigateway put-integration --rest-api-id "$API_ID" --resource-id "$res_id" --http-method "$method" --type AWS_PROXY --integration-http-method POST --uri "arn:aws:apigateway:$REGION:lambda:path/2015-03-31/functions/arn:aws:lambda:$REGION:000000000000:function:$func/invocations" --region "$REGION"
  # Permiso amplio: cualquier método y stage para esa ruta
  awslocal lambda add-permission \
    --function-name "$func" \
    --statement-id "apigw-${func}-${method}-${res_id}" \
    --action lambda:InvokeFunction \
    --principal apigateway.amazonaws.com \
    --source-arn "arn:aws:execute-api:$REGION:000000000000:$API_ID/*/*/$res_path" \
    --region "$REGION" 2>/dev/null || true
}

FULL_ITEM_PATH="${BASE_PATH_TRIMMED}/${RESOURCE}"   # ej: api/v1/items
echo "- Configuring CRUD methods under '/${FULL_ITEM_PATH}'"
add_proxy "$ITEMS_ID"   POST   "$FUNC_CRUD" "$FULL_ITEM_PATH"
add_proxy "$ITEM_ID_ID" GET    "$FUNC_CRUD" "${FULL_ITEM_PATH}/*"
add_proxy "$ITEM_ID_ID" PUT    "$FUNC_CRUD" "${FULL_ITEM_PATH}/*"
add_proxy "$ITEM_ID_ID" DELETE "$FUNC_CRUD" "${FULL_ITEM_PATH}/*"

# 5) Deploy
echo "- Deploying to stage 'dev'"
awslocal apigateway create-deployment --rest-api-id "$API_ID" --stage-name dev --region "$REGION" >/dev/null || true

# 6) Mostrar URLs y cURL de prueba (estilo nuevo de LocalStack)
API_BASE="http://localhost:4566/_aws/execute-api/${API_ID}/dev/${BASE_PATH_TRIMMED}"
API_URL="${API_BASE}/${RESOURCE}"

echo "✅ Initialization complete!"
echo "  • API (new style): ${API_BASE}"
echo "  • POST items:     ${API_URL}"
echo ""
echo "➡️  Test POST (1 línea):"
echo "curl -sS -i -X POST \"${API_URL}\" -H 'Content-Type: application/json' -d '{\"id\":\"123\",\"message\":\"¡hola mundo!\"}'"
echo ""
echo "➡️  Test GET (1 línea):"
echo "curl -sS -i \"${API_BASE}/${RESOURCE}/123\""
