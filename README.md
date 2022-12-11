# Charum Backend

## Project Information
Charum is application for forum group discussion for many topics. This application will be used by end user and admin. The end users are the main users who use the application to perform a discussion and admin are manager all data like discussions, users, and other related data. End user can discuss certain topics while admin can manage the discussions.

## Technology Stack
1. [Echo](https://echo.labstack.com/) - Web framework
2. [Mongo](https://www.mongodb.com/) - Database
3. [Testify](https://github.com/stretchr/testify) - Testing
4. [Cloudinary](https://cloudinary.com/) - Image storage
5. [JWT](https://jwt.io/) - Authentication Strategy
6. [AWS EC2](https://aws.amazon.com/ec2/) - Cloud Server
7. Github Action - CI/CD

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