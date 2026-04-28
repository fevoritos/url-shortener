# URL-Shortener Service

Сервис для сокращения ссылок. Поддерживает работу с PostgreSQL и хранение данных в памяти (In-Memory).

## Способы запуска

### 1. Docker Compose

```bash
docker-compose up --build
```

После запуска сервис будет доступен по адресу: `http://localhost:8080`

### 2. Запуск локально через параметры (флаги/env)

**Доступные параметры:**

- `-addr`: Адрес и порт (по умолчанию `:8080`).
- `-db-dsn`: DSN.
- `-storage`: Тип хранилища (`postgres` или `memory`).

**Пример запуска с использованием PostgreSQL:**

```bash
go run cmd/api/main.go -storage="postgres"
```

**Пример запуска в режиме In-Memory:**

```bash
go run cmd/api/main.go -storage="memory"
```

> **!** Приложение также поддерживает переменные окружения (`HTTP_ADDR`, `DATABASE_DSN`, `STORAGE_TYPE`).

---

## Документация API (Swagger)

`/docs/index.html`

### Перегенерация документации

`Makefile`

```bash
make swagger
```

---
