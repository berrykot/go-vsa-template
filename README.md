# VSA

## Запуск (обязательный порядок)

Перед первым запуском установи Wire (один раз):

```bash
git config core.hooksPath .githooks

go install github.com/google/wire/cmd/wire@latest

go install github.com/go-task/task/v3/cmd/task@latest
```

Далее при любом изменении зависимостей или новых фичах — сначала генерируй Wire, потом запускай приложение:

```bash
wire ./cmd/app
go run ./cmd/app
```

Через Makefile (необязательно):

```bash
make wire      # только wire
make run       # wire + go run
make build     # wire + go build -o bin/app ./cmd/app
make all       # tidy + wire + build
```

---

## Стек

- **Go 1.25**
- **Gin** — HTTP API + middleware
- **Wire** — DI
- **Zerolog** — логирование
- **caarlos0/env** + **godotenv** — конфиг из env
- **robfig/cron/v3** — фоновые джобы
---

## 📁 Структура проекта (для LLM)

```
gestero-backend/
├── cmd/
│   └── app/
│       ├── main.go           # Точка входа
│       ├── app.go            # Оркестратор жизненного цикла (Router + Scheduler + фичи)
│       ├── wire.go           # Чертёж зависимостей (Wire)
│       └── wire_gen.go       # Автогенерируемый файл Wire
├── internal/
│   ├── config/               # Типизированная конфигурация (env → struct)
│   │   └── config.go
│   ├── infrastructure/
│   │   ├── logger/           # Zerolog
│   │   ├── database/         # любая бд
│   │   ├── auth/             # Auth middleware + UserMetaData
│   │   ├── server/           # Gin Router: Engine + Public/Protected группы
│   │   └── scheduler/        # Cron jobs (robfig/cron/v3)
│   └── features/             # Бизнес-логика (Vertical Slice Architecture)
│       └── health/           # Пример: Health check
├── .env.example              # Пример файла с переменными окружения
├── Dockerfile                # Multi-stage Docker build
├── Makefile                  # Команды для разработки
├── go.mod
├── go.sum
└── README.md
```

---

## Как добавлять новую фичу (VSA, без отхода от архитектуры)

**Для LLM:** ориентируйся на готовые примеры:
- HTTP + DB: `internal/features/health` (`Handler`, метод `Register`, запрос в Databese)
- HTTP + Auth + decimal: `internal/features/report` (`Handler`, `Generate`, использование `auth.GetAuthUser` и `decimal.Decimal`)

### 1. Папка фичи

Создай `internal/features/<feature-name>/` (например `internal/features/orders/`).

### 2. API (HTTP-хендлер) — public или protected

- Файл: `handler.go`.
- Структура-хендлер с зависимостями (через Wire), конструктор `NewHandler(...)`.
- Метод **`Register(g *gin.RouterGroup)`** — принимает **группу роутов**:
  - **Public** (`router.Public`) — без JWT (пример: `health.Register`).
  - **Protected** (`router.Protected`) — с Bearer JWT Databese (пример: `report.Handler.Register`).
- В **`cmd/app/wire.go`**: добавь конструктор фичи (`orders.NewHandler`) в `wire.Build(...)`.
- В **`cmd/app/app.go`**:
  - добавь хендлер в сигнатуру `NewApp(..., ordersH *orders.Handler, ...)`;
  - в теле `NewApp` вызови `ordersH.Register(router.Public)` или `ordersH.Register(router.Protected)` в зависимости от того, требует ли фича авторизации.

### 3. Cron (фоновые джобы)

- Файл: `job.go`.
- Структура-обработчик джобов с зависимостями (например `logger`, `*Databese.Client`, cron-строка из конфига), конструктор `NewJobHandler(...)`.
- Метод **`Register(s *scheduler.Scheduler)`** — внутри только регистрацию функций: `s.Cron.AddFunc(cronExpr, h.someJob)`.
- В **`internal/config/config.go`**: добавь поле в `Config.Cron` с тегом `env:"MY_JOB_CRON,required"` (или `envDefault`).
- В **`.env.example`**: добавь переменную `MY_JOB_CRON=...`.
- В **`cmd/app/wire.go`**: добавь `orders.NewJobHandler`.
- В **`cmd/app/app.go`**: добавь JobHandler в `NewApp` и вызови `ordersJob.Register(scheduler)`.

### 4. Конфиг из env

- Все настройки фичи — в **`internal/config/config.go`** (расширяем `Config` или вложенные структуры), теги `env:"VAR_NAME"` / `envDefault:"..."` / `required`.
- Примеры: `Port`, `Databese.URL`, `Databese.Key`, `Cron.QuarterlyReportCron`. В `.env.example` дублируй все переменные.

### 5. Wire после изменений

После любых изменений в `wire.go` или новых провайдерах фич выполни:

```bash
wire ./cmd/app
```

Иначе `go run`/`go build` могут использовать устаревший или отсутствующий `wire_gen.go`.

### Правила VSA в этом проекте

- Одна фича = одна папка в `internal/features/<name>/`.
- Внутри слайса — только код этой фичи (хендлеры, джобы, локальные типы). Общие вещи — в `internal/infrastructure` или `internal/config`.
- Роуты регистрируются только через `Register(*gin.RouterGroup)` в `app.go`, не в `infrastructure/server`.
- Зависимости фич инжектятся через конструкторы и Wire; в `wire.go` перечислены все провайдеры.


---

## Docker

Multi-stage: в builder ставится Wire, выполняется `wire ./cmd/app`, собирается бинарник; в финальном образе — только alpine + ca-certificates + бинарник. Порт 8080, переменная `PORT` для Render. См. `Dockerfile`.

---

## Env

Скопируй `.env.example` в `.env` и заполни. Обязательные для текущего кода: `Databese_URL`, `Databese_ANON_KEY`, `QUARTERLY_REPORT_CRON`. Остальные — см. `internal/config/config.go` и `.env.example`.
