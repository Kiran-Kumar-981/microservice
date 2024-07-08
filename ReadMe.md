Getting Started

Prerequisites

-> Go 1.22.4
-> PostgreSQL latest

# Setup

1. Extract the code 

2. Install dependencies:

# Running the Server

To run the gRPC server and expose HTTP endpoints:

The server will start on port `8081` by default.

# Running Tests

To run tests for the HTTP handlers:

This will run all tests, including those for the HTTP handlers and any other packages in the project.

# Endpoints

# Get User

-> Endpoint: `/user/{id}`
-> Method: GET
-> Description: Retrieves user data by ID from the gRPC server.
-> Example: `/user/1`

# Get Users

-> Endpoint: `/users?ids={id1}&ids={id2}&...`
-> Method: GET
-> Description: Retrieves multiple users by IDs from the gRPC server.
-> Example: `/users?ids=1&ids=2`

# Search Users

-> Endpoint: `/search?address={address}&phone={phone}&married={true/false}`
-> Method: GET
-> Description: Searches users based on address, phone, and marital status from the gRPC server.
-> Example: `/search?address=123%20Main%20St&phone=1234567890&married=true`

# Dependencies

-> `github.com/google/go-cmp/cmp`: Used for deep comparison of expected and actual values in tests.
-> `google.golang.org/grpc`: Provides gRPC support for communication between the HTTP handlers and the gRPC server.
-> `google.golang.org/protobuf/testing/protocmp`: Helper functions for comparing protobuf messages in tests.


This template provides a basic structure

If every Thing goes well we will get the output

Provided unit test if the code changes that will be reflected in the test cases and will figure out the mistakes have done