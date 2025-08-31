# Архитектура и технологический стек форум-системы

## Обзор архитектуры

Данный проект представляет собой **микросервисную архитектуру** для форум-системы, построенную с использованием лучших практик разработки программного обеспечения 2025 года. Система состоит из нескольких независимых сервисов, каждый из которых отвечает за определенную бизнес-функцию.

## Технологический стек

### Основные технологии

| Компонент | Технология | Версия | Назначение |
|-----------|------------|---------|------------|
| **Язык программирования** | Go | 1.23 | Высокопроизводительный backend |
| **Web Framework** | Echo | v4.13.3 | HTTP API и middleware |
| **WebSocket** | Gorilla WebSocket | v1.5.3 | Реальное время соединения |
| **База данных** | PostgreSQL | latest | Основное хранилище данных |
| **DB Driver** | pgx | v5.7.2 | Высокопроизводительный PostgreSQL драйвер |
| **Message Broker** | RabbitMQ | 4-management | Асинхронная обработка сообщений |
| **Кэширование** | Redis | 7.0-alpine | Кэш и pub/sub |
| **Файловое хранилище** | SeaweedFS | latest | Распределенное хранение файлов |
| **Полнотекстовый поиск** | MeiliSearch | latest | Поиск по сообщениям |
| **Контейнеризация** | Docker | latest | Изоляция и развертывание |
| **Оркестрация** | Docker Compose | latest | Управление многоконтейнерными приложениями |
| **Обратный прокси** | Nginx | latest | Балансировка нагрузки и маршрутизация |

### Языки программирования

#### Backend
- **Go 1.23** - Основной язык для всех микросервисов
- **Python 3.12** - Сервис для обработки платежей и интеграций

## Архитектурные подходы и практики

### 1. Микросервисная архитектура (Microservices Architecture)

Система разделена на следующие сервисы:

- **forum-authentication** (порт 8000) - Аутентификация пользователей
- **forum-authorization** (порт 8001) - Авторизация и управление правами
- **forum-dialog** (порт 8002) - Система личных сообщений
- **forum-thread** (порт 8003) - Форумные темы и сообщения
- **forum-frame** (порт 8004) - Управление разделами форума
- **forum-status** (порт 8005) - Статусы пользователей
- **forum-user** (порт 8006) - Профили пользователей
- **forum-payment** (порт 8009) - Платежная система
- **forum-admin** (порт 8010) - Административная панель

### 2. Clean Architecture (Чистая архитектура)

Каждый микросервис следует принципам чистой архитектуры:

```
internal/
├── controller/          # Слой представления
│   ├── http/           # HTTP обработчики
│   ├── ws/             # WebSocket обработчики
│   └── amqp/           # AMQP обработчики
├── service/            # Бизнес-логика
├── repo/              # Слой данных
└── model/             # Модели данных
```

### 3. Domain-Driven Design (DDD)

- **Bounded Contexts** - каждый сервис представляет отдельный bounded context
- **Entities** - доменные сущности с четко определенными границами
- **Value Objects** - неизменяемые объекты-значения
- **Domain Services** - бизнес-логика, которая не принадлежит конкретной сущности

### 4. Event-Driven Architecture (EDA)

Система использует событийно-ориентированную архитектуру с RabbitMQ для асинхронной связи между сервисами:

- **Publishers** - сервисы публикуют события в очереди
- **Consumers** - сервисы подписываются на события и обрабатывают их
- **Event Sourcing** - некоторые операции записываются как последовательность событий

### 5. CQRS (Command Query Responsibility Segregation)

Разделение операций чтения и записи:
- **Commands** - операции изменения состояния
- **Queries** - операции чтения данных
- **Разные модели** для записи и чтения данных

## Лучшие практики разработки

### 1. Dependency Injection

Использование внедрения зависимостей для тестируемости и гибкости:

```go
type ServiceUser struct {
    userRepo           model.IUserRepo
    notificationClient model.INotificationClient
}

func New(
    userRepo model.IUserRepo,
    notificationClient model.INotificationClient,
) *ServiceUser {
    return &ServiceUser{
        userRepo:           userRepo,
        notificationClient: notificationClient,
    }
}
```

### 2. Interface Segregation

Использование интерфейсов для разделения ответственности:

```go
type IUserService interface {
    CreateUser(ctx context.Context, accountID int, login string) (int, error)
    UserByAccountID(ctx context.Context, accountID int) (*User, error)
}
```

### 3. Repository Pattern

Абстракция слоя данных:

```go
type IUserRepo interface {
    CreateUser(ctx context.Context, accountID int, login string) (int, error)
    UserByAccountID(ctx context.Context, accountID int) ([]*User, error)
}
```

### 4. Transaction Management

Управление транзакциями с контекстом:

```go
ctx, err := s.messageRepo.CtxWithTx(ctx)
if err != nil {
    return 0, err
}
defer s.messageRepo.RollbackTx(ctx)

// Бизнес-операции...

err = s.messageRepo.CommitTx(ctx)
```

### 5. Circuit Breaker Pattern

Обработка отказов внешних сервисов для повышения отказоустойчивости.

## Паттерны обмена данными

### 1. Request-Response

Синхронная связь через HTTP API между сервисами.

### 2. Publish-Subscribe

Асинхронная связь через RabbitMQ для событийной архитектуры:

```go
func (r *RabbitMQ) Publish(ctx context.Context, queueName string, body any) error {
    event := Event{Body: body}
    jsonData, _ := json.Marshal(event)
    
    return r.channel.PublishWithContext(ctx, "", queueName, false, false, 
        amqp091.Publishing{
            ContentType:  "application/json",
            Body:         jsonData,
            DeliveryMode: amqp091.Persistent,
        })
}
```

### 3. WebSocket для реального времени

Двунаправленная связь для чатов и уведомлений:

```go
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin:     func(r *http.Request) bool { return true },
}
```

## Управление данными

### 1. Database per Service

Каждый микросервис имеет собственную базу данных:
- **Изоляция данных** между сервисами
- **Технологическая независимость** выбора БД
- **Масштабируемость** каждого сервиса отдельно

### 2. Connection Pooling

Использование pgx с пулом соединений для оптимальной производительности:

```go
func New(username, password, host, port, dbname string) *Postgres {
    connString := "postgres://" + username + ":" + password + 
                  "@" + host + ":" + port + "/" + dbname + 
                  "?pool_max_conns=100"
    
    db, err := pgxpool.New(ctx, connString)
    return &Postgres{db}
}
```

### 3. Database Migrations

Автоматические миграции схемы БД:

```sql
CREATE TABLE IF NOT EXISTS threads (
    id SERIAL PRIMARY KEY,
    thread_name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

## Безопасность

### 1. JWT Authentication

Аутентификация на основе JWT токенов:

```go
type AuthorizationData struct {
    AccountID   int    `json:"accountID"`
    Role        string `json:"role"`
    TwoFaStatus bool   `json:"twoFaStatus"`
    Message     string `json:"message"`
    Code        int    `json:"code"`
}
```

### 2. Role-Based Access Control (RBAC)

Система ролей и разрешений:

```go
var StatusPermissionMap = map[int]StatusPermission{
    1: {
        PrivateThreads:       []string{"thread1", "thread2"},
        MaxTopicsPerDay:      2,
        MaxTopicsInSubthread: 10,
    },
}
```

### 3. Input Validation

Валидация входных данных:

```go
type CreateUserBody struct {
    AccountID int    `json:"accountID" validate:"required"`
    Login     string `json:"login"`
}
```

## Мониторинг и наблюдаемость

### 1. Structured Logging

Структурированное логирование с контекстом:

```go
logger := slog.With(
    slog.String("request_id", requestID),
    slog.String("method", request.Request().Method),
    slog.String("service_name", "forum-dialog"),
)
```

### 2. Health Checks

Проверки состояния сервисов для мониторинга.

### 3. Metrics Collection

Сбор метрик производительности и бизнес-метрик.

## Тестирование

### 1. Unit Testing

Изолированное тестирование компонентов:

```go
func TestControllerCreateUser(t *testing.T) {
    testConfig.PrepareDB()
    
    err := testConfig.userClient.CreateUser(1, "123456")
    assert.NoError(t, err)
}
```

### 2. Integration Testing

Интеграционное тестирование с Docker Compose:

```yaml
services:
  forum-user-unit-test:
    build:
      context: "../../../"
      dockerfile: ".github/test/unit/Dockerfile"
    depends_on:
      - forum-user-db-unit-test
```

### 3. API Testing

Тестирование HTTP API и WebSocket соединений.

## DevOps и развертывание

### 1. Containerization

Каждый сервис упакован в Docker контейнер:

```dockerfile
FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/main.go

FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/server .
CMD ["./server"]
```

### 2. Multi-stage Builds

Оптимизированные Docker образы с многоэтапной сборкой.

### 3. CI/CD Pipeline

Автоматизированная сборка и развертывание:

```yaml
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: self-hosted
    steps:
      - uses: actions/checkout@v2
      - name: Update production
        run: |
          ssh $USER@$SERVER << EOF
          cd wewall/wewall-system
          git pull
          EOF
```

### 4. Infrastructure as Code

Описание инфраструктуры в виде кода с Docker Compose:

```yaml
services:
  forum-authentication:
    build:
      context: "../forum-authentication"
      dockerfile: ".github/Dockerfile"
    ports:
      - "8000:8000"
    depends_on:
      - forum-authentication-db
```

## Масштабируемость

### 1. Horizontal Scaling

Горизонтальное масштабирование сервисов через репликацию контейнеров.

### 2. Load Balancing

Балансировка нагрузки через Nginx:

```nginx
upstream backend {
    server forum-service1:8000;
    server forum-service2:8000;
}
```

### 3. Caching Strategy

Кэширование с Redis для улучшения производительности:

```go
func (rc *ClientRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
    return rc.client.Set(ctx, key, value, expiration).Err()
}
```

## Производительность

### 1. Connection Pooling

Пулы соединений для всех внешних ресурсов (БД, Redis, RabbitMQ).

### 2. Lazy Loading

Отложенная загрузка данных для оптимизации памяти.

### 3. Async Processing

Асинхронная обработка тяжелых операций через очереди сообщений.

## Заключение

Данная архитектура представляет собой современный подход к построению масштабируемых, отказоустойчивых и легко поддерживаемых систем. Использование микросервисной архитектуры в сочетании с проверенными паттернами и практиками позволяет создать высокопроизводительную форум-систему, готовую к продакшн-нагрузкам.

Ключевые преимущества данного подхода:
- **Модульность** - независимая разработка и развертывание сервисов
- **Масштабируемость** - возможность масштабирования отдельных компонентов
- **Отказоустойчивость** - изоляция отказов между сервисами
- **Технологическое разнообразие** - возможность выбора оптимальных технологий для каждой задачи
- **Командная автономность** - разные команды могут работать над разными сервисами