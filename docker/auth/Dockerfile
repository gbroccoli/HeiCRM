# ---------- STAGE 1: build ----------
FROM golang:1.25 AS build

WORKDIR /src/services/auth

# сначала зависимости (для кэширования)
COPY ./services/auth/go.mod ./services/auth/go.sum ./
RUN go mod download

# копируем остальной код
COPY ./services/auth/ ./

# собираем main из cmd/auth
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/auth ./cmd/auth


# ---------- STAGE 2: runtime ----------
FROM ubuntu:24.04

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && useradd --create-home --shell /usr/sbin/nologin appuser \
    && mkdir -p /app/logs \
    && chown -R appuser:appuser /app

COPY ./config.yaml /app/config.yaml
COPY --from=build /out/auth /app/auth


EXPOSE 8080

USER appuser:appuser

ENTRYPOINT ["/app/auth"]
