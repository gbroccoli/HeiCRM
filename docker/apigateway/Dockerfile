# ---------- STAGE 1: build ----------
FROM golang:1.25 AS build

WORKDIR /src/services/apigateway

# сначала зависимости (для кэширования)
COPY ./services/apigateway/go.mod services/apigateway/go.sum ./
RUN go mod download

# копируем остальной код
COPY ./services/apigateway/ ./

# собираем main из cmd/apigateway
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/apigateway ./cmd/gateway


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
COPY --from=build /out/apigateway /app/apigateway

EXPOSE 8000

USER appuser:appuser

ENTRYPOINT ["/app/apigateway"]
