# city-counter

Counts cities by starting letter using data retrieved from api.geonames.org.
The current backend counts matches among up to 50 places returned by GeoNames,
not all cities worldwide.

Source code lives in `frontend/` (React and Nginx) and `backend/` (Go backend).

## Run with Docker Compose

Install Docker Engine with Docker Compose, then run:

```sh
docker compose up --build
```

Open http://localhost:8080.

- `frontend` builds the React app and serves it with Nginx.
- `backend` builds and runs the Go API.
- Nginx forwards `/cities/` requests to `backend:${PORT}` (8080 by default) over the Compose network.
- Only Nginx is published to the host, on `127.0.0.1:${UI_PORT}` (8080 by default).

Check the API through Nginx:

```sh
curl 'http://localhost:8080/cities/count?first_letter=L'
```

Stop the containers:

```sh
docker compose down
```

## Configuration

Copy `.env.example` to `.env` and edit it before starting Compose:

```sh
cp .env.example .env
docker compose up --build
```

| Variable | Default | Meaning |
| --- | --- | --- |
| `GEONAMES_URL` | URL in `.env.example` | Full upstream HTTP/HTTPS URL, including query parameters, returns the cities |
| `CACHE_ENABLED` | `true` | Set to `false` to fetch from GeoNames on every request, true - fetch GeoNames from memory cache |
| `CACHE_TTL` | `1m` | Positive Go duration, such as `30s` or `5m` |
| `BACKEND_PORT` | `8080` | Backend listening port |
| `FRONTEND_PORT` | `8080` | Local port exposed by the frontend container |

## Develop without Docker

Run the backend from its directory:

```sh
cd backend
go run ./cmd/server
```

In another terminal, with Node.js 22.12+ and npm installed:

```sh
cd frontend
npm ci
npm run dev
```

Open http://localhost:5173. Vite forwards `/cities` requests to Go on port 8080.

## Test/Check

```sh
go -C backend test -race ./...
npm --prefix frontend run build
```
