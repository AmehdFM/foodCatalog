# Food Catalog API 🍎

API RESTful para gestionar un catálogo de alimentos con información nutricional y sistema de etiquetado flexible para aplicaciones de gyms y entrenadores personales.

## 🚀 Características

- **Catálogo de Alimentos**: Información nutricional completa (calorías, proteínas, carbohidratos, grasas, fibra)
- **Sistema de Tags**: Etiquetado flexible para vitaminas, minerales, restricciones dietéticas, etc.
- **Filtros Avanzados**: Búsqueda por nombre, tags múltiples, fechas, estado
- **Paginación**: Control de resultados con metadata completa
- **Autenticación**: Sistema de tokens para proteger los endpoints
- **Batch Queries**: Consulta múltiples recursos en una sola petición
- **Seed System**: Carga masiva de datos desde archivo JSON

## 📋 Requisitos

- Go 1.21+
- PostgreSQL 12+

## 🛠️ Instalación

1. **Clonar el repositorio**
```bash
git clone <tu-repo>
cd food-catalog-api
```

2. **Instalar dependencias**
```bash
go mod download
```

3. **Configurar variables de entorno**
```bash
cp .env.example .env
# Editar .env con tus credenciales
```

Configuración necesaria en `.env`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=tupassword
DB_NAME=food_catalog
DB_SSLMODE=disable

SERVER_PORT=8080

# Tokens separados por comas
AUTH_TOKENS=token123,token456,token789
```

4. **Crear la base de datos**
```bash
createdb food_catalog
```

## 🎯 Uso

### Iniciar el servidor
```bash
go run cmd/server/main.go
```

El servidor estará disponible en `http://localhost:8080`

### Cargar datos desde JSON
```bash
go run cmd/seed/main.go -file=data.json
```

El comando seed:
- Crea nuevos registros si no existen
- Actualiza registros existentes si la fecha `updated_at` del JSON es más reciente
- Omite registros que no han cambiado

## 📡 API Endpoints

### Health Check
```http
GET /health
```

No requiere autenticación. Verifica el estado del servidor y la base de datos.

### Autenticación

Todos los endpoints bajo `/api/v1` requieren un header de autorización:
```http
Authorization: Bearer tu-token-aqui
```

O simplemente:
```http
Authorization: tu-token-aqui
```

### Alimentos

#### Listar alimentos
```http
GET /api/v1/alimentos
```

**Parámetros de query:**
- `page` (int): Número de página (default: 1)
- `limit` (int): Resultados por página (default: 20, max: 100)
- `search` (string): Búsqueda por nombre
- `tags` (string): Filtrar por tags (separados por comas). Ej: `vegano,proteina-alta`
- `activo` (bool): Filtrar por estado activo
- `date_from` (date): Fecha desde (formato: YYYY-MM-DD)
- `date_to` (date): Fecha hasta (formato: YYYY-MM-DD)
- `sort_by` (string): Campo para ordenar (ej: `nombre`, `calorias`, `created_at`)
- `sort_order` (string): Orden ascendente o descendente (`asc` o `desc`)

**Ejemplos:**
```bash
# Listar todos los alimentos
curl -H "Authorization: token123" http://localhost:8080/api/v1/alimentos

# Buscar alimentos veganos con proteína alta
curl -H "Authorization: token123" \
  "http://localhost:8080/api/v1/alimentos?tags=vegano,proteina-alta"

# Buscar "pollo" ordenado por proteínas
curl -H "Authorization: token123" \
  "http://localhost:8080/api/v1/alimentos?search=pollo&sort_by=proteinas&sort_order=desc"

# Paginación
curl -H "Authorization: token123" \
  "http://localhost:8080/api/v1/alimentos?page=2&limit=10"
```

**Respuesta:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "nombre": "Pechuga de pollo",
      "calorias": 165,
      "proteinas": 31,
      "carbohidratos": 0,
      "grasas": 3.6,
      "fibra": 0,
      "imagen_url": "https://example.com/images/pollo.jpg",
      "icon_url": "https://example.com/icons/pollo.svg",
      "activo": true,
      "tags": [
        {
          "id": 2,
          "nombre": "proteina-alta",
          "imagen_url": "https://example.com/images/proteina.jpg",
          "icon_url": "https://example.com/icons/proteina.svg",
          "color": "#FF5722",
          "created_at": "2024-01-15T10:00:00Z",
          "updated_at": "2024-01-15T10:00:00Z"
        }
      ],
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total_pages": 1,
    "total_items": 10
  }
}
```

#### Obtener alimento por ID
```http
GET /api/v1/alimentos/{id}
```

**Ejemplo:**
```bash
curl -H "Authorization: token123" http://localhost:8080/api/v1/alimentos/1
```

#### Obtener múltiples alimentos por IDs
```http
GET /api/v1/alimentos/batch?ids=1,2,3
```

**Ejemplo:**
```bash
curl -H "Authorization: token123" \
  "http://localhost:8080/api/v1/alimentos/batch?ids=1,2,5"
```

### Tags

#### Listar tags
```http
GET /api/v1/tags
```

**Parámetros de query:**
- `page` (int): Número de página (default: 1)
- `limit` (int): Resultados por página (default: 20, max: 100)
- `search` (string): Búsqueda por nombre
- `date_from` (date): Fecha desde (formato: YYYY-MM-DD)
- `date_to` (date): Fecha hasta (formato: YYYY-MM-DD)
- `sort_by` (string): Campo para ordenar
- `sort_order` (string): Orden (`asc` o `desc`)

**Ejemplo:**
```bash
curl -H "Authorization: token123" http://localhost:8080/api/v1/tags
```

#### Obtener tag por ID
```http
GET /api/v1/tags/{id}
```

**Ejemplo:**
```bash
curl -H "Authorization: token123" http://localhost:8080/api/v1/tags/1
```

#### Obtener tag por nombre
```http
GET /api/v1/tags/nombre/{nombre}
```

**Ejemplo:**
```bash
curl -H "Authorization: token123" http://localhost:8080/api/v1/tags/nombre/vegano
```

#### Obtener múltiples tags por IDs
```http
GET /api/v1/tags/batch?ids=1,2,3
```

**Ejemplo:**
```bash
curl -H "Authorization: token123" \
  "http://localhost:8080/api/v1/tags/batch?ids=1,3,5"
```

## 📊 Estructura del Proyecto
```
food-catalog-api/
├── cmd/
│   ├── server/         # Servidor principal
│   │   └── main.go
│   └── seed/           # Comando para cargar datos
│       └── main.go
├── config/             # Configuración
│   └── config.go
├── database/           # Conexión y migraciones
│   └── database.go
├── handlers/           # Handlers HTTP
│   ├── alimento_handlers.go
│   ├── tag_handlers.go
│   └── health_handlers.go
├── middleware/         # Middlewares
│   └── auth.go
├── models/             # Modelos de datos
│   ├── alimento.go
│   └── tag.go
├── utils/              # Utilidades
│   ├── response.go
│   └── query.go
├── data.json           # Datos de ejemplo
├── .env.example        # Ejemplo de variables de entorno
├── .gitignore
├── go.mod
└── README.md
```

## 💾 Formato del archivo data.json
```json
{
  "tags": [
    {
      "id": 1,
      "nombre": "vegano",
      "imagen_url": "https://example.com/images/vegano.jpg",
      "icon_url": "https://example.com/icons/vegano.svg",
      "color": "#4CAF50",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ],
  "alimentos": [
    {
      "id": 1,
      "nombre": "Pechuga de pollo",
      "calorias": 165,
      "proteinas": 31,
      "carbohidratos": 0,
      "grasas": 3.6,
      "fibra": 0,
      "imagen_url": "https://example.com/images/pollo.jpg",
      "icon_url": "https://example.com/icons/pollo.svg",
      "activo": true,
      "tags": ["proteina-alta", "sin-gluten"],
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ]
}
```

## 🔐 Seguridad

- Los tokens se almacenan en variables de entorno
- Todos los endpoints (excepto `/health`) requieren autenticación
- CORS habilitado para permitir peticiones desde cualquier origen

## 🚀 Despliegue

### Compilar binarios
```bash
# Servidor
go build -o bin/server cmd/server/main.go

# Seed
go build -o bin/seed cmd/seed/main.go
```

### Ejecutar en producción
```bash
./bin/server
```

## 📝 Ejemplos de Casos de Uso

### 1. App de gym busca alimentos veganos altos en proteína
```bash
curl -H "Authorization: token123" \
  "http://localhost:8080/api/v1/alimentos?tags=vegano,proteina-alta&limit=50"
```

### 2. Nutricionista busca alimentos con vitamina B12
```bash
curl -H "Authorization: token123" \
  "http://localhost:8080/api/v1/alimentos?tags=vitamina-b12"
```

### 3. App necesita información de alimentos específicos
```bash
curl -H "Authorization: token123" \
  "http://localhost:8080/api/v1/alimentos/batch?ids=1,5,8,12,15"
```

### 4. Buscar alimentos recién agregados
```bash
curl -H "Authorization: token123" \
  "http://localhost:8080/api/v1/alimentos?date_from=2024-01-01&sort_by=created_at&sort_order=desc"
```