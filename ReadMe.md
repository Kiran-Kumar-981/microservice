Project Title: Totality Corp

Description

This project is a Gin server that provides a REST API for user data. It uses a PostgreSQL database to store user information. its good as the gin implements protobuf as a prerequisites

Getting Started

Prerequisites

- Go 1.22
- PostgreSQL latest
- Gin latest
- gRPC latest

Installation

1. Clone the repository: `git clone 
2. Install dependencies: go get
3. Start the server: go run main.go

Usage


API Endpoints

- GET /users: Returns a list of all users
- GET /user/:id: Returns a single user by ID
- POST /user: Creates a new user
- PUT /user/:id: Updates a single user
- DELETE /user/:id: Deletes a single user

gRPC Endpoints

- UserService/GetUsers: Returns a list of all users
- UserService/GetUser: Returns a single user by ID
- UserService/CreateUser: Creates a new user
- UserService/UpdateUser: Updates a single user
- UserService/DeleteUser: Deletes a single user


Acknowledgments
