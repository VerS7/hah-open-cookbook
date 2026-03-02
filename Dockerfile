# 1. Build frontend
FROM node:22-alpine AS frontend-builder
WORKDIR /build

COPY frontend/package*.json ./

RUN npm ci

COPY frontend/ ./

RUN npm run build

# 2. Build backend
FROM golang:1.25-alpine AS backend-builder
WORKDIR /build

RUN apk add --no-cache git

COPY backend/go.mod backend/go.sum ./

RUN go mod download

COPY backend/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# 3. Final
FROM nginx:alpine

RUN apk add --no-cache supervisor sqlite

RUN adduser -D -u 1000 -g 'app' app

COPY --from=frontend-builder /build/dist /usr/share/nginx/html
COPY --from=backend-builder /build/main /app/main
COPY nginx.conf /etc/nginx/nginx.conf
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

RUN mkdir -p /app/data && \
    chown -R app:app /app && \
    chmod +x /app/main

EXPOSE 80

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
