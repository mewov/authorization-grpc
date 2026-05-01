# Authorization Service (gRPC) 🔐

**Authorization Service** — микросервис авторизации и управления пользователями, написанный на **Golang (gRPC)**.
Он обеспечивает **регистрацию**, **вход**, **обновление токенов**, **валидацию** и **логаут**, используя **PostgreSQL** для хранения данных.

**Authorization Service** is a microservice for user authentication and session management, built with **Go (gRPC)**.
It supports **user registration**, **login**, **token refresh**, **token validation**, and **logout**, backed by **PostgreSQL**.

---

## 🧩 Core Features / Основные возможности

| Feature / Функция           | Description (EN)                                          | Описание (RU)                                              |
| --------------------------- | --------------------------------------------------------- | ---------------------------------------------------------- |
| **User Registration**       | Create user with login, password, email, role, and client | Регистрация пользователя с логином, паролем, email и ролью |
| **Login**                   | Authenticate user and issue JWT tokens                    | Авторизация пользователя и выдача JWT                      |
| **Access / Refresh Tokens** | Short-lived access + long-lived refresh pair              | Короткоживущий access и долговечный refresh токен          |
| **Token Refresh**           | Obtain new tokens via refresh request                     | Получение новой пары токенов по refresh                    |
| **User Info**               | Get user info from access token                           | Получение данных пользователя из access токена             |
| **Logout**                  | Revoke refresh token and terminate session                | Выход и удаление активной сессии                           |
| **Structured Logging**      | Per-request logging via custom slog-based logger          | Логирование каждого вызова с контекстом                    |

---

## ⚙️ Technologies / Технологии

| Layer              | Stack                                      |
| ------------------ | ------------------------------------------ |
| **Language**       | Go 1.25+                                   |
| **Protocol**       | gRPC (proto3)                              |
| **Database**       | PostgreSQL                                 |
| **Auth**           | Access / Refresh                           |
| **Middleware**     | RateLimiter, Logger                        |
| **Sessions**       | Stored in PostgreSQL via `SessionsService` |
| **Config**         | Custom `LocalConfig` structure             |
| **Docker Compose** | Used for container orchestration           |

---

## 🛰️ gRPC Services / Сервисы

### **Authorization**

Обрабатывает регистрацию и вход пользователей.

| RPC Method | Request Message   | Response Message | Description (EN)                      | Описание (RU)                             |
| ---------- | ----------------- | ---------------- | ------------------------------------- | ----------------------------------------- |
| `Register` | `RequestRegister` | `ResponseToken`  | Register a new user and return tokens | Регистрация пользователя и выдача токенов |
| `Login`    | `RequestLogin`    | `ResponseToken`  | Login user and return tokens          | Авторизация и получение токенов           |

---

### **Session**

Управляет сессиями, обновлением и валидацией токенов.

| RPC Method | Request Message  | Response Message | Description (EN)                | Описание (RU)                                        |
| ---------- | ---------------- | ---------------- | ------------------------------- | ---------------------------------------------------- |
| `Refresh`  | `RequestRefresh` | `ResponseToken`  | Refresh token pair              | Обновление пары токенов                              |
| `Logout`   | `RequestLogout`  | `ResponseLogout` | Revoke refresh token            | Удаление refresh токена и завершение сессии          |
| `Info`     | `RequstInfo`     | `ResponseInfo`   | Get user info from access token | Получение информации о пользователе из access токена |

---

## 🧾 Message Definitions / Структуры сообщений

### Request / Response Types

```proto
message RequestRegister {
    string login   = 1;
    string password = 2;
    string email   = 3;
    string client  = 4;
    string role    = 5;
}

message RequestLogin {
    string login   = 1;
    string password = 2;
}

message RequestRefresh {
    string refresh = 1;
}

message Tokens {
    string refresh = 1;
    string access  = 2;
}

message ResponseToken {
    string status  = 1;
    string message = 2;
    Tokens tokens  = 3;
}

message RequestLogout {
    string refresh = 1;
}

message ResponseLogout {
    string status  = 1;
    string message = 2;
}

message RequstInfo {
    string access = 1;
}

message ResponseInfo {
    string status  = 1;
    string message = 2;
    User user      = 3;
}

message User {
    int64  user_id = 1; 
    string login   = 2;
    string email   = 3;
    string role    = 4;
}
```

---

## 🔒 Authentication Flow / Процесс авторизации

1. **Register** — клиент вызывает `Authorization.Register` с логином, паролем, email, client и role → получает `access` и `refresh` токены.
2. **Login** — вызывает `Authorization.Login` → получает новую пару токенов.
3. **Access Token** — короткоживущий (примерно 15 мин), используется для доступа.
4. **Refresh Token** — живёт дольше (7 дней), хранится в таблице `sessions`.
5. **Refresh** — вызывает `Session.Refresh` с refresh токеном → получает новую пару.
6. **Info** — вызывает `Session.Info` с access токеном → получает информацию о пользователе.
7. **Logout** — вызывает `Session.Logout` → refresh токен удаляется, сессия завершается.

---

## ⚡ Quick Start / Быстрый старт

### Local

```bash
go run cmd/server.go
```

### Docker

```bash
docker-compose up --build
```

### Default Address

```
localhost:44044
```

---

## 📄 License / Лицензия

**MIT License** — свободное использование и модификация проекта.
