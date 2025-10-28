# Practical 5 - Integration Testing with TestContainers

This directory demonstrates integration testing using TestContainers with PostgreSQL, featuring a complete CRUD API for user management with comprehensive test coverage.

## Requirements

- Go 1.19 or higher
- Docker installed and running
- Internet connection (first run will download Docker images)

## Project Structure

```
practical5/
├── models/
│   └── user.go          # User data model
├── repository/
│   ├── user_repository.go      # Database operations
│   └── user_repository_test.go # Integration tests
├── migrations/
│   └── init.sql         # Database schema and seed data
└── README.md
```

## Running Tests

Basic test run:
```bash
cd Practicals/Practical_5
go test ./... -v
```

Run with coverage:
```bash
go test ./... -cover -v
```

Current test coverage: 84.5%

## What to Expect

1. First Run:
   - Docker will pull `postgres:15-alpine` image (~80MB)
   - Tests will start a PostgreSQL container
   - Initial run may take 1-2 minutes

2. Test Output:
   - You'll see container lifecycle logs
   - All CRUD operations are tested
   - 6 test functions with multiple sub-tests
   - Tests run in isolation using fresh containers

3. Test Cases:
   - Create: New user creation and duplicate email handling
   - Read: Fetch by ID and email
   - Update: Modify existing users and handle non-existent users
   - Delete: Remove users and handle non-existent users
   - List: View all users, verify counts and data

## Troubleshooting

1. "Cannot connect to Docker daemon"
   - Ensure Docker Desktop is running
   - Check `docker ps` works in terminal

2. "Container startup timeout"
   - Increase Docker resources in Docker Desktop
   - Check network connectivity
   - Try running `docker pull postgres:15-alpine` separately

## Implementation Details

The project uses:
- TestContainers for managing PostgreSQL containers
- SQL transactions for test isolation
- Automatic container cleanup
- Initialization SQL for database setup

