# Game Integration API

A Go service for integrating third-party casino games with wallet and transaction management.

## Features
- Clean Architecture
- RESTful API
- PostgreSQL for persistence
- Wallet service integration
- Dockerized for easy setup

## Prerequisites
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Setup

1. **Clone the repository:**
   ```sh
   git clone <your-repo-url>
   cd GameIntegrationAPI
   ```

2. **Copy and configure environment variables:**
   ```sh
   cp .env.example .env
   # Edit .env as needed
   ```

3. **Build and run with Docker Compose:**
   ```sh
   docker-compose up --build
   ```

4. **Access services:**
   - Game Integration API: [http://localhost:8080](http://localhost:8080)
   - Wallet Service Swagger: [http://localhost:8000/swagger/index.html](http://localhost:8000/swagger/index.html)

## Development

- Main code: `cmd/main.go`
- Clean architecture: see `internal/`

## Database
- Default user: `gameuser`
- Default password: `gamepass`
- Default db: `gamedb`

## License
MIT 