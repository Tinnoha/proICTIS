# Этап сборки
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY Backend/Database/go.mod Backend/Database/go.sum ./
RUN go mod download

COPY Backend/Database/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY Frontend ./frontend

EXPOSE 8080

CMD ["./main"]