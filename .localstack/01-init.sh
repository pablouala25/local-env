#!/usr/bin/env bash
set -euo pipefail

echo "🚀 Bootstrapping: CRUD API + SQS producer/consumer en LocalStack..."

# — Parámetros configurables (puedes exportarlos o fijarlos aquí) —
REGION=${AWS_REGION:-us-east-1}
QUEUE_NAME=${SQS_QUEUE_NAME:-items-queue}
API_NAME=${API_NAME:-items-api}
RESOURCE_NAME=${RESOURCE_NAME:-items}
LAMBDA_CRUD=${LAMBDA_CRUD:-lambda-crud}            # Función que maneja CRUD de items
LAMBDA_PRODUCER=${LAMBDA_PRODUCER:-lambda-producer}# Función que publica en SQS
LAMBDA_CONSUMER=${LAMBDA_CONSUMER:-lambda-consumer}# Función que procesa mensajes SQS

# 1) Crear cola SQS
echo "- Creando cola SQS '${QUEUE_NAME}'"
awslocal sqs create-queue \
  --queue-name "${QUEUE_NAME}" \
  --region "${REGION}"

QUEUE_URL="http://localhost:4566/000000000000/${QUEUE_NAME}"
QUEUE_ARN=$(awslocal sqs get-queue-attributes \
  --queue-url "${QUEUE_URL}" \
  --attribute-names QueueArn \
  --region "${REGION}" \
  --query 'Attributes.QueueArn' --output text)

# 2) Conectar la cola al Lambda consumer
echo "- Mapeando SQS→Lambda (${LAMBDA_CONSUMER})"
awslocal lambda create-event-source-mapping \
  --function-name "${LAMBDA_CONSUMER}" \
  --batch-size 10 \
  --event-source-arn "${QUEUE_ARN}" \
  --region "${REGION}"

# 3) Crear REST API en API Gateway
echo "- Creando API Gateway '${API_NAME}'"
API_ID=$(awslocal apigateway create-rest-api \
  --name "${API_NAME}" \
  --region "${REGION}" \
  --query 'id' --output text)

# 4) Obtener el resource-id raíz (“/”)
ROOT_ID=$(awslocal apigateway get-resources \
  --rest-api-id "${API_ID}" \
  --region "${REGION}" \
  --query 'items[0].id' --output text)

# 5) Crear recursos: /items y /items/{id} y /publish
echo "- Creando recursos API"
ITEMS_ID=$(awslocal apigateway create-resource \
  --rest-api-id "${API_ID}" \
  --parent-id "${ROOT_ID}" \
  --path-part "${RESOURCE_NAME}" \
  --region "${REGION}" \
  --query 'id' --output text)

ITEM_ID_ID=$(awslocal apigateway create-resource \
  --rest-api-id "${API_ID}" \
  --parent-id "${ITEMS_ID}" \
  --path-part "{id}" \
  --region "${REGION}" \
  --query 'id' --output text)

PUBLISH_ID=$(awslocal apigateway create-resource \
  --rest-api-id "${API_ID}" \
  --parent-id "${ROOT_ID}" \
  --path-part publish \
  --region "${REGION}" \
  --query 'id' --output text)

# --- Función auxiliar para métodos proxy a Lambda ---
add_method() {
  local res_id=$1 method=$2 lambda=$3
  awslocal apigateway put-method \
    --rest-api-id "${API_ID}" \
    --resource-id "${res_id}" \
    --http-method "${method}" \
    --authorization-type NONE \
    --region "${REGION}"
  awslocal apigateway put-integration \
    --rest-api-id "${API_ID}" \
    --resource-id "${res_id}" \
    --http-method "${method}" \
    --type AWS_PROXY \
    --integration-http-method POST \
    --uri "arn:aws:apigateway:${REGION}:lambda:path/2015-03-31/functions/arn:aws:lambda:${REGION}:000000000000:function:${lambda}/invocations" \
    --region "${REGION}"
  # Permiso para API Gateway invocar la Lambda
  awslocal lambda add-permission \
    --function-name "${lambda}" \
    --statement-id "apigw-${lambda}-${method}-${res_id}" \
    --action "lambda:InvokeFunction" \
    --principal "apigateway.amazonaws.com" \
    --source-arn "arn:aws:execute-api:${REGION}:000000000000:${API_ID}/*/${method}/${RESOURCE_NAME}" \
    --region "${REGION}" >/dev/null 2>&1 || true
}

echo "- Configurando métodos CRUD (/items)"
add_method "${ITEMS_ID}" POST   "${LAMBDA_CRUD}"
add_method "${ITEM_ID_ID}" GET   "${LAMBDA_CRUD}"
add_method "${ITEM_ID_ID}" PUT   "${LAMBDA_CRUD}"
add_method "${ITEM_ID_ID}" DELETE "${LAMBDA_CRUD}"

echo "- Configurando método PUBLISH (/publish)"
add_method "${PUBLISH_ID}" POST  "${LAMBDA_PRODUCER}"

# 6) Desplegar al stage 'dev'
echo "- Desplegando API al stage 'dev'"
awslocal apigateway create-deployment \
  --rest-api-id "${API_ID}" \
  --stage-name dev \
  --region "${REGION}"

# 7) Información final
echo "✅ Inicialización completa!"
echo "  • SQS QueueURL: ${QUEUE_URL}"
echo "  • API base URL: http://localhost:4566/restapis/${API_ID}/dev/_user_request_"
echo "  • CRUD endpoints:"
echo "     POST   /${RESOURCE_NAME}"
echo "     GET    /${RESOURCE_NAME}/{id}"
echo "     PUT    /${RESOURCE_NAME}/{id}"
echo "     DELETE /${RESOURCE_NAME}/{id}"
echo "  • Publish endpoint:"
echo "     POST   /publish  → envía mensajes a SQS"
echo "  • Lambda consumer (${LAMBDA_CONSUMER}) suscrito a SQS"
