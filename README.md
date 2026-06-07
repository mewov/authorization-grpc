# Authorization Service (gRPC) 🔐

A user authentication and session management microservice built with **Go (gRPC)** and backed by **PostgreSQL**.

---

## ✨ Core Features & gRPC API

### Authorization Service (Auth & Signup)

* `Register` (`RequestRegister` ➡️ `ResponseToken`) — Registers a new user (login, password, email, role, client) and issues tokens.
* `Login` (`RequestLogin` ➡️ `ResponseToken`) — Authenticates a user and returns a token pair.

### Session Service (Session Management)

* `Refresh` (`RequestRefresh` ➡️ `ResponseToken`) — Renews the access/refresh token pair using a valid refresh token.
* `Info` (`RequstInfo` ➡️ `ResponseInfo`) — Retrieves user profile data from the access token.
* `Logout` (`RequestLogout` ➡️ `ResponseLogout`) — Terminates the session and revokes the refresh token.

---

## ⚙️ Tech Stack

* **Language:** Go 1.25+
* **Protocol:** gRPC (proto3)
* **Database:** PostgreSQL (stores users and active sessions)
* **Middleware:** RateLimiter, structured logging (`slog`)
* **Orchestration:** Docker Compose

---

## 🧾 Message Definitions (Protobuf)

```proto
message RequestRegister { string login = 1; string password = 2; string email = 3; string client = 4; string role = 5; }
message RequestLogin    { string login = 1; string password = 2; }
message RequestRefresh  { string refresh = 1; }
message RequestLogout   { string refresh = 1; }
message RequstInfo      { string access = 1; }

message Tokens          { string refresh = 1; string access = 2; }
message ResponseToken   { string status = 1; string message = 2; Tokens tokens = 3; }
message ResponseLogout  { string status = 1; string message = 2; }
message ResponseInfo    { string status = 1; string message = 2; User user = 3; }
message User            { int64 user_id = 1; string login = 2; string email = 3; string role = 4; }

```

---

## ⚡ Quick Start

Default service address: `localhost:44044`

**Run Locally:**

```bash
go run cmd/server.go

```

**Run via Docker:**

```bash
docker-compose up --build

```

---

## 📄 License

**MIT License** — free to use and modify.
