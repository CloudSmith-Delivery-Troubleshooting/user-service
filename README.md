# User Service

A simple user management system built in Go.

## Features

- Create, read, update, and delete users
- In-memory database using HashiCorp's go-memdb
- RESTful API with Gorilla Mux
- Structured logging with Zap

## Getting Started

### Prerequisites

- Go 1.20 or higher

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/saurabhbothra22/user-service.git
   cd user-service
   ```

2. Install dependencies:
   ```
   GOPROXY=direct go mod tidy
   ```

3. Build the application:
   ```
   go build -v ./cmd/server
   ```

4. Run the server:
   ```
   ./server
   ```

The server will start on port 8080.

## API Endpoints

- `POST /users` - Create a new user
- `GET /users` - List all users
- `GET /users/{email}` - Get a user by email
- `PUT /users/{email}` - Update a user
- `DELETE /users/{email}` - Delete a user

## User Model

```json
{
  "email": "user@example.com",
  "name": "User Name",
  "age": 30
}
```

## Development

### Running Tests

```
go test -v ./...
```

## CI/CD

This project uses GitHub Actions for continuous integration. The workflow runs on every push to the `ci-build` branch and on pull requests to this branch.
