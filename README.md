# Food Catalog API | Catálogo nutricional para aplicaciones

API REST en Go para consultar alimentos, información nutricional y etiquetas desde aplicaciones de entrenamiento, nutrición o gestión de gimnasios.

## El problema que resuelve

Separa el catálogo de alimentos del resto de una aplicación. Un cliente puede buscar alimentos, recuperar sus nutrientes y reutilizar etiquetas sin mantener esa lógica dentro de cada interfaz.

El repositorio implementa consultas y carga de datos por línea de comandos; no incluye una interfaz de administración ni un servicio de recomendaciones nutricionales.

## Funcionalidades

- Consulta de calorías, proteínas, carbohidratos, grasas y fibra.
- Relación de alimentos con múltiples etiquetas.
- Búsqueda por nombre, etiquetas, fechas y estado.
- Paginación y metadatos de resultados, con un máximo de 100 elementos por página.
- Consultas por ID y por lotes.
- Protección de rutas mediante tokens configurados.
- Comando para cargar o actualizar alimentos y etiquetas desde JSON.
- Estado de servicio y base de datos mediante `/health`.

## Tecnologías y arquitectura

**Go · Gorilla Mux · PostgreSQL · GORM · godotenv · Docker**

| Ruta | Responsabilidad |
|---|---|
| [cmd/server/](cmd/server/) | Punto de entrada de la API |
| [cmd/seed/](cmd/seed/) | Carga de datos |
| [config/](config/) | Variables de entorno y tokens |
| [database/](database/) | Conexión y migraciones automáticas |
| [handlers/](handlers/) | Endpoints de alimentos, etiquetas y estado |
| [middleware/](middleware/) | Autorización por token y CORS |
| [models/](models/) | Modelos y relaciones |
| [utils/](utils/) | Parámetros y respuestas |
| [data.json](data.json) | Datos de ejemplo |

**El servidor del catálogo está en `cmd/server/main.go`.** El `main.go` de la raíz es un servidor de ejemplo separado y no expone esta API.

## Ejecutar localmente

Requisitos: Go compatible con [go.mod](go.mod), que declara Go 1.21, y una base PostgreSQL disponible. Crear una base de desarrollo `food_catalog` y su usuario antes de arrancar.

```bash
git clone https://github.com/AmehdFM/foodCatalog.git
cd foodCatalog
go mod download
```

Crear manualmente un archivo `.env` en la raíz:

```dotenv
DATABASE_URL=postgres://usuario:password@localhost:5432/food_catalog?sslmode=disable
SERVER_PORT=8080
AUTH_TOKENS=reemplazar-por-un-token-propio
```

No se necesita copiar un `.env.example`: ese archivo no está incluido. La configuración lee `DATABASE_URL`; no utiliza variables `DB_HOST`, `DB_USER` o `DB_PASSWORD` por separado. El archivo `.env` es opcional si estas variables ya están definidas en el proceso.

```bash
go run ./cmd/seed -file=data.json
go run ./cmd/server
```

El servidor y el comando seed conectan a la base y ejecutan la migración de modelos al iniciar. El seed identifica los registros por nombre y compara `updated_at` para decidir las actualizaciones; usar datos de desarrollo.

## Configuración

| Variable | Uso |
|---|---|
| `DATABASE_URL` | Conexión PostgreSQL, obligatoria |
| `AUTH_TOKENS` | Uno o varios tokens válidos separados por comas, obligatorio |
| `SERVER_PORT` | Puerto HTTP; predeterminado `8080` |

Los tokens identifican acceso al servicio; no implementan cuentas de usuario, roles ni sesiones individuales.

## API

Todas las rutas de `/api/v1` requieren autorización. `/health` no requiere token.

| Método | Ruta | Función |
|---|---|---|
| GET | `/health` | Estado |
| GET | `/api/v1/alimentos` | Alimentos con filtros |
| GET | `/api/v1/alimentos/{id}` | Detalle |
| GET | `/api/v1/alimentos/batch?ids=1,2,3` | Consulta múltiple |
| GET | `/api/v1/tags` | Etiquetas |
| GET | `/api/v1/tags/{id}` | Etiqueta por ID |
| GET | `/api/v1/tags/nombre/{nombre}` | Etiqueta por nombre |
| GET | `/api/v1/tags/batch?ids=1,2,3` | Etiquetas por lote |

Filtros de alimentos: `page`, `limit`, `search`, `tags`, `activo`, `date_from`, `date_to`, `sort_by` y `sort_order`. Las fechas usan `YYYY-MM-DD`.

```bash
curl http://localhost:8080/health
curl -H "Authorization: Bearer reemplazar-por-un-token-propio" \
  "http://localhost:8080/api/v1/alimentos?search=pollo&page=1&limit=10"
curl -H "Authorization: Bearer reemplazar-por-un-token-propio" \
  "http://localhost:8080/api/v1/alimentos?tags=vegano"
```

Los valores de búsqueda son ilustrativos; los resultados dependen del conjunto cargado.

### Formato de carga

[data.json](data.json) contiene las colecciones `tags` y `alimentos`. Los alimentos incluyen nutrientes, estado, etiquetas por nombre y fecha `updated_at` en RFC 3339. Las etiquetas se procesan antes que los alimentos.

## Compilación y Docker

```bash
go build -o server ./cmd/server
go build -o seed ./cmd/seed
docker build -t food-catalog-api .
docker run --rm --env-file .env -p 8080:8080 food-catalog-api
```

Para Docker, modificar la conexión para que PostgreSQL sea accesible desde el contenedor. Una dirección `localhost` dentro del contenedor no apunta al equipo anfitrión.

El Dockerfile incluye los binarios de servidor y seed. La versión de Go de la etapa de compilación debe ser compatible con las dependencias resueltas.

## Estado y mejoras pendientes

API de consulta con separación por responsabilidades y documentación de integración. No hay una suite de pruebas automatizadas en el repositorio.

Antes de exponerla públicamente, restringir `sort_by` a columnas permitidas: actualmente se incorpora a una expresión de ordenamiento SQL. También revisar CORS, gestión de tokens, manejo de errores y tiempos de espera HTTP. La información nutricional depende del origen y revisión de los datos cargados.

## Autor

[Amehd Mendez](https://github.com/AmehdFM) · Backend e integración de sistemas.
