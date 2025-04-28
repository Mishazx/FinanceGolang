# Этап сборки
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY ./src .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# Финальный этап
FROM alpine:latest

WORKDIR /app

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/main .
COPY --from=builder /app/bank_config.json .
COPY --from=builder /app/private_key.asc .
COPY --from=builder /app/public_key.asc .

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates

EXPOSE 8080

CMD ["./main"] 