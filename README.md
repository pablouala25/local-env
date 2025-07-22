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
      --table-name MiTabla \
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
   DYNAMODB_TABLE_NAME=MiTabla
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