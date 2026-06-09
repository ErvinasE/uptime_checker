# Deploy to Railway (recommended)

Railway can host the full stack (MySQL + Go API + Next.js) with a free trial credit. One public URL for the frontend is enough for your GitHub profile.

## 1. Push code to GitHub

Make sure `main` is up to date on `github.com/ErvinasE/uptime_checker`.

## 2. Create a Railway project

1. Go to [railway.app](https://railway.app) and sign in with GitHub.
2. **New Project** → **Deploy from GitHub repo** → select `uptime_checker`.

## 3. Add MySQL

1. In the project, click **+ New** → **Database** → **MySQL**.
2. Open the MySQL service → **Variables** and note:
   - `MYSQLHOST`
   - `MYSQLPORT`
   - `MYSQLUSER`
   - `MYSQLPASSWORD`
   - `MYSQLDATABASE`

## 4. Deploy the Go API

1. **+ New** → **GitHub Repo** → same repo (or **Empty Service** and link repo).
2. **Settings** → **Root Directory**: leave empty (repo root).
3. **Settings** → **Config file**: `railway.json` (auto-detected).
4. **Variables** → add:

| Variable | Value |
|----------|-------|
| `DB_HOST` | `${{MySQL.MYSQLHOST}}` (use Railway reference to your MySQL service name) |
| `DB_PORT` | `${{MySQL.MYSQLPORT}}` |
| `DB_USER` | `${{MySQL.MYSQLUSER}}` |
| `DB_PASSWORD` | `${{MySQL.MYSQLPASSWORD}}` |
| `DB_NAME` | `${{MySQL.MYSQLDATABASE}}` |
| `CHECK_INTERVAL_MINUTES` | `144` |

Replace `MySQL` with your MySQL service name if different.

5. **Settings** → **Networking** → **Generate Domain** (e.g. `uptime-api-production.up.railway.app`).

## 5. Deploy the frontend

1. **+ New** → **GitHub Repo** → same repo again.
2. **Settings** → **Root Directory**: `frontend`
3. **Variables**:

| Variable | Value |
|----------|-------|
| `API_URL` | `https://YOUR-API-DOMAIN` (from step 4, no trailing slash) |

Example: `https://uptime-api-production.up.railway.app`

4. **Settings** → **Networking** → **Generate Domain** (e.g. `uptime-checker-production.up.railway.app`).

This is the URL you share on your GitHub profile.

## 6. Verify

- API: `https://YOUR-API-DOMAIN/health` → `{"status":"ok"}`
- Site: `https://YOUR-FRONTEND-DOMAIN` → dashboard loads
- Search and website detail pages work (browser uses `/api` proxy to the backend)

## 7. Link on GitHub

### Repository website

Repo → **Settings** → **General** → **Website** → paste your **frontend** Railway URL.

### Profile README

In your profile repo (`ErvinasE/ErvinasE` or similar), add:

```markdown
- [Server Uptime Checker](https://YOUR-FRONTEND-DOMAIN) — live status of 50 popular websites
```

### README badge (optional)

In this repo's `README.md`, update the Live Demo link once deployed.

## Notes

- The Go API does not need to be public for the site to work, but a public domain helps with debugging.
- Railway may sleep free/trial services; the first visit after idle can be slow.
- Migrations run automatically when the API starts.

## Local Docker vs production

| Environment | `API_URL` (frontend) |
|-------------|----------------------|
| Docker Compose | `http://app:8080` |
| Railway | `https://your-api.up.railway.app` |

Browser requests always use `/api/*` on the frontend, which proxies to `API_URL`.
