## ✅ **T Ciclo:**

1. **Levantar entorno la primera vez:**

   ```sh
   make dev
   ```
2. **Después, por cada cambio de código:**

   ```sh
   make reload
   ```



---

### Pasos previos

1. **Crear la tabla en DynamoDB Local**
   Ejecuta (ajusta región/endpoint si cambia):

   ```bash
   aws dynamodb create-table \
     --table-name MiTabla \
     --attribute-definitions AttributeName=id,AttributeType=S \
     --key-schema AttributeName=id,KeyType=HASH \
     --billing-mode PAY_PER_REQUEST \
     --endpoint-url http://localhost:8000 \
     --region us-east-1
   ```

2. **Variables de entorno** (en tu `.env` o en la configuración de SAM/Toolkit):

   ```
   AWS_REGION=us-east-1
   DYNAMODB_ENDPOINT=http://host.docker.internal:8000    # o http://localhost:8000
   DYNAMODB_TABLE_NAME=MiTabla
   ```

3. **Construir y levantar localmente**

   * Con SAM CLI:

     ```bash
     sam local start-api \
       --env-vars env.json \
       --docker-network tu_red_docker
     ```
   * O usando tu herramienta preferida (Toolkit de VSCode, `make dev`, etc.)

4. **Probar con curl**

   ```bash
   curl -X POST \
     -H "Content-Type: application/json" \
     http://localhost:3000/2015-03-31/functions/function/invocations \
     -d '{}'
   ```

   Deberías ver algo como:

   ```
   "Item recuperado → id=1, message=¡Hola desde DynamoDB!"
   ```
