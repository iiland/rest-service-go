# Используем официальный golang образ с alpine для сборки
FROM golang:1.21-alpine AS builder

# Устанавливаем переменные окружения для компиляции статического бинарника
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Копируем файлы зависимостей и скачиваем модули
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN go build -o main .

# Финальный минимальный образ на alpine, добавляем ca-certificates
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder
COPY --from=builder /app/main .

# Открываем порт HTTP на 8080 
EXPOSE 8080

# Запускаем бинарник
CMD ["./main"]
