# ---- Build stage ----
FROM golang:1.18-bullseye AS builder

RUN apt-get update && apt-get install -y gcc musl-dev && rm -rf /var/lib/apt/lists/*

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o /js8web .

# ---- Runtime stage ----
FROM debian:bullseye-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /js8web /usr/local/bin/js8web

# Default database path inside the volume
ENV JS8WEB_DB_PATH=/data/js8web.db

# Database volume
VOLUME /data

EXPOSE 8080

ENTRYPOINT ["js8web"]
