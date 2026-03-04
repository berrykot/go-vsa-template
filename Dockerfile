# --- Stage 1: Зависимости + код + генерация (wire) ---
FROM golang:1.25-alpine AS generator

RUN go install github.com/go-task/task/v3/cmd/task@v3.48.0

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux
RUN task generate

FROM generator AS tester
RUN task test

FROM tester AS builder
RUN task build

FROM alpine:3.23.3

# Сертификаты (HTTPS) + tzdata (time.LoadLocation, напр. Europe/Madrid для cron)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /

# Копируем скомпилированный файл из первого этапа
COPY --from=builder /build/bin/app /usr/local/bin/app

#The default value of PORT is 10000 for all Render web services.
EXPOSE 8080

# Запускаем приложение
CMD ["/usr/local/bin/app"]