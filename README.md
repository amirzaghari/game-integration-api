# Game Integration API

A Game Integration API for casino games with wallet management. Provides authentication, player information, bet placement (withdraw), bet settlement (deposit), and transaction cancellation endpoints.

## Prerequisites
- [Docker](https://www.docker.com/)
- [docker-compose](https://docs.docker.com/compose/)

## Setup & Running

1. **Start all services (Postgres, Wallet, App):**
   ```sh
   make up
   ```
   This will start the database, wallet service, and the app using Docker Compose.

2. **Build the app (if needed):**
   ```sh
   make build
   ```
   This will build the app container and copy the .env.example to .env if present.

3. **Run the app (if already built):**
   ```sh
   make run
   ```
   This will run the app inside the container.

4. **Stop all services:**
   ```sh
   make down
   ```

## API Documentation

Once the app is running, access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

- The Swagger UI provides interactive documentation for all endpoints.
- The `/auth/login` endpoint can be tested with the following sample credentials (from the seeded users):
  ```json
  {
    "username": "testuser1",
    "password": "testpass"
}
```
- For protected endpoints, use the "Authorize" button and paste the JWT token returned from `/auth/login` as:
  ```
  Bearer <your_token>
```

## Running Tests

To run all tests in the `test/` directory inside the Docker container:

```
make test
```

---

- Environment variables are managed via Docker Compose and `.env` files.
- The app will auto-migrate and seed the database on startup.
- For local development with hot reload, you can use `make local-dev` (requires [air](https://github.com/cosmtrek/air)).
