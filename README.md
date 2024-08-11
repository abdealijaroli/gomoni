# Gomoni

Gomoni is a barebones banking API built from scratch in Go, emphasizing minimal dependencies and raw implementation.

## Features

- User authentication with JWT
- Account Management
- Money transfer between accounts
- PostgreSQL database integration
- Comprehensive unit tests

## Prerequisites

- Go 1.16 or higher
- PostgreSQL database (Can also just use a Docker Image for Postgres)

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/abdealijaroli/gomoni.git
   cd gomoni
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Set up your PostgreSQL database and update the `.env` file with your database URL:
   ```
   DB_URL=postgres://username:password@localhost:5432/dbname?sslmode=disable
   ```

4. Set your JWT secret in the `.env` file:
   ```
   JWT_SECRET=your_jwt_secret_here
   ```

## Usage

1. Run the server:
   ```
   go run *.go
   ```

2. The API will be available at `http://localhost:8008`

## API Endpoints

- `POST /login`: User login
- `GET /account`: Get all accounts (requires authentication)
- `POST /account`: Create a new account (requires authentication)
- `GET /account/{id}`: Get account by ID (requires authentication)
- `DELETE /account/{id}`: Delete an account (requires authentication)
- `POST /transfer`: Transfer money between accounts (requires authentication)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.
