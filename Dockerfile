# 1. Build backend
FROM golang:1.25-alpine AS backend-builder
WORKDIR /build

RUN apk add --no-cache git

COPY backend/go.mod backend/go.sum ./

RUN go mod download

COPY backend/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# 2. Final
FROM alpine:latest

RUN apk add --no-cache sqlite

RUN adduser -D -u 1000 -g 'app' app

COPY --from=backend-builder /build/main /app/main

RUN mkdir -p /app/data && \
    chown -R app:app /app && \
    chmod +x /app/main

EXPOSE 8080

CMD ["/app/main", " --port=8080"]
