# HaH Open Cookbook (BACKEND-ONLY, NO AUTH)

Cookbook service for [Haven & Hearth](https://www.havenandhearth.com/) food recipes data. This project doesn't have any public instances, you should self-host it.

## What The Project Does

The application stores and serves cookbook recipes, grouped by cookbook version.

## Stack

- Backend: Go 1.24+/1.25, standard `net/http`, `sqlx`, `modernc.org/sqlite`, `logrus`
- Storage: SQLite
- Runtime: Docker

## Architecture Overview

### Backend

Entry point: `backend/cmd/main.go`

Main responsibilities:

- load env vars, optionally from `--env`
- scan `DATA_DIR` recursively for `.db` files
- create one storage per cookbook database

Database naming rules:

- any `<name>.db` becomes a cookbook version with alias `<name>`
- filenames containing `archived` are exposed as archived cookbook versions and are read-only from the API perspective

## Backend Details

### Storage Model

The backend uses SQLite with three logical schemas:

- `recipes`: main recipe metadata and all FEP columns
- `ingredients`: ingredient name/percentage rows linked to recipes

Important behavior:

- recipes are deduplicated by a computed MD5 hash
- exports can rebuild a portable cookbook SQLite file
- recipe timestamps are used in API responses to show date ranges in the UI

### API Endpoints

| Method | Path                     | Purpose                                       |
| ------ | ------------------------ | --------------------------------------------- |
| `GET`  | `/versions`              | Get all cookbooks versions                    |
| `GET`  | `/api/{version}/recipes` | Filtered, sorted, paginated recipe query      |
| `GET`  | `/api/{version}/export`  | Export cookbook as `default` or `json`        |
| `POST` | `/api/{version}/recipe`  | Insert recipe batch for non-archived versions |

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

### Prerequisites

- Go
- Docker

### Backend development

```bash
cd backend
go run ./cmd/main.go --env=.env --port=8080
```

Notes:

- `DATA_DIR` is scanned recursively for `.db` files
- when `DEBUG=true`, the backend enables permissive CORS for local frontend development

### Image Build

```bash
docker build -t hah-open-cookbook:backend-only .
```

The `Dockerfile` is multi-stage:

1. build the backend binary with Go Alpine
2. copy everything into an Alpine runtime image

### Runtime Behavior

```bash
docker run -v $(pwd)/data:/app/data \
           -e DEBUG="true" \
           -e DATA_DIR="/app/data" \
           -u $(id -u):$(id -g) \
           -p 8080:8080 \
           hah-open-cookbook:backend-only
```

The final container:

- exposes port `8080`
- stores cookbook data under `/app/data`
- expects `DEBUG`, `DATA_DIR`
