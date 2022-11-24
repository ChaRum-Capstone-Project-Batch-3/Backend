# Charum Backend

## Project Information
Charum is application for forum group discussion for many topics. This application will be used by end user and admin. The end users are the main users who use the application to perform a discussion and admin are manager all data like discussions, users, and other related data. End user can discuss certain topics while admin can manage the discussions.

## Technology Stack
1. [Echo](https://echo.labstack.com/) - Web framework
2. [Mongo](https://www.mongodb.com/) - Database
3. [Testify](github.com/stretchr/testify) - Testing

## How to use

### Prerequisites
1. [Air](https://github.com/cosmtrek/air)

### Setup
1. Copy environment file and fill it with your own values
    `cp .env.example .env`
2. Install dependencies
    `go mod download`

### Run on Local Development
1. Run the server
    `air`

### Run Unit Testing on Use Case
1. Run the test with cover
    `go test -coverprofile=cover ./business/...`
2. Show the coverage
    `go tool cover -html=cover`