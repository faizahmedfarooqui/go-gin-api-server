# Documentation

This document provides a detailed explanation of the project's structure, covering key components such as validators, database migrations, services, controllers, cron jobs (schedules) and the overall development workflow.

This project is designed to offer a solid foundation for developing an API server using Go and the Gin framework.

## Table of Contents

- [Project Structure](#project-structure)
  - [Explanation of Directories](#explanation-of-directories)
- [Validators](#validators)
  - [RegisterUserValidator](#registeruservalidator)
  - [LoginUserValidator](#loginuservalidator)
  - [Custom Error Messages](#custom-error-messages)
- [Database Migrations](#database-migrations)
  - [Migration Structure](#migration-structure)
  - [Running Migrations](#running-migrations)
  - [Rolling Back Migrations](#rolling-back-migrations)
- [Scheduled Tasks (Cron Jobs)](#scheduled-tasks-cron-jobs)
- [Auto-Generated Code](#auto-generated-code)
  - [Using `go generate`](#using-go-generate)
- [Development Workflow](#development-workflow)
  - [Live Reloading with Air](#live-reloading-with-air)

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

## Database Migrations

This project uses SQL-based migrations to manage schema changes. Migration files are located in the `database/migrations/` directory.

### Migration Structure

Migration files are organized with the following structure:
- **Up Migrations**: Files ending with `.up.sql` are applied when migrating **up** (applying changes).
- **Down Migrations**: Files ending with `.down.sql` are used when rolling back changes (migrating **down**).

For example:
```bash
database/migrations/
  ├── 000001_create_items_table.up.sql
  ├── 000001_create_items_table.down.sql
  ├── 000002_create_users_table.up.sql
  └── 000002_create_users_table.down.sql
```

### Running Migrations

To apply all pending migrations (up), run the following command:

```bash
make migrate-up
```

This command will apply all `*.up.sql` migration files that have not yet been applied to the database.

### Rolling Back Migrations

To rollback the last applied migration (down), use:

```bash
make migrate-down
```

This command will execute the last migration's corresponding `*.down.sql` file, reverting the changes made by the last migration.

## Scheduled Tasks (Cron Jobs)

The project uses cron jobs to handle recurring tasks. These tasks are decoupled from the main server and can be run independently.

**Structure**
- **Cron Job Registration**: Cron jobs are registered in `cmd/schedules/main.go`.
- **Task Logic**: Task-specific logic is located in `schedules/tasks.go`.

### Running Cron Jobs

To start the cron jobs, use:
```bash
make run-cron
```

### Adding a New Cron Job

You can add new cron jobs by:
1. **Defining a new task** in `schedules/tasks.go`.
2. **Registering the task** in `cmd/schedules/main.go` with the desired cron schedule.

**Example Cron Expression**:
- `0 */1 * * * *`: Runs every minute.
- `0 0 0 * * *`: Runs daily at midnight.

## Auto-Generated Code

Certain files, such as validators, can be dynamically registered using the `go generate` command.
Auto-generated code is stored in files like `validators/auto_generated.go`.

### Using `go generate`

To generate the auto-generated code, run:

```bash
make generate
```

This command ensures that all validator structs in the `validators` packages are automatically registered.

## Development Workflow

The project includes a development workflow that simplifies the process of building and running the server in development mode.

### Live Reloading with Air

The project uses **Air** for live-reloading during development. This means that the server automatically restarts whenever changes are made to the codebase.

To start the server in development mode with live-reloading:

```bash
make develop
```

**Air Configuration**: The live-reloading behaviour is configured in the `air.toml` file, where you can specify which directories to watch and which files types to trigger a rebuild on changes.

For example, in `air.toml`:

```toml
[watch]
  includes = ["./controllers", "./routes", "./services", "./validators", "./config"]
  include_ext = ["go", "html", "tmpl", "tpl"]
  exclude_dir = ["vendor", "tmp", "database/migrations"]
```

This configuration ensures that changes to `.go` files and template files will trigger a rebuild.

## Summary

This documentation provides an overview of key aspects of the project, including the validation system, migration system, and development workflow using live reloading. Use this as a reference for understanding and working with the codebase.

For any further questions or contributions, please check the project's [README.md](README.md) or submit a pull request!