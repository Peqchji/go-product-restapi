# Golang Product API

REST API for product, built with Go (Echo), PostgreSQL, and Docker.

## Prerequisites

- [Docker](https://www.docker.com/) & Docker Compose
- [Go 1.25+](https://go.dev/) (optional, for local development)

## Getting Started

1. **Clone the repository**
   ```bash
   git clone https://github.com/Peqchji/go-product-restapi.git
   cd go-product-restapi
   ```

2. **Configure Environment**
   Copy the example environment file:
   ```bash
   cp example.env .env
   ```
   Update `.env` values if necessary.

## Running the Service

### Using Docker (Recommended)

Start all services (API, Database, Migrations, PgAdmin):
```bash
docker-compose up -d --build
```

- **API URL**: `http://localhost:8080`
- **Health Check**: `http://localhost:8080/health`
- **Swagger Docs**: `http://localhost:8080/api-docs/index.html`
- **PgAdmin**: `http://localhost:5050`

To stop the services:
```bash
docker-compose down
```

### Running Locally

If you prefer to run the Go application locally:

1. **Start dependencies (PostgreSQL)**:
   ```bash
   docker-compose up -d postgres
   ```

2. **Run Migrations**:
   ```bash
   # Using golang-migrate CLI or running the migrate container
   docker-compose up migrate
   ```

3. **Run the Application**:
   ```bash
   go run cmd/httpserver/main.go
   ```

## Running Tests

Run the full test suite (Unit, Integration, Component):
```bash
go test -v ./...
```


## API Input Constraints

### Create Product
- **Name**: Required and can be empty string.
- **Price**: Must be greater than or equal to 0.
- **Sale Price**: Must be greater than or equal to 0 (if provided).

### Update Product (Patch)
- **Name**: Optional and can be empty string.
- **Price**: Must be greater than or equal to 0 (if provided).
- **Sale Price**: Must be greater than or equal to 0 (if provided).

