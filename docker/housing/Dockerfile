# ---------- STAGE 1: build ----------
FROM golang:1.25 AS build

WORKDIR /src

# сначала зависимости (для кэширования)
COPY ./services/housing/go.mod ./services/housing/go.sum ./services/housing/

# копируем локальные пакеты
COPY ./pkg/ ./pkg/

WORKDIR /src/services/housing
RUN go mod download

# копируем остальной код сервиса
COPY ./services/housing/ ./

# собираем
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/housing ./cmd/housing


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
COPY --from=build /out/housing /app/housing


EXPOSE 8082

USER appuser:appuser

ENTRYPOINT ["/app/housing"]
