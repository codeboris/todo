# Используем официальный образ Go
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum (если есть) для установки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod tidy

# Копируем весь исходный код
COPY . .

# Компилируем бинарный файл приложения
RUN go build -o main ./cmd/main.go

# Создаём минимальный образ для работы
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем скомпилированный бинарник из стадии сборки
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs/

# Открываем порт, на котором будет работать приложение
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
