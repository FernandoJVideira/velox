# Velox CLI Documentation

## Table of Contents

- [Velox CLI Documentation](#velox-cli-documentation)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Installation](#installation)
    - [Prerequisites](#prerequisites)
    - [Installation Steps](#installation-steps)
  - [CLI Commands Overview](#cli-commands-overview)
  - [Command Reference](#command-reference)
    - [Project Creation](#project-creation)
      - [`new <appname>`](#new-appname)
    - [Server Management](#server-management)
      - [`up`](#up)
      - [`down`](#down)
    - [Database Operations](#database-operations)
      - [`migrate`](#migrate)
      - [`migrate down`](#migrate-down)
      - [`migrate reset`](#migrate-reset)
    - [Code Generation](#code-generation)
      - [`make migration <name> <format>`](#make-migration-name-format)
      - [`make auth`](#make-auth)
      - [`make handler <name>`](#make-handler-name)
      - [`make model <name>`](#make-model-name)
      - [`make session`](#make-session)
      - [`make mail <name>`](#make-mail-name)
    - [Utility Commands](#utility-commands)
      - [`make key`](#make-key)
      - [`help`](#help)
      - [`version`](#version)
  - [Implementation Examples](#implementation-examples)
    - [Creating a New Project](#creating-a-new-project)
    - [Building a Basic CRUD Application](#building-a-basic-crud-application)
    - [Implementing Authentication](#implementing-authentication)
    - [Working with Email Templates](#working-with-email-templates)
    - [Database Migration Workflow](#database-migration-workflow)
  - [Best Practices](#best-practices)
  - [Troubleshooting](#troubleshooting)

## Introduction

Velox is an open-source web framework built in Go, designed to provide a Laravel-like experience with the performance benefits of Go. The framework offers:

- Simple and intuitive CLI for project management
- Web page rendering with Go templates and Jet template engine
- Support for multiple database types (MySQL/MariaDB and PostgreSQL)
- Database migration system
- Session management with multiple storage options
- Cache management (Badger or Redis)
- CSRF protection
- Comprehensive email system
- Authentication system with SSO support
- Remote file storage systems (S3, SFTP, WebDAV, MinIO)
- RPC support
- Testing utilities

For detailed code examples demonstrating how to use the Velox framework in practice, please refer to the [EXAMPLES.md](EXAMPLES.md) file, which contains sample code snippets for application structure, routing, handlers, models, and more.

## Installation

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Make](https://www.gnu.org/software/make/)

### Installation Steps

1. Clone the Velox repository:

```bash
git clone https://github.com/FernandoJVideira/velox.git
```

2. Change directory to the Velox folder:

```bash
cd velox
```

3. Build the CLI:

For macOS/Linux:

```bash
make build
```

For Windows:
First, open the `Makefile` and change line 15 to:

```
@go build -o dist/velox.exe ./cmd/cli
```

Then run:

```bash
make build
```

4. Add the CLI to your PATH:

For macOS/Linux:

```bash
export PATH=$PATH:/path/to/velox/dist
```

For Windows:
Add the `dist` directory to your Path Environment Variable.

5. Verify installation:

```bash
velox version
```

## CLI Commands Overview

```
Available commands:

help                           - Shows this help message
new <appname>                  - Creates a new Velox application
down                           - Put the Server in maintenance mode
up                             - Take the Server out of maintenance mode
version                        - Shows the current version of the CLI
migrate                        - Runs all up migrations that have not been run yet
migrate down                   - Reverts the most recent migration
migrate reset                  - Runs all down migrations in reverse order and then runs all up migrations
make migration <name> <format> - Creates up and down migration files in the migrations folder; format = sql/fizz (default fizz)
make auth                      - Creates the auth tables, files and middleware
make handler <name>            - Creates a stub handler file in the handlers directory
make model <name>              - Creates a stub model file in the data directory
make session                   - Creates a table in the database to store session data
make key                       - Generates a random 32 character key
make mail <name>               - Generates 2 starter mail templates in the mail directory
```

## Command Reference

### Project Creation

#### `new <appname>`

Creates a new Velox application with the given name. This command:

1. Clones the starter application template
2. Configures environment variables
3. Updates import paths and package names
4. Sets up the Go modules

**Example:**

```bash
velox new myapp
```

This will create a new directory called `myapp` with a ready-to-use Velox application structure.

**Output Directory Structure:**

```
myapp/
├── config/
├── data/
├── handlers/
├── migrations/
├── middleware/
├── public/
├── templates/
├── views/
├── .env
├── go.mod
├── go.sum
├── Makefile
└── main.go
```

### Server Management

#### `up`

Takes the server out of maintenance mode. Uses the RPC client to communicate with the running application.

**Example:**

```bash
velox up
```

#### `down`

Puts the server in maintenance mode. Uses the RPC client to communicate with the running application.

**Example:**

```bash
velox down
```

### Database Operations

#### `migrate`

Runs all pending up migrations that haven't been applied yet.

**Example:**

```bash
velox migrate
```

#### `migrate down`

Reverts the most recently applied migration.

**Example:**

```bash
velox migrate down
```

#### `migrate reset`

Runs all down migrations in reverse order and then runs all up migrations. This effectively resets the database to a clean state.

**Example:**

```bash
velox migrate reset
```

### Code Generation

#### `make migration <name> <format>`

Creates up and down migration files in the migrations folder. The format can be either `fizz` (default) or `sql`.

**Example:**

```bash
velox make migration create_products fizz
```

Creates:

- `migrations/TIMESTAMP_create_products.fizz.up.sql`
- `migrations/TIMESTAMP_create_products.fizz.down.sql`

#### `make auth`

Generates authentication tables, files, and middleware. This command:

1. Creates users, tokens, and remember_tokens tables
2. Adds user model code
3. Adds authentication middleware
4. Creates authentication handlers
5. Adds login and password reset views
6. Sets up email templates for password resets

**Example:**

```bash
velox make auth
```

#### `make handler <name>`

Creates a stub handler file in the handlers directory.

**Example:**

```bash
velox make handler products
```

Creates: `handlers/products.go`

#### `make model <name>`

Creates a stub model file in the data directory.

**Example:**

```bash
velox make model product
```

Creates: `data/product.go`

#### `make session`

Creates a database table to store session data, enabling database-backed sessions.

**Example:**

```bash
velox make session
```

#### `make mail <name>`

Generates two starter mail templates in the mail directory.

**Example:**

```bash
velox make mail welcome
```

Creates:

- `mail/welcome.html.tmpl`
- `mail/welcome.plain.tmpl`

### Utility Commands

#### `make key`

Generates a random 32-character key, suitable for encryption purposes.

**Example:**

```bash
velox make key
```

#### `help`

Shows a help message with all available commands.

**Example:**

```bash
velox help
```

#### `version`

Shows the current version of the CLI tool.

**Example:**

```bash
velox version
```

## Implementation Examples

### Creating a New Project

This example demonstrates creating a new project and setting it up:

```bash
# Create a new application
velox new mywebapp

# Change to the application directory
cd mywebapp

# Generate a secure key and update .env file
velox make key
# Copy the generated key and update KEY= in your .env file

# Set up database connection in .env file
# DATABASE_TYPE=postgres
# DATABASE_HOST=localhost
# DATABASE_PORT=5432
# DATABASE_USER=postgres
# DATABASE_PASS=password
# DATABASE_NAME=mywebapp
# DATABASE_SSL_MODE=disable

# Run the application using the Makefile
make start
```

When creating a new application, Velox generates a Makefile with useful commands to manage your project. Common commands include:

```bash
# Start the application
make start

# Build the application
make build

# Run tests
make test

# Clean generated files
make clean
```

Use these Makefile commands instead of running Go commands directly to ensure consistent behavior across development environments.

### Building a Basic CRUD Application

Here's how to implement a basic product management system:

```bash
# Create database migration for products
velox make migration create_products

# Edit the generated up migration file to define your schema
# Example for fizz format:
# add_column("products", "id", "integer", {"primary": true})
# add_column("products", "name", "string", {"size": 255})
# add_column("products", "description", "text", {})
# add_column("products", "price", "decimal", {"precision": 10, "scale": 2})
# add_column("products", "created_at", "timestamp", {})
# add_column("products", "updated_at", "timestamp", {})

# Run the migration
velox migrate

# Generate product model
velox make model product

# Edit the generated model to match your schema

# Create a handler for product operations
velox make handler products

# Edit the handler to implement CRUD operations

# Update your routes in your application
```

### Implementing Authentication

Set up a complete authentication system:

```bash
# Generate authentication components
velox make auth

# Run migrations to create auth tables
velox migrate

# Add authentication routes to your application
# Update your templates to include login/register forms

# Protect routes using the generated middleware
# Example in your routes setup:
# app.Routes.Use(middleware.AuthenticateMiddleware)
```

### Working with Email Templates

Create and use email templates:

```bash
# Generate email templates
velox make mail password_reset

# Edit the templates in the mail directory to customize the content
# Example HTML template:
# <h1>Password Reset</h1>
# <p>Dear {{.Name}},</p>
# <p>Click the link below to reset your password:</p>
# <p><a href="{{.Link}}">Reset Password</a></p>

# Use the mail service in your application:
# data := map[string]any{
#    "Name": user.FirstName,
#    "Link": resetLink,
# }
# msg := mailer.Message{
#    To:       user.Email,
#    Subject:  "Password Reset",
#    Template: "password_reset",
#    Data:     data,
# }
# app.Mail.SendSMTPMessage(msg)
```

### Database Migration Workflow

Manage database schema changes:

```bash
# Create initial migration
velox make migration initial_schema

# Update migration files with your schema

# Apply migrations
velox migrate

# Make changes to schema
velox make migration add_category_to_products

# Edit the new migration files

# Apply the new migration
velox migrate

# If you need to roll back
velox migrate down

# Complete reset if needed during development
velox migrate reset
```

## Best Practices

1. **Environment Configuration**:

   - Keep sensitive information in the `.env` file
   - Use different environments for development and production

2. **Project Structure**:

   - Organize handlers by feature area
   - Keep models in the data directory
   - Use middleware for cross-cutting concerns

3. **Authentication**:

   - Use the built-in authentication system
   - Implement proper authorization checks
   - Use HTTPS in production

4. **Database Operations**:

   - Use migrations for all schema changes
   - Write down migrations both for up and down operations
   - Test migrations before applying to production

5. **Session Management**:

   - Choose the appropriate session store based on your needs
   - For distributed systems, use Redis or database-backed sessions

6. **Error Handling**:
   - Implement consistent error handling
   - Log errors appropriately
   - Return user-friendly error messages

## Troubleshooting

**Problem**: `velox: command not found`
**Solution**: Ensure the velox executable is in your PATH. Check installation steps.

**Problem**: Database connection errors
**Solution**: Verify database credentials in the `.env` file and ensure the database server is running.

**Problem**: Migration errors
**Solution**: Check the syntax in your migration files. For SQL migrations, ensure compatibility with your database.

**Problem**: RPC commands (up/down) not working
**Solution**: Make sure your application is running and the RPC port specified in your .env file is correct.

**Problem**: Authentication not working
**Solution**: Verify that you've run the auth migrations and the middleware is correctly applied to your routes.

---

This documentation provides a comprehensive overview of the Velox CLI tool and its capabilities. For additional support or to report issues, please visit the project repository.
