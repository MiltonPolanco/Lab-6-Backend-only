# Series Tracker API

## Descripción

La **Series Tracker API** es una aplicación RESTful desarrollada en Go para gestionar un backlog de series. Permite realizar operaciones CRUD (crear, leer, actualizar y eliminar series) y actualizaciones parciales (usando PATCH) en campos específicos, tales como:
- **Status**: Actualizar el estado de una serie.
- **Episode**: Incrementar el número del último episodio visto.
- **Upvote/Downvote**: Incrementar o disminuir la puntuación (ranking) de una serie.

Además, el proyecto incluye un frontend (ubicado en la carpeta `static`) que consume esta API para mostrar, filtrar, ordenar y actualizar la información de las series.

---

## Endpoints de la API

La API se encuentra disponible en `http://localhost:8080/api` y cuenta con los siguientes endpoints:

### Endpoints Básicos

- **GET /api/series**  
  **Descripción:** Obtiene la lista de series.  
  **Parámetros de query (opcional):**
  - `search`: Texto para buscar en el título.
  - `status`: Filtrar por estado (por ejemplo, "Watching", "Completed", etc.).
  - `sort`: Ordenar por ranking; valores posibles: `asc` o `desc`.

- **GET /api/series/{id}**  
  **Descripción:** Obtiene los detalles de la serie con el ID especificado.

- **POST /api/series**  
  **Descripción:** Crea una nueva serie.  
  **Body JSON Ejemplo:**
  ```json
  {
    "title": "Mi Serie",
    "status": "Watching",
    "lastEpisodeWatched": 0,
    "totalEpisodes": 10,
    "ranking": 5
  }
