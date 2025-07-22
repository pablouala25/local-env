## Resumen

A continuación encontrarás una secuencia de pasos simplificados y verificados para configurar tu entorno, crear un proyecto SAM, compilar correctamente el ejecutable `bootstrap` (cuando uses el runtime **provided.al2023**) y ejecutar tu función Lambda de forma local. Cada comando ha sido validado contra su documentación oficial para garantizar que obtengas un flujo limpio y reproducible.

---

## 1. Verificar instalaciones previas

Antes de comenzar, asegúrate de que las herramientas necesarias estén instaladas y funcionando:

1. **Docker CLI**

   ```bash
   docker --version
   ```

   Muestra la versión de Docker instalada en tu sistema ([Docker Documentation][1]).

2. **Prueba de contenedor “Hello World”**

   ```bash
   docker run hello-world
   ```

   Descarga y ejecuta la imagen `hello-world` para validar tu instalación ([Docker Documentation][2]) ([A Docker Tutorial for Beginners][3]).

3. **AWS CLI**

   ```bash
   aws --version
   ```

   Verifica la versión del cliente de AWS en tu entorno ([AWS Documentation][4]).

4. **AWS SAM CLI**

   ```bash
   sam --version
   ```

   Comprueba la versión del AWS SAM CLI instalada ([AWS Documentation][5]).

---

## 2. Inicializar un nuevo proyecto SAM

Crea la estructura básica de tu aplicación serverless:

```bash
sam init
```

Este comando te guiará para seleccionar un runtime (p. ej. Go), un template y configuraciones iniciales ([AWS Documentation][6]).

---

## 3. Compilar el ejecutable para **provided.al2023**

Si optas por el runtime **provided.al2023**, debes producir manualmente un archivo llamado `bootstrap` en la raíz de tu función:

```bash
cd hello-world
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/bootstrap main.go
```

Lambda buscará **siempre** el binario `bootstrap` en `/var/task`, y fallará con `InvalidEntrypoint` si no lo encuentra ([Stack Overflow][7]) ([Amazon Web Services, Inc.][8]).

---

## 4. Construir y ejecutar localmente

Una vez tengas el proyecto (y, si aplica, el `bootstrap`), construye y arranca tu API local:

1. **Construcción**

   ```bash
   sam build
   ```

   Empaqueta y prepara tu función según el runtime y el `CodeUri` configurado ([AWS Documentation][9]).

2. **Inicio de API local**

   ```bash
   sam local start-api
   ```

   Levanta un servidor HTTP en `http://127.0.0.1:3000` que enruta hacia tu función Lambda ([AWS Documentation][10]).

---

## 5. Probar el endpoint

Envía una petición para verificar que todo responde correctamente:

```bash
curl http://localhost:3000/hello
```

Deberías recibir tu saludo JSON o texto sin errores de “InvalidEntrypoint” ni de “no such file or directory”.


DYNAMO


## Resumen

Para conectar tu aplicación SAM local con DynamoDB Local corriendo en Docker Compose, debes:

1. **Levantar DynamoDB Local** dentro de una red Docker compartida.
2. **Inyectar el endpoint** de DynamoDB Local en tu función Lambda mediante variables de entorno.
3. **Arrancar SAM Local** usando la misma red y el archivo de variables de entorno (`--docker-network` y `--env-vars`).

Con estos pasos, tu función podrá comunicarse con la base de datos sin desplegar nada en AWS.

---

## 1. Levantar DynamoDB Local con Docker Compose

Usa el servicio `dynamodb-local` en tu `docker-compose.yml`, incluyendo la definición de red `app-network` (o la que prefieras):

```yaml
services:
  dynamodb-local:
    image: amazon/dynamodb-local:latest
    container_name: ${DYNAMODB_CONTAINER_NAME}
    ports:
      - "${DYNAMODB_PORT}:${DYNAMODB_PORT}"
    env_file:
      - .env
    command: ["-jar", "DynamoDBLocal.jar", "-sharedDb", "-port", "${DYNAMODB_PORT}"]
    volumes:
      - dynamodb_data:${DYNAMODB_DATA_PATH}
    networks:
      - ${APP_NETWORK}
    healthcheck:
      test: ["CMD-SHELL", "curl -s http://localhost:${DYNAMODB_PORT} || exit 1"]
      interval: 5s
      timeout: 5s
      retries: 5
      start_period: 10s
networks:
  app-network:
    driver: bridge
volumes:
  dynamodb_data:
```

* Docker Compose crea automáticamente la red `app-network` para los servicios ([Medium][1]).
* DynamoDB Local escuchará en `http://dynamodb-local:${DYNAMODB_PORT}` dentro de la red.

---

## 2. Exponer el endpoint a tu función Lambda

### 2.1 Definir variable en `template.yaml`

En tu plantilla SAM, añade el endpoint de DynamoDB Local en el bloque `Environment` de tu función:

```yaml
Resources:
  HelloWorldFunction:
    Type: AWS::Serverless::Function
    Properties:
      # ...
      Environment:
        Variables:
          DYNAMODB_ENDPOINT: "http://dynamodb-local:${DYNAMODB_PORT}"
          TABLE_NAME: !Ref MyTable
```

SAM Local solo inyecta variables declaradas aquí o pasadas vía `--env-vars` .

### 2.2 Crear `env.json` para parámetros dinámicos

En la raíz del proyecto, define un archivo `env.json`:

```json
{
  "HelloWorldFunction": {
    "DYNAMODB_ENDPOINT": "http://dynamodb-local:8000",
    "TABLE_NAME": "YourTableName"
  }
}
```

* Este JSON sobreescribe las variables de entorno en runtime ([AWS Documentation][2]).

---

## 3. Arrancar SAM Local conectado a la red Docker

Usa las opciones `--docker-network` y `--env-vars` de SAM CLI:

```bash
sam build                                                   # Empaqueta tu app SAM :contentReference[oaicite:3]{index=3}
sam local start-api \
  --docker-network app-network \
  --env-vars env.json
```

* `--docker-network app-network`: conecta el contenedor Lambda a la red donde corre DynamoDB Local ([AWS Documentation][2]).
* `--env-vars env.json`: inyecta `DYNAMODB_ENDPOINT` y `TABLE_NAME` en tu función ([AWS Documentation][2]).
* Verás en consola:

  ```
  Mounting HelloWorldFunction at http://127.0.0.1:3000/hello [POST,OPTIONS]
  ```

---

## 4. Probar la integración

Con la API en marcha y DynamoDB Local levantado:

```bash
curl -X POST http://127.0.0.1:3000/hello \
  -H "Content-Type: application/json" \
  -d '{"id":"123","message":"hola Dynamo"}'
```

Tu función ahora usará la variable `DYNAMODB_ENDPOINT` para conectar a `dynamodb-local:8000` dentro de la red Docker y escribirá en la tabla local ([GitHub][3]).

---

## Referencias

* `sam local start-api` opciones `--docker-network` y `--env-vars` ([AWS Documentation][2])
* Introducción a `sam local start-api` en AWS SAM CLI ([AWS Documentation][4])
* Problema de red Docker con SAM Local (StackOverflow) ([Stack Overflow][5])
* `--docker-network` explicado en Fig.io ([Fig][6])
* Soporte de red Docker en georgmao/aws-sam-local ([GitHub][3])
* Cómo ejecutar API Gateway, Lambda y DynamoDB Local juntos ([Medium][1])
* Acceso a recursos locales con SAM API Gateway ([blowstack.com][7])
* Creación de red Docker para recursos locales ([Medium][1])
* Uso de `DynamoDBCrudPolicy` para permisos en SAM ([Medium][8])
* Configuración de endpoint personalizado en AWS SDK Go v2 ([blowstack.com][9])

[1]: https://andmoredev.medium.com/how-to-run-api-gateway-aws-lambda-and-dynamodb-locally-91b75d9a54fe?utm_source=chatgpt.com "How to run API Gateway, AWS Lambda and DynamoDB locally"
[2]: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-local-start-api.html?utm_source=chatgpt.com "sam local start-api - AWS Serverless Application Model"
[3]: https://github.com/georgmao/aws-sam-local?utm_source=chatgpt.com "georgmao/aws-sam-local - GitHub"
[4]: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/using-sam-cli-local-start-api.html?utm_source=chatgpt.com "Introduction to testing with sam local start-api - AWS Documentation"
[5]: https://stackoverflow.com/questions/79216214/sam-local-api-docker-network-issue?utm_source=chatgpt.com "SAM Local Api Docker Network Issue - Stack Overflow"
[6]: https://fig.io/manual/sam/local/start-api?utm_source=chatgpt.com "sam local start-api - Fig.io"
[7]: https://blowstack.com/blog/how-to-access-local-resources-within-an-api-gateway-run-on-aws-sam?utm_source=chatgpt.com "How to access local resources within an API Gateway run on AWS ..."
[8]: https://chris-hart.medium.com/aws-sam-in-a-docker-container-arcadian-cloud-d69c3781ccfb?source=user_profile---------2----------------------------&utm_source=chatgpt.com "AWS SAM in a Docker Container — Arcadian Cloud | by Chris Hart"
[9]: https://blowstack.com/blog/how-to-run-multiple-api-gateway-instances-locally-using-sam?utm_source=chatgpt.com "How to run multiple API Gateway instances locally using SAM"
