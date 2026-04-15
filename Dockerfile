# Этап сборки
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum из подпапки Backend/Database
COPY Backend/Database/go.mod Backend/Database/go.sum ./
RUN go mod download

# Копируем весь исходный код бэкенда
COPY Backend/Database/ .

# Собираем main.go из cmd/
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/main .

# Копируем статику фронтенда (папка Frontend в корне проекта)
COPY Frontend ./frontend

EXPOSE 8080

CMD ["./main"]