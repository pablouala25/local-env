# 🚀 Ejemplos Rápidos de CURL - Rules API

Colección de comandos curl listos para usar con la API de Rules Management.

## 📋 Índice Rápido

- [Crear Regla](#crear-regla)
- [Obtener Regla](#obtener-regla)
- [Listar Reglas](#listar-reglas)
- [Filtrar por Iniciativa](#filtrar-por-iniciativa)
- [Actualizar Regla](#actualizar-regla)
- [Eliminar Regla](#eliminar-regla)

---

## ➕ Crear Regla

```bash
curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "POST",
    "path": "/api/v1/",
    "body": "{\"initiative_id\": \"mi-iniciativa\", \"category\": \"auth\", \"description\": \"Mi regla\", \"rule\": \"user.age >= 18\"}",
    "headers": {"Content-Type": "application/json"}
  }'
```

---

## 🔍 Obtener Regla

```bash
# Reemplaza RULE_ID con el ID real de la regla
curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "GET",
    "path": "/api/v1/RULE_ID",
    "headers": {"Content-Type": "application/json"}
  }'
```

---

## 📋 Listar Reglas

```bash
# Listar todas las reglas
curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "GET",
    "path": "/api/v1/rules",
    "headers": {"Content-Type": "application/json"}
  }'

# Con paginación (limit=5)
curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "GET",
    "path": "/api/v1/rules",
    "queryStringParameters": {"limit": "5"},
    "headers": {"Content-Type": "application/json"}
  }'
```

---

## 🎯 Filtrar por Iniciativa

```bash
curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "GET",
    "path": "/api/v1/rules/by-initiative",
    "queryStringParameters": {"initiative_id": "mi-iniciativa"},
    "headers": {"Content-Type": "application/json"}
  }'
```

---

## ✏️ Actualizar Regla

```bash
# Reemplaza RULE_ID con el ID real de la regla
curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "PUT",
    "path": "/api/v1/RULE_ID",
    "body": "{\"category\": \"security\", \"description\": \"Regla actualizada\", \"rule\": \"user.role == admin\"}",
    "headers": {"Content-Type": "application/json"}
  }'
```

---

## 🗑️ Eliminar Regla

```bash
# Reemplaza RULE_ID con el ID real de la regla
curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "DELETE",
    "path": "/api/v1/RULE_ID",
    "headers": {"Content-Type": "application/json"}
  }'
```

---

## 🔧 Script de Prueba Rápida

```bash
#!/bin/bash

# Crear regla y obtener su ID
echo "Creando regla..."
RESPONSE=$(curl -s -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "POST",
    "path": "/api/v1/",
    "body": "{\"initiative_id\": \"test-rapido\", \"category\": \"test\", \"description\": \"Prueba rápida\", \"rule\": \"user.test == true\"}",
    "headers": {"Content-Type": "application/json"}
  }')

RULE_ID=$(echo $RESPONSE | jq -r '.body' | jq -r '.id')
echo "Regla creada con ID: $RULE_ID"

# Obtener la regla
echo "Obteniendo regla..."
curl -s -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d "{
    \"httpMethod\": \"GET\",
    \"path\": \"/api/v1/$RULE_ID\",
    \"headers\": {\"Content-Type\": \"application/json\"}
  }" | jq .

# Eliminar la regla
echo "Eliminando regla..."
curl -s -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d "{
    \"httpMethod\": \"DELETE\",
    \"path\": \"/api/v1/$RULE_ID\",
    \"headers\": {\"Content-Type\": \"application/json\"}
  }" | jq .

echo "¡Prueba completada!"
```

---

## 📊 Verificar Estado

```bash
# Verificar contenedores
make check-containers

# Verificar base de datos
aws dynamodb scan --table-name col-stage-reauth-rulesdb --endpoint-url http://localhost:8000

# Ver logs de la aplicación
docker logs lambda-go-dev --tail 20
```

---

## 🚨 Solución de Problemas

### Error: "Connection refused"
```bash
# Verificar que la aplicación esté corriendo
docker ps | grep lambda-go-dev

# Si no está corriendo, iniciar el entorno
make start
```

### Error: "jq: command not found"
```bash
# Instalar jq
sudo apt-get install jq
# o
brew install jq
```

### Error: "Table not found"
```bash
# Verificar que DynamoDB esté corriendo
docker ps | grep dynamodb-local

# Si no está, iniciar servicios
make compose-up
```

---

*Documentación de ejemplos rápidos - Para uso en desarrollo*
