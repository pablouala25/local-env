## 📦 Ciclo de desarrollo

1. **Inicializar el entorno** (solo la primera vez)

   ```bash
   make dev
   ```
2. **Recargar cambios** (tras cada modificación de código)

   ```bash
   make reload
   ```

---

## 🔧 Requisitos previos

Antes de arrancar tus Lambdas, asegúrate de:

1. **Crear la tabla en DynamoDB Local**

   ```bash
      aws dynamodb create-table \
      --table-name Items \
      --attribute-definitions AttributeName=id,AttributeType=S \
      --key-schema AttributeName=id,KeyType=HASH \
      --billing-mode PAY_PER_REQUEST \
      --endpoint-url http://localhost:8000 \
      --region us-east-1
   ```

2. **Definir las variables de entorno** (archivo `.env` o `env.json`):

   ```dotenv
   AWS_REGION=us-east-1
   DYNAMODB_ENDPOINT=http://localhost:8000   # o http://host.docker.internal:8000
   DYNAMODB_TABLE_NAME=Items
   ```

3. **Verificar con curl**

   ```bash
   curl -X POST \
     -H "Content-Type: application/json" \
     http://localhost:9000/2015-03-31/functions/function/invocations \
     -d '{}'
   ```

   Deberías recibir:

   ```
   "Item recuperado → id=1, message=¡Hola desde DynamoDB!"
   ```









   Localstack


   Para exponer un **GET** en tu API Gateway (LocalStack o AWS real) debes hacer dos cosas:

1. **Configurar el método GET en API Gateway** (igual que con POST, pero cambiando el http-method).
2. **Adaptar tu handler Go** para leer `req.HTTPMethod == "GET"` y, si quieres, procesar parámetros de consulta.

---

## 1. Añadir GET en el bootstrap de LocalStack

En tu script `01-init.sh`, justo después de crear el recurso `/echo`, añade los pasos para GET:

```bash
# … después de crear POST …

# 4.b) Configurar el método GET
awslocal apigateway put-method \
  --rest-api-id "$API_ID" \
  --resource-id "$ECHO_ID" \
  --http-method GET \
  --authorization-type NONE \
  --region ${AWS_REGION}

# 5.b) Integración proxy GET → Lambda
awslocal apigateway put-integration \
  --rest-api-id "$API_ID" \
  --resource-id "$ECHO_ID" \
  --http-method GET \
  --type AWS_PROXY \
  --integration-http-method POST \
  --uri arn:aws:apigateway:${AWS_REGION}:lambda:path/2015-03-31/functions/arn:aws:lambda:${AWS_REGION}:000000000000:function:lambda-go-dev/invocations \
  --region ${AWS_REGION}
```

Luego vuelve a desplegar:

```bash
awslocal apigateway create-deployment \
  --rest-api-id "$API_ID" \
  --stage-name dev \
  --region ${AWS_REGION}
```

Ahora tu ruta `/echo` responderá tanto a **POST** como a **GET**.

---

## 2. Adaptar tu handler Go para GET

Si usas el handler con proxy (events.APIGatewayProxyRequest), detecta el método y lee parámetros de consulta:

```go
import (
  "github.com/aws/aws-lambda-go/events"
  // …
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  switch req.HTTPMethod {
  case "POST":
    // tu lógica actual de insertar/leer Dynamo…
    // body := req.Body
  case "GET":
    // ejemplar: leer ?id=123 de la URL
    id := req.QueryStringParameters["id"]
    if id == "" {
      return events.APIGatewayProxyResponse{StatusCode: 400, Body: "missing id"}, nil
    }
    // lee Dynamo con ese id
    out, err := client.GetItem(ctx, &dynamodb.GetItemInput{
      TableName: aws.String(table),
      Key: map[string]types.AttributeValue{
        "id": &types.AttributeValueMemberS{Value: id},
      },
    })
    if err != nil {
      return events.APIGatewayProxyResponse{StatusCode:500}, err
    }
    // serializa y devuelve JSON
    resp, _ := json.Marshal(map[string]string{
      "id":   id,
      "message": out.Item["message"].(*types.AttributeValueMemberS).Value,
    })
    return events.APIGatewayProxyResponse{
      StatusCode: 200,
      Headers:    map[string]string{"Content-Type":"application/json"},
      Body:       string(resp),
    }, nil

  default:
    return events.APIGatewayProxyResponse{
      StatusCode: 405,
      Body:       fmt.Sprintf("method %s not allowed", req.HTTPMethod),
    }, nil
  }
}
```

Con esto:

* **GET /echo?id=1** invoca el mismo Lambda, entra en el caso `"GET"`, lee el parámetro `id` y retorna el ítem de Dynamo.
* **POST /echo** sigue operando tal como ya lo tienes implementado.

---

## 3. Probar el GET localmente

Una vez redeployado el API en LocalStack, prueba con:

```bash
curl "http://localhost:4566/restapis/${API_ID}/dev/_user_request_/echo?id=1"
```

Y deberías recibir tu JSON con `{ "id": "1", "message": "¡Hola desde DynamoDB!" }`.

¡Así tendrás un endpoint GET funcionando exactamente igual que un servidor HTTP tradicional, pero dentro de tu Lambda en LocalStack!


awslocal apigateway get-stages --rest-api-id <API_ID>

------------------------------------------------

curl -i -X POST \
  http://localhost:4566/restapis/p4njryy8fp/dev/_user_request_/items \
  -H "Content-Type: application/json" \
  -d '{"id":"123","message":"¡hola mundo!"}'

HTTP/1.1 201 CREATED
Server: TwistedWeb/24.3.0
Date: Sun, 27 Jul 2025 15:07:02 GMT
Content-Type: application/json
Connection: keep-alive
Content-Length: 31
x-amzn-RequestId: f247ced3-a014-43cb-b5d9-580d34d438f7
x-amz-apigw-id: 218aee53=
X-Amzn-Trace-Id: Root=1-68864096-8afd73c595024e1956dde23f;Parent=1a766eee5ecb617e;Sampled=0
x-localstack: true

{"id":"123","status":"created"}%     

probar: URL de invocación de API Gateway (LocalStack ≥ 3.8):
/restapis/<api_id>/<stage>/_user_request_ está deprecated. Debes usar
/_aws/execute-api/<api_id>/<stage>.

------------------------------


❯ curl -i -X GET \
  "http://localhost:4566/restapis/p4njryy8fp/dev/_user_request_/items/123"

HTTP/1.1 200 OK
Server: TwistedWeb/24.3.0
Date: Sun, 27 Jul 2025 15:11:20 GMT
Content-Type: application/json
Connection: keep-alive
Content-Length: 38
x-amzn-RequestId: 8c020306-9f98-477e-9891-6fce08a156d2
x-amz-apigw-id: c75264e1=
X-Amzn-Trace-Id: Root=1-68864198-4686755b7f7b4f57dbc921cc;Parent=9b347d43897536a0;Sampled=0
x-localstack: true

{"id":"123","message":"¡hola mundo!"}%   