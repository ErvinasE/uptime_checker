# Server Uptime Checker

Monitor uptime of 50 popular websites. Built with Go, MySQL, Next.js, and Docker.

**Live demo:** _Add your Railway URL after deploying — see [DEPLOY.md](DEPLOY.md)_

## Features

- Checks ~10 times per day per site (every ~144 minutes with jitter)
- Home page showing currently up and down websites
- Search with autocomplete suggestions
- Detail page with 24-hour check history (date, time, response time in ms)
- Downtime incident tracking

## Quick start (local)

```bash
cp .env.example .env
docker compose up --build
```

- Frontend: http://localhost:3000
- API: http://localhost:8080

## Deploy online

See **[DEPLOY.md](DEPLOY.md)** for step-by-step Railway hosting (free trial). After deploy, add the URL to your GitHub repo **Website** field and profile README.

## API endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/websites` | All websites with latest check |
| GET | `/websites/:id` | Single website |
| GET | `/websites/:id/history` | Check history |

## Configuration

See `.env.example` for all options:

- `CHECK_INTERVAL_MINUTES` — minutes between check cycles (default: 144)
- `CHECK_RETENTION_DAYS` — how long to keep check records (default: 7)

## Local development (without Docker)

### Backend

```bash
export DB_HOST=localhost DB_PORT=3306 DB_USER=uptime DB_PASSWORD=uptimepassword DB_NAME=uptime
go run ./cmd/server
```

### Frontend

```bash
cd frontend
npm install
API_URL=http://localhost:8080 npm run dev
```

## Tests

```bash
go test ./...
```

## How uptime is determined

A site is marked **up** when the server responds with any HTTP status below 500 (including 403). Many popular sites block automated clients with bot protection even when they are online in a browser. A **down** result means a timeout, connection error, or HTTP 5xx.

## Troubleshooting

**"Could not reach the API" on the home page**

1. Rebuild after config changes: `docker compose up --build`
2. Confirm the API is up: `curl http://localhost:8080/health`
3. Check app logs: `docker compose logs app`
4. In Docker, set frontend `API_URL=http://app:8080`

**Production (Railway)**

1. Frontend `API_URL` must be the public Go API URL (https).
2. Redeploy frontend after changing `API_URL`.
