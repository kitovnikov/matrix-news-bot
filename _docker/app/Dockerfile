FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o httpApp .


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/httpApp .

EXPOSE 8083

CMD ["./httpApp"]