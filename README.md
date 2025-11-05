# Go Todo REST API 
## About

It is a RESTful API with Go using **echo**, **goose** migrations, **docker-compose** and **gorm**. 

## Features

- User registration and login with JWT token  
- CRUD operations for tasks (Create, Read, Update, Delete) tied to authenticated user  
- Authentication middleware for protected routes  
- Database migrations using **Goose**  
- Dockerised setup with PostgreSQL  
- Full OpenAPI 3.0 documentation  

## Installation & Run
```bash
# Download this project
go get github.com/daniltaro/ToDo-API
```

Before running API server, you should create the **.env** with your variable values 
#### Example of .env:
```.env
SECRET=123
PORT=8080
DB_PASS=password
DB_USER=root
DB_NAME=postgres
DB_HOST=db
DB_PORT=5432
GOOSE_DRIVER=postgres
GOOSE_DBSTRING="host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASS} sslmode=disable database=${DB_NAME}"
```
#### Run:
```bash
# Build and Run
cd ToDo-API
make up ## 	docker compose up --build

# API Endpoint : http://127.0.0.1:8080

## if you want to run tests
make tests ## 	go test ./internal/middleware -v 
           ##   go test ./internal/handler -v

## if you want to remove the containers: 
make down ## docker compose down -v
```

## Documentation
The OpenAPI 3.0 specification is provided in the file:
```
docs/swagger.yaml
```

## API

#### /tasks
* `GET` : Get all tasks
* `POST` : Create a new task
#### /projects/:id
* `PUT` : Update a task
* `DELETE` : Delete a task
#### /login
* `POST` : Login to account
#### /signup
* `POST` : Create user
#### /validate 
* `GET` : Validate current JWT token

## Tech Stack
*   Language: Go (version 1.23+ recommended)
* 	Framework: Echo
*   Database: PostgreSQL 
*   ORM: GORM
*   Authentication: JWT (HMAC-SHA256)
*   Migrations: Goose
*   Documentation: Swagger (OpenAPI 3.0)
*   Containerisation: Docker & Docker Compose	
*   Testing: Go’s testing package + github.com/stretchr/testify