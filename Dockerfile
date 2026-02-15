# --- Stage 1: Build ---
FROM golang:1.25-alpine AS builder

# Устанавливаем wire для генерации зависимостей
RUN go install github.com/google/wire/cmd/wire@v0.7.0

WORKDIR /app

# Сначала копируем только файлы зависимостей (для кэширования слоев)
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . .

# Генерируем wire_gen.go прямо внутри контейнера
RUN wire ./cmd/app

# Собираем статический бинарник (CGO_ENABLED=0 критично для alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/app

# --- Stage 2: Final ---
FROM alpine:3.23.3

# Добавляем сертификаты (нужны для запросов к Supabase по HTTPS)
RUN apk --no-cache add ca-certificates

WORKDIR /

# Копируем только скомпилированный файл из первого этапа
COPY --from=builder /main /main

# Render сам подставит нужный порт в переменную PORT
EXPOSE 8080

# Запускаем приложение
CMD ["/main"]