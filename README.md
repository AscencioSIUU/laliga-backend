# La Liga Tracker Backend & Frontend

Este proyecto es el backend y frontend de **La Liga Tracker**, una aplicación para gestionar partidos de La Liga. El backend está desarrollado en Go (usando Gin) y se conecta a una base de datos MySQL. El frontend (archivo `LaLigaTracker.html`) se sirve de forma estática desde el mismo backend.

---

## Requisitos

- Docker y Docker Compose
- Go 1.21 (o superior)
- Git

---

## Estructura del Proyecto

```
lab1and2/
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── init.sql
├── main.go
└── public/
    └── LaLigaTracker.html
└── imgs/
    └── imgs/e594b348-c193-4b46-b2bb-ac4a58d751a0.jpeg
```

---

## Descripción

### Archivos Principales

- **Dockerfile:** Define el proceso de compilación del backend y la inclusión del frontend.
- **docker-compose.yml:** Orquesta los servicios de MySQL (base de datos) y el backend.
- **init.sql:** Script de inicialización que configura el usuario, plugin de autenticación y crea la tabla `matches`.
- **main.go:** Código del backend (API REST) desarrollado con Gin.
- **public/LaLigaTracker.html:** Archivo estático que representa el frontend de la aplicación.

### Rutas creadas

## Endpoints del Backend

| Endpoint           | Método     | Descripción                               | Ejemplo de Petición / Cuerpo                                                            |
| ------------------ | ---------- | ----------------------------------------- | --------------------------------------------------------------------------------------- |
| `/api/matches`     | **GET**    | Obtiene todos los partidos.               | N/A                                                                                     |
| `/api/matches/:id` | **GET**    | Obtiene un partido específico por su ID.  | `GET /api/matches/1`                                                                    |
| `/api/matches`     | **POST**   | Crea un nuevo partido.                    | <pre>{"homeTeam": "Barcelona","awayTeam": "Real Madrid","matchDate":"2025-04-01"}</pre> |
| `/api/matches/:id` | **PUT**    | Actualiza un partido existente por su ID. | <pre>{"homeTeam": "Atletico","awayTeam": "Sevilla","matchDate": "2025-05-10"}</pre>     |
| `/api/matches/:id` | **DELETE** | Elimina un partido específico por su ID.  | `DELETE /api/matches/1`                                                                 |

---

![Screenshot de La Liga Tracker](imgs/e594b348-c193-4b46-b2bb-ac4a58d751a0.jpeg)

---

## Pasos para Levantar el Proyecto

### 1. Clonar el Repositorio

```bash
git clone https://github.com/AscencioSIUU/laliga-backend.git
cd laliga-backend
```

## configurar las variables del .env

```bash
DB_HOST=localhost
DB_PORT=3306
DB_USER=laliga
DB_PASSWORD=laliga
DB_NAME=laliga
```

## Levantar en docker compose

```bash
docker-compose down -v
docker-compose up --build -d
```

## Documentación de la API (Swagger)

El archivo [`swagger.yaml`](./swagger.yaml) describe los endpoints y parámetros de la API siguiendo el estándar OpenAPI 3.0.0.

### Visualizar con Swagger UI

1. Ejecutar el contenedor oficial de Swagger UI:
   ```bash
   docker run --rm -p 8081:8080 \
     -v $(pwd)/swagger.yaml:/usr/share/nginx/html/swagger.yaml \
     swaggerapi/swagger-ui
   ```

### Abrir swagger en el navegador

```
http://localhost:8081
```
