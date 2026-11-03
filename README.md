# SecAudit Hub

Локальная платформа непрерывного AppSec. В этом MVP работают фундамент, HttpOnly-аутентификация и промышленный модуль SAST. DAST, IaC и DB Security оставлены как архитектурные заглушки.

## Стек

- Frontend: React, Vite, TypeScript, TanStack Router, TanStack Query, Axios, Radix UI, SCSS Modules
- Backend: Go, chi, GORM, PostgreSQL через драйвер pgx, JWT в HttpOnly cookie, bcrypt
- SAST: 12 правил (CWE/OWASP), демо-корпус, вставка кода и shallow-clone публичного Git

## Запуск

```bash
make up
make dev
```

- Интерфейс: http://localhost:5173
- API: http://localhost:8080/api/health
- Postgres: `localhost:5432`, база `secaudit`, пользователь и пароль `postgres`

На экране входа есть «Быстрый вход демо-аудитора». JWT не попадает в `localStorage` или `sessionStorage`: сервер кладёт его в cookie `secaudit_session` (`HttpOnly`, `SameSite=Strict`). Vite проксирует `/api` на Go, поэтому cookie остаётся на одном origin.

`APP_ENV=production` включает `Secure` у cookie. Для локальной разработки флаг выключен.

## Проверки

```bash
cd backend && go test ./... && go build ./...
cd frontend && npm run build
```
