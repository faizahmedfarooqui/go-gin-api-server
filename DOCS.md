# Documentation

This document provides an overview of the key components in the project, including validators, database migrations, and other boilerplate features. This project is designed to offer a solid foundation for developing an API server using Go and the Gin framework.

## Table of Contents

- [Validators](#validators)
  - [RegisterUserValidator](#registeruservalidator)
  - [LoginUserValidator](#loginuservalidator)
  - [Custom Error Messages](#custom-error-messages)
- [Database Migrations](#database-migrations)
  - [Migration Structure](#migration-structure)
  - [Running Migrations](#running-migrations)
  - [Rolling Back Migrations](#rolling-back-migrations)
- [Auto-Generated Code](#auto-generated-code)
  - [Using `go generate`](#using-go-generate)
- [Development Workflow](#development-workflow)
  - [Live Reloading with Air](#live-reloading-with-air)

---

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