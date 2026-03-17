# HaH Open Cookbook

Cookbook service for [Haven & Hearth](https://www.havenandhearth.com/) food recipes data. This project doesn't have any public instances, you should self-host it.

## What The Project Does

The application stores and serves cookbook recipes, grouped by cookbook version. Users can:

- log in and receive a bearer token
- browse available cookbook versions
- search recipes with a small query DSL
- sort and paginate results
- inspect ingredients and FEP composition
- export cookbook data as standard JSON or a SQLite custom database for sharing
- capture a screenshot of the current table view for sharing

Admin-related behavior is only partially implemented at the moment.

## Stack

- Backend: Go 1.24+/1.25, standard `net/http`, `sqlx`, `modernc.org/sqlite`, `logrus`
- Frontend: Vue 3, TypeScript, Vue Router, Vuetify, Vite, `html2canvas`
- Storage: SQLite
- Runtime: Nginx, Docker

## Architecture Overview

### Backend

Entry point: `backend/cmd/main.go`

Main responsibilities:

- load env vars, optionally from `--env`
- scan `DATA_DIR` recursively for `.db` files
- create one storage per cookbook database
- create a dedicated users storage from `users.db`
- auto-create the default admin user from env vars
- expose public, authenticated, and admin API routes

Database naming rules:

- `users.db` is the users/auth database
- any other `<name>.db` becomes a cookbook version with alias `<name>`
- filenames containing `archived` are exposed as archived cookbook versions and are read-only from the API perspective

### Frontend

Entry point: `frontend/src/main.ts`

The SPA has three routes:

- `/login` for authentication
- `/` for the main cookbook interface
- `/admin` reserved for future admin UI

The main screen includes:

- cookbook version switcher
- recipe table
- search/filter input
- sortable metrics
- pagination

### Runtime

The production image serves the frontend through Nginx and proxies `/api/*` to the Go backend on `127.0.0.1:8080`.

Inside the container:

- `nginx` serves the SPA
- `/app/main` runs the backend
- `crond` executes nightly backups
- `supervisord` keeps all services alive

## Repository Layout

```text
.
├── backend/
│   ├── cmd/main.go
│   ├── go.mod
│   └── internal/
│       ├── api/
│       ├── auth/
│       ├── filter/
│       ├── logger/
│       └── storage/
├── frontend/
│   ├── package.json
│   ├── vite.config.ts
│   ├── public/
│   └── src/
│       ├── components/
│       ├── composables/
│       ├── styles/
│       └── views/
├── .github/workflows/
├── Dockerfile
├── nginx.conf
├── supervisord.conf
├── backup.sh
└── crontab.root
```

## Backend Details

### Storage Model

The backend uses SQLite with three logical schemas:

- `recipes`: main recipe metadata and all FEP columns
- `ingredients`: ingredient name/percentage rows linked to recipes
- `users`: application users, hashed passwords, access tokens, admin flag

Important behavior:

- recipes are deduplicated by a computed MD5 hash
- exports can rebuild a portable cookbook SQLite file
- recipe timestamps are used in API responses to show date ranges in the UI

### Authentication

Authentication is simple and token-based:

- `POST /api/login` checks `username` and `password`
- the response returns `username`, `token`, and `isAdmin`
- authenticated endpoints expect `Authorization: Bearer <token>`

The default admin user is inserted on startup from:

- `ADMIN_USERNAME`
- `ADMIN_PASSWORD`

### API Endpoints

| Method | Path                     | Auth | Purpose                                          |
| ------ | ------------------------ | ---- | ------------------------------------------------ |
| `POST` | `/api/login`             | No   | Login with form fields `username` and `password` |
| `GET`  | `/api/versions`          | No   | List cookbook version aliases                    |
| `GET`  | `/api/{version}/recipes` | User | Filtered, sorted, paginated recipe query         |
| `GET`  | `/api/{version}/export`  | User | Export cookbook as `default` or `json`           |
| `POST` | `/api/{version}/recipe`  | User | Insert recipe batch for non-archived versions    |

### Recipe Query Parameters

`GET /api/{version}/recipes` requires:

- `filter`: query DSL string, must not be empty
- `sort`: `asc` or `desc`
- `by`: sort key
- `p`: page number, minimum `1`
- `l`: page size, `1..500`

Supported sort keys in the backend:

- `name`
- `energy`
- `hunger`
- `total`
- `feph`
- individual FEP keys such as `str1`, `str2`, `agi1`, `agi2`, ...
- grouped FEP keys such as `str`, `agi`, `int`, `con`, `prc`, `csm`, `dex`, `wil`, `psy`

### Filter DSL

The filter language is parsed in `backend/internal/filter/parser.go` and validated in the storage layer.

Supported operators:

- `=`
- `!=`
- `>`
- `<`
- `>=`
- `<=`

Supported fields visible in the code:

- `name`
- `from`
- `hunger`
- `energy`
- `total`
- individual FEP fields such as `str1`, `str2`, `agi1`, ...

Value types:

- integers: `10`
- floats: `3.5`
- percentages: `30%`
- strings: `"steak"`
- string lists: `"dill","pepper","sage"`

Examples:

```text
name=steak;str2>30%;energy>=300;from!="dill","pepper","sage"
```

### Export Formats

The backend currently supports:

- `type=default`: SQLite custom-schema export with `recipes` and `ingredients` tables
- `type=json`: Standard JSON (used in [Nurgling2](https://github.com/aleksandrsvoboda/nurgling2) or [Hurricane](https://github.com/Nightdawg/Hurricane))

## Frontend Details

### Main UX

The main page is centered around `RecipesTable.vue`. It provides:

- cookbook version selection
- URL-synced state for selected version, query, page size, and page
- table sorting by totals, hunger, energy, FEP/H, and FEP groups
- ingredient chip collapsing with popup overflow
- hover expansion for truncated recipe names
- screenshot mode that widens the table for clean captures

Query state is stored in the URL with:

- `v`: cookbook version
- `q`: filter string
- `l`: page size
- `p`: page number

### Frontend Environment

Frontend API base URL is taken from `VITE_API_URL`.

Project defaults:

- `frontend/.env.development` points to `http://localhost:8080`
- `frontend/.env.production` is empty, which makes the frontend use same-origin `/api` behind Nginx

### Export And Token Utilities

The top bar lets a logged-in user:

- copy the current bearer token to the clipboard
- export the currently selected cookbook version
- log out

## Local Development

### Prerequisites

- Go
- Node.js matching `frontend/package.json` engines
- npm

### Backend development

Create an ignored env file such as `backend/.env`:

```env
DEBUG=true
DATA_DIR=./data/
ADMIN_USERNAME=admin
ADMIN_PASSWORD=change-me
```

Then run:

```bash
cd backend
go run ./cmd/main.go --env=.env --port=8080
```

Notes:

- `DATA_DIR` is scanned recursively for `.db` files
- using a trailing slash in `DATA_DIR` is the safest choice because the fallback `users.db` path is built by string concatenation in code
- when `DEBUG=true`, the backend enables permissive CORS for local frontend development

### Frontend development

Development env points to the local backend by default.

```bash
cd frontend
npm install
npm run dev
```

## Docker And Production Runtime

### Image Build

The `Dockerfile` is multi-stage:

1. build the frontend with Node 22 Alpine
2. build the backend binary with Go Alpine
3. copy everything into an Nginx Alpine runtime image

### Runtime Behavior

The final container:

- exposes port `80`
- stores cookbook data under `/app/data`
- expects `DEBUG`, `DATA_DIR`, `ADMIN_USERNAME`, `ADMIN_PASSWORD`
- runs daily backups via cron

### Backups

`backup.sh`:

- archives `BACKUP_SOURCE_DIR` or `/app/data`
- writes backups to `${DEPLOY_DATA_DIR}/backups` or `/app/backups`
- deletes archives older than 7 days

`crontab.root` schedules this at midnight every day.

## CI/CD

### CI Workflow

`.github/workflows/ci.yml` runs:

- frontend dependency install and build
- backend `go test ./...`
- backend binary build
- Docker image build validation

### Deploy Workflow

`.github/workflows/deploy.yml` performs manual VPS deployment.

Required GitHub secrets:

- `VPS_HOST`
- `VPS_USER`
- `VPS_SSH_KEY`
- `APP_ADMIN_USERNAME`
- `APP_ADMIN_PASSWORD`

Optional GitHub variables:

- `VPS_PORT`
- `APP_DEBUG`
- `DEPLOY_DATA_DIR`
- `APP_DATA_DIR`

Deployment flow:

- build Docker image in GitHub Actions
- save and gzip the image
- copy the archive to the VPS over SSH
- load the image remotely
- stop and remove the previous container
- run the new container with a mounted data directory

## Current Limitations to be fixed

- `frontend/src/views/AdminView.vue` is still a stub
- backend admin user-management handlers are placeholders
- the project has CI build checks, but there are currently no meaningful automated test suites in the source tree
- auth is intentionally lightweight and based on stored access tokens rather than a full session or identity system
