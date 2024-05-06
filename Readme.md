# Go-http-gorm-example

This project is a simple example demonstrating how to use Go (Golang) with the GORM library to build a basic net/http server with CRUD (Create, Read, Update, Delete) functionality for managing resources in a database.

## Prerequisite

Before you begin, ensure you have the following installed on your machine:

- Go (Version 1.22 or higher)
- Docker

## Installation

1. Clone the repository to your local machine

```bash
git clone https://github.com/yourusername/go-http-gorm-example.git
```

2. Navigate into the project directory

```bash
cd go-http-gorm-example
```

3. Install dependency

```bash
go mod tidy
```

## Configuration

1. Set up your database configuration in configs/config.go. You can change the database connection details, application port and other settings according to your environment.

## Usage

You can run the application via docker-compose file in case not having database being setup via following command

```bash
docker-compose -f docker-compose.dev.yaml up -d
```

or running a normal main package

```bash
go run main.go
```

## API Endpoints

- GET api/v1/books query all books information in a server (pagination query)
- GET /api/v1/books/{id} query book information of the given ID
- POST /api/v1/books add a new book into a server
- PUT /api/v1/books/{id} update book information of the given ID
- DELETE /api/v1/books/{id} delete a book from a server.

You can see all endpoints or try to call APIs via swagger at `http://localhost:${Config.App.Port}/api/v1/swagger/index.html`
