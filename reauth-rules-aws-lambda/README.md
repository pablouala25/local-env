# README

Este proyecto emula un entorno AWS completo en local, combinando contenedores Docker y SAM CLI. A continuación encontrarás:

* 🛠️ **Cómo levantar el entorno**
* 📦 **Qué hace cada contenedor**
* 🔄 **Flujo de arranque**
* ⚙️ **Descripciones de archivos clave**

---

## Requisitos

* Docker & Docker Compose (v2+)
* Go (v1.18+)
* AWS SAM CLI
* AWS CLI (opcional para debugging)

---

## Estructura de archivos

```text
├── bin/                 # Binarios compilados (bootstrap)
├── cmd/                 # Código Go de la Lambda
├── docker-compose.yml   # Servicios Docker para DynamoDB, LocalStack y Lambda
├── template.yaml        # SAM template (CloudFormation/SAM)
├── Makefile             # Tareas para build, despliegue local y teardown
└── README.md            # Esta documentación
```

---

## 1. Compilar la Lambda Go

```bash
make build
```

* Genera `bin/bootstrap` (Go Linux/amd64, sin CGO).
* Es el ejecutable que correrá dentro del contenedor Lambda.

---

## 2. Levantar servicios Docker

```bash
make compose-up
```

Inicia en background:

1. **dynamodb-local**: emula DynamoDB en `http://dynamodb-local:8000`.
2. **dynamodb-admin**: UI web para explorar tablas DynamoDB (`http://localhost:8001`).
3. **lambda-go-dev**: contenedor con tu función Go.
4. **localstack**: emula servicios AWS secundarios (S3, SQS, SNS, etc.).
5. **sam-local**: emula API Gateway y enruta HTTP a la Lambda.

---

## 3. Iniciar API Gateway local (SAM CLI)

```bash
make sam-up
```

* Ejecuta `sam local start-api --host 0.0.0.0 --port 3000 --docker-network app-network`.
* Escucha peticiones HTTP y las envía al contenedor `lambda-go-dev`.

---

## 4. Arranque todo en un solo paso

```bash
make all
```

Equivale a:

1. `make build`
2. `make compose-up`
3. `make sam-up`

Al finalizar verás:

```
✅ Development environment ready! Visit http://localhost:3000
```

---

## 5. Probar la API

### Prueba Rápida
```bash
# Ejecutar script de prueba completa
./scripts/quick-test.sh

# O verificar contenedores
make check-containers
```

### Documentación de la API
- 📚 **[Documentación Completa](API_DOCUMENTATION.md)** - Guía detallada con ejemplos
- 🚀 **[Ejemplos Rápidos](CURL_EXAMPLES.md)** - Comandos curl listos para usar

### Endpoints Disponibles
```
POST   /api/v1/                    # Crear regla
GET    /api/v1/{id}                # Obtener regla por ID
GET    /api/v1/rules               # Listar reglas
GET    /api/v1/rules/by-initiative # Obtener reglas por iniciativa
PUT    /api/v1/{id}                # Actualizar regla
DELETE /api/v1/{id}                # Eliminar regla
```

### Ejemplo de Uso
```bash
# Crear una regla
curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "POST",
    "path": "/api/v1/",
    "body": "{\"initiative_id\": \"test\", \"category\": \"auth\", \"description\": \"Mi regla\", \"rule\": \"user.age >= 18\"}",
    "headers": {"Content-Type": "application/json"}
  }'
```

Deberías recibir un HTTP 200 y ver la regla guardada en DynamoDB Local.

---

## Descripción de los contenedores

### 1. DynamoDB Local

* Imagen: `amazon/dynamodb-local:latest`.
* Expone DynamoDB en `http://dynamodb-local:${DYNAMODB_PORT}`.
* Se monta el volumen `dynamodb_data` para persistencia.
* Healthcheck con CURL para verificar que arranque.

### 2. DynamoDB Admin

* Imagen: `aaronshaf/dynamodb-admin`.
* UI web en `http://localhost:${DYNAMODB_ADMIN_PORT}`.
* Conecta contra `dynamodb-local`.

### 3. Lambda Go Dev

* Imagen: `public.ecr.aws/lambda/provided:al2023`.
* Monta tu binario `bootstrap` como runtime `/var/runtime/bootstrap`.
* Recibe envs: `AWS_REGION`, `DYNAMODB_ENDPOINT`, `DYNAMODB_TABLE_NAME`.
* Actúa como tu función Lambda real.

### 4. LocalStack

* Imagen: `localstack/localstack:latest`.
* Emula múltiples servicios AWS (S3, SQS, SNS, etc.)
* Sin Lambda ni DynamoDB (dejamos esos servicios fuera para evitar duplicidad).

### 5. SAM Local (API Gateway)

* Contenedor que ejecuta `sam local start-api`.
* Emula API Gateway y enruta HTTP a la función Lambda.
* Usa tu `template.yaml` para leer rutas, métodos y variables de entorno.

---

## Diagrama de interacción

```text
[you / host]
   │
   │  make all
   ▼

┌─────────────────┐     ┌───────────────┐     ┌───────────────┐
│ sam-local       │◀───▶│ lambda-go-dev │◀───┐│ dynamodb-local│
│ (API Gateway)   │     │ (Go Lambda)   │     │ (DynamoDB)    │
└─────────────────┘     └───────────────┘     └───────────────┘
       │
       ▼
┌─────────────────┐
│ localstack      │
│ (otros servicios│
│  AWS emulados)  │
└─────────────────┘
```

---

## Detalles adicionales

* El `template.yaml` define la función, la tabla DynamoDB y las variables de entorno.
* Al usar `!Ref Items` en SAM, `DYNAMODB_TABLE_NAME` se resuelve a `Items`.
* Puedes expandir `template.yaml` para agregar más funciones, triggers (SQS, SNS, EventBridge), buckets S3, etc.
* Para una integración end-to-end con SQS, activa `sqs` en `SERVICES` de LocalStack y declara el EventSourceMapping en `template.yaml`.

---

¡Listo! Ahora tienes un entorno local 100% reproducible, versionable y alineado con AWS.
