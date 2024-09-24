# Go-Gin API Server

This project is a simple API server built with [Go](https://golang.org/) using the [Gin Web Framework](https://github.com/gin-gonic/gin).

The project supports live reloading during development using **Air**, has an integrated migration system and cron scheduling system.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Development](#development)
- [Project Structure](#project-structure)
  - [Explanation of Directories](#explanation-of-directories)
- [Database Migrations](#database-migrations)
  - [Apply Migrations](#apply-migrations)
  - [Rollback Migrations](#rollback-migrations)
- [Scheduled Tasks (Cron Jobs)](#scheduled-tasks-cron-jobs)
  - [Running Cron Jobs](#running-cron-jobs)
  - [Adding a New Cron Job](#adding-a-new-cron-job)
- [Validators](#validators)
  - [RegisterUserValidator](#registeruservalidator)
  - [LoginUserValidator](#loginuservalidator)
  - [Custom Error Messages](#custom-error-messages)
- [Auto-Generated Code](#auto-generated-code)
  - [Using `go generate`](#using-go-generate)
- [Configuration](#configuration)
- [Development Workflow](#development-workflow)
  - [Live Reloading with Air](#live-reloading-with-air)
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

## Project Structure

```bash
.
├── main.go                   # Entry point of the application
├── tools.go                  # Track tools used in the project
├── Makefile                  # Makefile for common tasks
├── air.toml                  # Air configuration for live-reloading
├── go.mod                    # Go module file
├── go.sum                    # Go module file
├── cmd
│   ├── generate_validators   # Tool to auto-generate validators
│   │   └── main.go
│   ├── migrate               # Tool to run database migrations
│   │   └── main.go
│   └── schedules             # Tool to register and run cron jobs
│       └── main.go
├── schedules                 # Package for cron tasks
│   └── tasks.go              # Decoupled task logic for cron jobs
├── config                    # Configuration files
│   └── database.go
├── controllers               # API route handlers
│   ├── auth_controller.go
│   └── item_controller.go
├── database                  # Database-related code
│   ├── migrate.go            # Migration logic
│   └── migrations            # SQL migration files
│       ├── 000001_create_items_table.up.sql
│       ├── 000001_create_items_table.down.sql
│       └── 000002_create_users_table.up.sql
│       ├── 000002_create_users_table.down.sql
├── middlewares               # Middleware logic
│   └── error_handler.go
├── models                    # Data models
│   ├── item.go
│   └── user.go
├── repositories              # Data access layer
│   ├── item_repository.go
│   └── user_repository.go
├── routes                    # API routes
│   └── routes.go
├── services                  # Business logic
│   ├── auth_service.go
│   └── item_service.go
├── tmp                       # Temporary files (excluded from version control)
│   └── main
├── utils                     # Utility functions
│   └── password.go
└── validators                # Input validation logic
    ├── auth_validator.go
    ├── item_validator.go
    ├── auto_generated.go     # Auto-generated file
    └── register.go           # Handles go:generate directive
```

### Explanation of Directories

- `cmd/`: This directory contains subdirectories for command-line tools. Currently, there are two:
  - `generate_validators`: Contains `main.go`, which is responsible for auto-generating the validator registration.
  - `migrate`: Contains `main.go`, which handles database migration commands such as `migrate-up` and `migrate-down`.
  - `schedules`: Contains `main.go`, which is responsible for registering and running cron jobs.

- `config/`: Contains configuration-related files, such as `database.go`, which is responsible for initializing the database connection.

- `controllers/`: This directory contains the handlers for your API endpoints. Each file corresponds to a different part of the API:
  - `auth_controller.go`: Handles authentication-related API routes (e.g., login, register).
  - `item_controller.go`: Handles item-related routes (e.g., CRUD operations for items).

- `database/`: Contains all database-related code:
  - `migrate.go`: The migration logic, handling applying and rolling back migrations.
  - `migrations/`: Directory containing SQL migration files, including both `.up.sql` (for applying migrations) and `.down.sql` (for rolling back).

- `middlewares/`: This directory contains middleware logic, such as `error_handler.go`, which is responsible for handling validation and binding errors.

- `models/`: Defines the data models for the application:
  - `item.go`: Defines the structure for the `Item` model.
  - `user.go`: Defines the structure for the `User` model.

- `repositories/`: Contains the data access layer, which abstracts database queries for different models:
  - `item_repository.go`: Provides the database access methods for the `Item` model.
  - `user_repository.go`: Provides the database access methods for the `User` model.

- `routes/`: Responsible for setting up the API routes:
  - `routes.go`: Contains the function that configures all the routes for the application.

- `services/`: This directory contains the business logic of the application:
  - `auth_service.go`: Contains the logic for user authentication, such as login and registration.
  - `item_service.go`: Contains the business logic for managing items.

- `tmp/`: Temporary files created during development, such as the Go binary generated by Air for live-reloading. This directory is excluded from version control.

- `tools.go`: A Go file that is used to track tools like Air. This file ensures development tools are included in `go.mod` and can be installed by others working on the project.

- `utils/`: This directory contains utility functions, such as `password.go`, which includes password hashing and validation logic.

- `validators/`: This directory contains all the input validation logic for your application:
  - `auth_validator.go`: Defines validators for authentication-related data (e.g., login, register).
  - `item_validator.go`: Defines validators for item-related data (e.g., item creation).
  - `auto_generated.go`: This file is auto-generated and contains dynamic registration of validators.
  - `register.go`: Handles the go:generate directive for generating the auto_generated.go file.

- `schedules/`: Contains the logic for scheduling and running cron jobs:
  - `tasks.go`: Contains the decoupled task logic for cron jobs.

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

## Scheduled Tasks (Cron Jobs)

The project uses cron jobs to perform scheduled tasks such as database cleanups or other recurring operations.

### Running Cron Jobs

Cron jobs are defined separately from the server and run independently. To start the cron job scheduler, run:

```bash
make run-cron
```

This will start the cron scheduler and execute tasks based on their defined schedule.

### Adding a New Cron Job

1. **Define the task** in the `schedules/tasks.go`. For example:
```go
func NewTask(db *sql.DB) {
  log.Println("Running a new task...")
  // Perform the task here
}
```

2. **Register the task** in the `cmd/schedules/main.go` with the cron expression:
```go
_, err := c.AddFunc("0 0 * * *", func() {
  tasks.NewTask(db)
})
```

**Example Cron Expression**:
- `* * * * *`: Every minute
- `0 0 0 * * *`: Runs at midnight every day

## Validators

Validators ensure that incoming data (such as user input) meets the necessary requirements before it is processed by the server. The project uses `go-playground/validator` to handle validation.

Validators are defined in the `validators/` directory.

### RegisterUserValidator

This validator is used when creating a new user (e.g., during registration). It ensures that the `username`, `email`, and `password` fields meet certain criteria.

**Location**: `validators/auth_validator.go`

```go
type RegisterUserValidator struct {
	Username string `json:"username" binding:"required,min=3,max=50" message:"Username is required with minimum of 3 characters and maximum of 50 characters"`
	Email    string `json:"email" binding:"required,email" message:"Email is required and must be a valid email address"`
	Password string `json:"password" binding:"required,min=6" message:"Password is required with minimum of 6 characters"`
}
```

- `Username`: Required, with a minimum length of 3 and a maximum of 50 characters.
- `Email`: Required, must be a valid email format.
- `Password`: Required, with a minimum length of 6 characters.

### LoginUserValidator

This validator is used when logging in a user. It ensures that the `email` and `password` fields are provided.

**Location**: `validators/auth_validator.go`

```go
type LoginUserValidator struct {
  Email    string `json:"email" binding:"required,email" message:"Email is required and must be a valid email address"`
  Password string `json:"password" binding:"required" message:"Password is required"`
}
```

- `Email`: Required, must be a valid email format.
- `Password`: Required.

### Custom Error Messages

Custom error messages for validators are set in the validators structs using the message tag. If validation fails, custom messages will be returned in the response.

For example:
```go
Username string `json:"username" binding:"required,min=3,max=50" message:"Username is required with minimum of 3 characters and maximum of 50 characters"`
```

In case a custom message is not provided, default validation messages are returned by the validation middleware.

## Auto-Generated Code

Certain files, such as validators, can be dynamically registered using the `go generate` command.
Auto-generated code is stored in files like `validators/auto_generated.go`.

### Using `go generate`

To generate the auto-generated code, run:

```bash
make generate
```

This command ensures that all validator structs in the `validators` packages are automatically registered.

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

## Development Workflow

The project includes a development workflow that simplifies the process of building and running the server in development mode.

### Live Reloading with Air

The project uses **Air** for live-reloading during development. This means that the server automatically restarts whenever changes are made to the codebase.

To start the server in development mode with live-reloading:

```bash
make develop
```

## Contributing

If you would like to contribute to the project:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -am 'Add new feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Create a new Pull Request.

All contributions are welcome!

## License

This project is open-source and available under the [MIT License](LICENSE).
