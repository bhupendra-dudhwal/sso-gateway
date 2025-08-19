# sso-gateway
A centralized Single Sign-On (SSO) gateway using JSON Web Tokens (JWT) for secure, stateless authentication across multiple services and applications.


---

## ⚙️ Tech Stack

- **Language**: Go 1.24.6
- **Architecture**: Hexagonal / Clean Architecture
- **Database**: PostgreSQL 16
- **Cache**: KeyDB (Redis-compatible)
- **Configuration**: YAML-based
- **Containerization**: Docker & Docker Compose
- **Task Runner**: Makefile

---

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/bhupendra-dudhwal/sso-gateway.git
cd sso-gateway
```

### 2. Install Go

Download and install Go from https://go.dev/dl based on your OS.

### 3. Set Up Config File

Copy the sample config:

```bash
cp config/config.sample.yaml config/config.yaml
```
Edit config.yaml as needed.

## ▶️ Run with Docker (Full Setup)

## Build & Run Services

```bash
make build     # Build Docker images
make run       # Start the app, Postgres, and KeyDB
```

The service will be available at: http://localhost:8080

## 💻 Run Locally (App on Host, Dependencies in Docker)
If you want to run only the Go app on your machine and use Docker only for services like Postgres and KeyDB:

### 1. Start Dependencies Only
```bash
make deps
```
This will start PostgreSQL and KeyDB via Docker.

### 2. Run the App Locally
```bash
go run cmd/http/main.go
```
Ensure your config/config.yaml is correctly set up to match Docker hostnames:

```yaml
database:
  host: localhost
  port: 5432
  ...

redis:
  host: localhost
  port: 6379
  ...
```

## 🛠️ Makefile Commands

| Command              | Description                                        |
|----------------------|----------------------------------------------------|
| `make build`         | Build all services via Docker Compose              |
| `make run`           | Start the app along with Postgres and KeyDB        |
| `make stop`          | Stop all running containers                        |
| `make logs`          | Tail logs of the app container (`go-hexagonal`)   |
| `make clean`         | Remove containers, volumes, and networks           |
| `make db-psql`       | Open a psql shell inside the Postgres container    |
| `make redis-cli`     | Open a CLI session inside the KeyDB container      |
| `make deps`          | Start only Postgres and KeyDB                      |
| `make deps-down`     | Stop only Postgres and KeyDB                       |
| `make postgres-up`   | Start only the Postgres service                    |
| `make postgres-down` | Stop only the Postgres service                     |
| `make keydb-up`      | Start only the KeyDB service                       |
| `make keydb-down`    | Stop only the KeyDB service                        |


## 🧪 Health Check
To verify the service is running, you can hit the health endpoint:

```bash
curl http://localhost:8080/healthz/readiness
curl http://localhost:8080/healthz/liveness
```

## 🧼 Clean Up Everything
```bash
make clean
```
This will stop all services, remove volumes and orphan containers, and free up space.

## 📁 Project Structure
```plaintext
go-hexagonal/
├── cmd/
│   └── http/                  # Application entry point (main.go)
├── config/
│   ├── config.sample.yaml     # Sample config
│   └── config.yaml            # Actual config
├── docker-compose.yml         # Service definitions
├── Dockerfile                 # Multi-stage Go Dockerfile
├── go.mod / go.sum            # Go dependencies
├── makefile                   # Dev tasks and helpers

├── internal/                  # Application core
│   ├── builder/               # Object builders / constructors
│   ├── constants/             # App-wide constants and helpers
│   ├── core/
│   │   ├── models/            # Data models
│   │   ├── ports/             # Interfaces (ingress/egress)
│   │   └── services/          # Business logic
│   ├── egress/
│   │   ├── cache/             # Redis (KeyDB) connections
│   │   ├── database/          # PostgreSQL connections
│   │   └── repository/        # Repository implementations
│   ├── ingress/
│   │   ├── handler/           # HTTP request handlers
│   │   ├── middleware/        # Custom middleware
│   │   └── response/          # Response formatting
│   └── utils/                 # Utility functions

├── pkg/
│   └── logger/                # Custom logger package
└── README.md                  
```

## ❓ Troubleshooting
- Make sure Docker and Docker Compose are installed and running.
- Check if the ports 5432, 6379, and 8080 are available.
- For config issues, verify the config.yaml file paths and values.