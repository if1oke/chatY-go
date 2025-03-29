# 🧠 ChatY-Go

Лёгкий, чистый TCP-чат на Go — как в старые добрые, но с архитектурой нового века.

## 🚀 Быстрый старт

```bash
go run cmd/server/main.go
```
По умолчанию сервер стартует по адресу и порту из .env (например, localhost:9000).

## 💬 Команды чата

Вводятся прямо в строку ввода. Начинаются с `/`

| Команда       | Описание                                        |
|---------------|-------------------------------------------------|
| `/nick <имя>` | Изменить свой никнейм                           |
| `/list`       | Посмотреть список всех активных пользователей   |
| `/exit`       | Выйти из чата                                   |
| `/help`       | Показать список команд                          |

### 🔄 Примеры:
```
/nick IgorKilla
/list
/exit
```

## 🔐 Поведение

- По умолчанию пользователю присваивается ник вида `User_127.0.0.1:XXXX`
- Все события (вход/выход/смена ника) — рассылаются остальным
- Все сообщения логируются и рассылаются в формате `[nickname]: сообщение`

## 💡 Архитектура
Проект разделён на слои:

* `cmd/server` — точка входа
* `internal/api/tcp` — TCP-сервер, принимает подключения
* `internal/domain/` — интерфейсы, изолированы
* `internal/application/session` — реализация сессии: регистрация, рассылка, чтение сообщений
* `internal/application/server` — DI-контейнер

Все зависимости внедряются через интерфейсы — чистая архитектура без лишнего.

## 📦 Текущий функционал
* TCP-сервер с поддержкой нескольких клиентов
* Асинхронная широковещательная рассылка сообщений
* Безопасный доступ к общим структурам через sync.Mutex
* Логика клиента изолирована внутри ChatSession
* Поддержка команд

## 🤝 Автор
Братан, который ебашил не в теории, а в бою 💪 Постигал основы DDD и чистоты в архитектуре.