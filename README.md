# Go-Gin API Server

This project is a simple API server built with [Go](https://golang.org/) using the [Gin Web Framework](https://github.com/gin-gonic/gin). The project supports live reloading during development using **Air** and has an integrated migration system.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Development](#development)
- [Database Migrations](#database-migrations)
  - [Apply Migrations](#apply-migrations)
  - [Rollback Migrations](#rollback-migrations)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

---

## Prerequisites

Before setting up the project, ensure you have the following installed:

- [Go 1.17+](https://golang.org/dl/)
- [PostgreSQL](https://www.postgresql.org/) (or your preferred database)
- [Air](https://github.com/air-verse/air) for live-reloading (installed as a dependency)

## Installation

1. **Clone the repository**:
  ```bash
  git clone https://github.com/yourusername/go-gin-api-server.git
  cd go-gin-api-server
  ```

2. **Install dependencies**: Ensure all Go dependencies and tools (including **Air**) are installed.
  ```bash
  go mod tidy
  go install github.com/air-verse/air@latest
  ```

3. **Setup environment variables**: Create a `.env` file with your database and other environment variables.
  ```bash
  cp .env.example .env
  ```

  Example `.env` file:
  ```ini
  # Available modes: debug, release or test
  GIN_MODE=debug

  # The port that the application will run on
  APP_PORT=3000

  # Database configuration
  DB_HOST=localhost
  DB_PORT=5432
  DB_USER=postgres
  DB_PASSWORD=your_password
  DB_NAME=postgres
  DB_SSLMODE=disable
  ```

## Development

To start the server in development mode with live-reloading (using **Air**):

```bash
make develop
```

## Database Migrations

The project includes a migration system for managing database schema changes:

### Apply Migrations

To apply all pending migrations, run:

```bash
make migrate-up
```

### Rollback Migrations

To rollback the last migration, run:

```bash
make migrate-down
```

Migration files are located in the `database/migrations/` directory.
- **Up Migration**: Files ending in `.up.sql` are used for applying changes.
- **Down Migration**: Files ending in `.down.sql` are used for rolling back changes.

## Configuration

### air.toml (Development Mode)

The `air.toml` file is used for configuring the **Air** live-reloading tool. It watches specific directories and file types, such as `.go` and `.html`, to automatically rebuild and restart the server during development.

You can modify `air.toml` to suit your development workflow, for example:

```toml
[build]
  cmd = "go build -o ./tmp/main ."
  bin = "./tmp/main"
  delay = 1000
  tmp_dir = "tmp"

[watch]
  includes = ["./controllers", "./routes", "./services", "./validators", "./config"]
  include_ext = ["go", "html", "tmpl", "tpl"]
  exclude_dir = ["vendor", "tmp", "database/migrations"]

[log]
  level = "info"
```

## Contributing

If you would like to contribute to the project:

1. Fork the repository.
2. Create a new branch (git checkout -b feature-branch).
3. Commit your changes (git commit -am 'Add new feature').
4. Push to the branch (git push origin feature-branch).
5. Create a new Pull Request.

All contributions are welcome!

## License

This project is open-source and available under the [MIT License](LICENSE).

