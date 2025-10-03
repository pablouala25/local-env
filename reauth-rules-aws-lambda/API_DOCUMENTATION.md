# 📚 Documentación de la API - Rules Management

Esta documentación contiene ejemplos de uso de todos los endpoints de la API de gestión de reglas.

## 🚀 Configuración del Entorno

### Contenedores Requeridos
```bash
# Verificar que todos los contenedores estén corriendo
make check-containers

# O iniciar el entorno completo
make start
```

### Endpoint Base
- **Lambda Directo:** `http://localhost:9000/2015-03-31/functions/function/invocations`
- **SAM Local:** `http://localhost:3000/api/v1/` (requiere autenticación)

---

## 📋 Endpoints Disponibles

### 1. ➕ **CREAR REGLA**
**Crear una nueva regla de reautenticación**

```bash
curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "POST",
    "path": "/api/v1/",
    "body": "{\"initiative_id\": \"auth-team-2024\", \"category\": \"authentication\", \"description\": \"Regla de edad mínima para transacciones\", \"rule\": \"user.age >= 18\"}",
    "headers": {"Content-Type": "application/json"}
  }'
```

**Respuesta exitosa:**
```json
{
  "statusCode": 200,
  "body": "{\"id\":\"01K5PQBH324TY8KF4D89Z80BRQ\",\"status\":\"created\"}"
}
```

**Campos requeridos:**
- `initiative_id`: ID de la iniciativa (string)
- `category`: Categoría de la regla (string)
- `description`: Descripción de la regla (string)
- `rule`: Lógica de la regla (string)

---

### 2. 🔍 **OBTENER REGLA POR ID**
**Obtener una regla específica por su ID**

```bash
curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "GET",
    "path": "/api/v1/01K5PQBH324TY8KF4D89Z80BRQ",
    "headers": {"Content-Type": "application/json"}
  }'
```

**Respuesta exitosa:**
```json
{
  "statusCode": 200,
  "body": "{\"debug\":\"Handler ejecutado sin errores\",\"id\":\"01K5PQBH324TY8KF4D89Z80BRQ\",\"message\":\"Get rule handler funcionando correctamente\",\"status\":\"ok\"}"
}
```

---

### 3. 📋 **LISTAR REGLAS**
**Obtener todas las reglas con paginación**

```bash
# Listar todas las reglas
curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "GET",
    "path": "/api/v1/rules",
    "headers": {"Content-Type": "application/json"}
  }'

# Listar con paginación (limit y next_token)
curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "GET",
    "path": "/api/v1/rules",
    "queryStringParameters": {"limit": "10", "next_token": ""},
    "headers": {"Content-Type": "application/json"}
  }'
```

**Respuesta exitosa:**
```json
{
  "statusCode": 200,
  "body": "{\"items\":[{\"id\":\"01K5PQBH324TY8KF4D89Z80BRQ\",\"initiative_id\":\"auth-team-2024\",\"category\":\"authentication\",\"description\":\"Regla de edad mínima\",\"rule\":\"user.age >= 18\"}],\"next_token\":\"\"}"
}
```

**Parámetros opcionales:**
- `limit`: Número máximo de reglas a retornar (default: 50)
- `next_token`: Token para paginación

---

### 4. 🎯 **OBTENER REGLAS POR INICIATIVA**
**Filtrar reglas por iniciativa específica**

```bash
curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "GET",
    "path": "/api/v1/rules/by-initiative",
    "queryStringParameters": {"initiative_id": "auth-team-2024"},
    "headers": {"Content-Type": "application/json"}
  }'
```

**Respuesta exitosa:**
```json
{
  "statusCode": 200,
  "body": "{\"debug\":\"Handler ejecutado sin errores\",\"initiative_id\":\"auth-team-2024\",\"message\":\"Get rules by initiative handler funcionando correctamente\",\"status\":\"ok\"}"
}
```

**Parámetros requeridos:**
- `initiative_id`: ID de la iniciativa a filtrar

---

### 5. ✏️ **ACTUALIZAR REGLA**
**Actualizar una regla existente**

```bash
curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "PUT",
    "path": "/api/v1/01K5PQBH324TY8KF4D89Z80BRQ",
    "body": "{\"category\": \"security\", \"description\": \"Regla actualizada de seguridad\", \"rule\": \"user.role == admin AND user.age >= 21\"}",
    "headers": {"Content-Type": "application/json"}
  }'
```

**Respuesta exitosa:**
```json
{
  "statusCode": 200,
  "body": "{\"id\":\"01K5PQBH324TY8KF4D89Z80BRQ\",\"status\":\"updated\"}"
}
```

**Nota:** La actualización reemplaza completamente el documento en DynamoDB.

---

### 6. 🗑️ **ELIMINAR REGLA**
**Eliminar una regla por ID**

```bash
curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "DELETE",
    "path": "/api/v1/01K5PQBH324TY8KF4D89Z80BRQ",
    "headers": {"Content-Type": "application/json"}
  }'
```

**Respuesta exitosa:**
```json
{
  "statusCode": 200,
  "body": "{\"id\":\"01K5PQBH324TY8KF4D89Z80BRQ\",\"status\":\"deleted\"}"
}
```

---

## 🔧 Scripts de Prueba

### Script Completo de Pruebas
```bash
#!/bin/bash

echo "🧪 Ejecutando pruebas completas de la API..."

# 1. Crear regla
echo "=== CREAR REGLA ==="
CREATE_RESPONSE=$(curl -s -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "POST",
    "path": "/api/v1/",
    "body": "{\"initiative_id\": \"test-script\", \"category\": \"test\", \"description\": \"Regla de prueba\", \"rule\": \"user.test == true\"}",
    "headers": {"Content-Type": "application/json"}
  }')

echo $CREATE_RESPONSE | jq .
RULE_ID=$(echo $CREATE_RESPONSE | jq -r '.body' | jq -r '.id')

# 2. Obtener regla
echo "=== OBTENER REGLA ==="
curl -s -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d "{
    \"httpMethod\": \"GET\",
    \"path\": \"/api/v1/$RULE_ID\",
    \"headers\": {\"Content-Type\": \"application/json\"}
  }" | jq .

# 3. Listar reglas
echo "=== LISTAR REGLAS ==="
curl -s -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "GET",
    "path": "/api/v1/rules",
    "headers": {"Content-Type": "application/json"}
  }' | jq .

# 4. Actualizar regla
echo "=== ACTUALIZAR REGLA ==="
curl -s -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d "{
    \"httpMethod\": \"PUT\",
    \"path\": \"/api/v1/$RULE_ID\",
    \"body\": \"{\\\"category\\\": \\\"updated\\\", \\\"description\\\": \\\"Regla actualizada\\\", \\\"rule\\\": \\\"user.updated == true\\\"}\",
    \"headers\": {\"Content-Type\": \"application/json\"}
  }" | jq .

# 5. Eliminar regla
echo "=== ELIMINAR REGLA ==="
curl -s -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d "{
    \"httpMethod\": \"DELETE\",
    \"path\": \"/api/v1/$RULE_ID\",
    \"headers\": {\"Content-Type\": \"application/json\"}
  }" | jq .

echo "✅ Pruebas completadas!"
```

---

## 📊 Códigos de Respuesta HTTP

| Código | Descripción |
|--------|-------------|
| `200` | ✅ Operación exitosa |
| `400` | ❌ Error en la petición (JSON inválido, campos faltantes) |
| `404` | ❌ Recurso no encontrado |
| `500` | ❌ Error interno del servidor |

---

## 🗄️ Estructura de Datos

### Regla (Rule)
```json
{
  "id": "01K5PQBH324TY8KF4D89Z80BRQ",
  "initiative_id": "auth-team-2024",
  "category": "authentication",
  "description": "Regla de edad mínima para transacciones",
  "rule": "user.age >= 18"
}
```

### Campos
- `id`: Identificador único generado automáticamente (ULID)
- `initiative_id`: ID de la iniciativa (requerido)
- `category`: Categoría de la regla (requerido)
- `description`: Descripción de la regla (requerido)
- `rule`: Lógica de la regla (requerido)

---

## 🚨 Manejo de Errores

### Error de Validación
```json
{
  "statusCode": 400,
  "body": "invalid JSON body"
}
```

### Error de Recurso No Encontrado
```json
{
  "statusCode": 404,
  "body": "{\"error\":{\"code\":\"not_found\",\"message\":\"rule not found\"}}"
}
```

### Error Interno
```json
{
  "statusCode": 500,
  "body": "{\"error\":{\"code\":\"unknown\",\"message\":\"internal error\"}}"
}
```

---

## 🔍 Verificación del Entorno

### Verificar Contenedores
```bash
make check-containers
```

### Verificar Base de Datos
```bash
aws dynamodb scan --table-name col-stage-reauth-rulesdb --endpoint-url http://localhost:8000
```

### Verificar Logs
```bash
docker logs lambda-go-dev --tail 20
```

---

## 📝 Notas Importantes

1. **IDs Únicos**: Los IDs se generan automáticamente usando ULID
2. **Paginación**: El endpoint de listar soporta `limit` y `next_token`
3. **Actualización**: La actualización reemplaza completamente el documento
4. **Validación**: Todos los campos son requeridos para crear reglas
5. **Base de Datos**: Los datos se persisten en DynamoDB local
6. **Autenticación**: SAM local requiere autenticación, el contenedor lambda directo no

---

*Documentación generada automáticamente - Última actualización: $(date)*
