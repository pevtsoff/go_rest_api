## REST API (Go + Gin + GORM + Postgres)

### Prerequisites
- Go 1.25+
- Docker and Docker Compose (for containerized Postgres)

### Environment
The app reads DB settings from `.env` at startup:

```env
PORT=3000 - port to raise the go rest api on
TAG=go-rest-api - tag suffix for all the docker images
DB_CONNECTION_STRING=host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=UTC
```

Auto-migration runs on startup (users, posts).

## Local setup

### Option A: Without Docker Compose (Postgres in Docker)
1) Start Postgres:
```bash
docker compose up -d postgres
```

2) Ensure `.env` is present (see above), then run the app:
```bash
go run ./main.go
```

The server listens on `http://localhost:3000`.

### Option B: With Docker Compose (external/local Postgres)
2) Set `.env` with your connection string, for example:
```env
docker compose up -d
```

## API
Swagger/OpenAPI is generated at runtime:
- OpenAPI JSON: `GET /openapi.json`
- Swagger UI: `GET /swagger/index.html` (via `GET /swagger/*any`)

### Endpoints
- Users
  - `POST /users/`
  - `GET /users/:id`
- Posts
  - `POST /posts/`
  - `GET /posts/`
  - `GET /posts/:id`
  - `PATCH /posts/:id`
  - `DELETE /posts/:id`

### cURL examples

Create user:
```bash
curl -sS -X POST http://localhost:3000/users/ \
  -H 'Content-Type: application/json' \
  -d '{"name":"John Doe"}'
```

Get user:
```bash
curl -sS http://localhost:3000/users/1
```

Create post:
```bash
curl -sS -X POST http://localhost:3000/posts/ \
  -H 'Content-Type: application/json' \
  -d '{"title":"Hello","body":"World"}'
```

List posts:
```bash
curl -sS http://localhost:3000/posts/
```

Show post:
```bash
curl -sS http://localhost:3000/posts/1
```

Update post:
```bash
curl -sS -X PATCH http://localhost:3000/posts/1 \
  -H 'Content-Type: application/json' \
  -d '{"title":"Updated","body":"Updated body"}'
```

Delete post:
```bash
curl -sS -X DELETE http://localhost:3000/posts/1
```

There are also sample HTTP files in `http/` you can use with REST clients.

## Tests

### What the tests do
- Use real Postgres (not sqlite).
- Each test runs inside a DB transaction which is rolled back at the end for isolation and speed.
- Builders create test data on-the-fly:
  - `tests/testutils/builders.go` exposes `NewUserBuilder` and `NewPostBuilder`.

### Running tests locally
1) Start Postgres and create a `test` database (once):
```bash
docker compose up -d postgres
docker compose exec -T postgres psql -U postgres -c "CREATE DATABASE test;" || true
```

2) Run tests (point to the test DB):
```bash
DB_CONNECTION_STRING="host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable TimeZone=UTC" \
go test ./... -v
```

3) With coverage:
```bash
DB_CONNECTION_STRING="host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable TimeZone=UTC" \
go test ./tests -v -covermode=atomic -coverpkg="$(go list ./... | grep -v '/tests' | tr '\n' ',' | sed 's/,$//')" -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Notes
- If you see GORM "record not found" logs during not-found test cases, that is expected. To silence in tests, you can set the GORM logger to Silent in the test DB session.

## CI (GitHub Actions)
- Workflow at `.github/workflows/tests.yml`:
  - Starts Postgres service.
  - Loads `.env` for `POSTGRES_*` defaults (falls back to `postgres` / `test`).
  - Runs tests against `127.0.0.1:5432`.
  - Generates coverage (`coverage.out`, `coverage.html`).
  - Uploads coverage artifacts and writes a coverage summary to the job summary.
