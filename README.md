# Vinyl Store API

A RESTful API for managing a vinyl record store inventory, built with Go using the Gin framework.

## Project Description

This project consists of the development of a RESTful API in Golang that simulates the inventory management of a vinyl store. The system provides endpoints for user authentication, album management, and system status monitoring. Access to most endpoints is restricted through a token-based authentication mechanism.



## Getting Started

### Prerequisites

- Go 1.21 or higher  
- Git  

### Installation and Execution

```bash
# Clone the repository
git clone https://github.com/<your-username>/vinyl-store-api.git
cd vinyl-store-api

# Install dependencies
go mod tidy

# Run the server
go run main.go
```

The server runs on:

http://localhost:8080

---

## Default Users

| Username | Password  |
|----------|-----------|
| admin    | admin123  |
| user1    | pass1     |
| user2    | pass2     |

---

## API Usage

The API can be accessed using tools such as curl, Postman, or any HTTP client that allows custom headers.

### Authentication Flow

1. The user sends credentials to the `/login` endpoint.
2. The server returns an access token.
3. The token must be included in the `Authorization` header for protected routes.
4. The user can revoke the token through the `/logout` endpoint.

---

## Endpoints

### Login

```bash
curl -u admin:admin123 http://localhost:8080/login
```

### Logout

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/logout
```

### Get All Albums

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/albums
```

### Get Album by ID

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/albums/2
```

### Create Album

```bash
curl -X POST http://localhost:8080/createAlbum   -H "Authorization: Bearer <TOKEN>"   -H "Content-Type: application/json"   -d '{"id": "4", "title": "Sample Album", "artist": "Sample Artist", "price": 49.99}'
```

### System Status

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/status
```

---

## Concurrency Handling

The API supports multiple users interacting with the system at the same time. Shared resources such as the album list and active tokens are protected using `sync.RWMutex`.

---

## Error Handling

- 400: Invalid request
- 401: Unauthorized access
- 404: Resource not found
- 409: Conflict due to duplicate data

---

## Future Improvements

- Persistent database integration  
- Search and filtering  
- Role-based access control  
- Token expiration  
- Update and delete endpoints  

---

## Technologies Used

- Go  
- Gin Framework  
- HTTP Basic Authentication  
- Token-based authorization  

---

## Summary

This project demonstrates the implementation of a RESTful API with authentication, concurrent access handling, and basic inventory management functionality using Go.