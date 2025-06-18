# Wallet App

This repository contains the working version (local env) of Wallet-App. Test version on Google Cloud Run is avaiable at (TBD).

## Overview

The wallet-app provides RESTful APIs for managing customers wallets including:

*   Getting Wallet Balance
*   Deposit money
*   Withdrawing money
*   Transfer money to another wallet
*   Retrieve past transactions of a wallet
*   Unit Tests (./internal/service/wallet_test.go)


This project excluding user authentication & authorisation but is designed to be able to support this if needed in future.

### System Design
API Endpoints:
- https://documenter.getpostman.com/view/28966723/2sB2x9jAUe

Database Design:
- ./migrations/ddl/001_Initialisation.sql

### What is pending (in roadmap but not yet available):
- Adding middleware mechanism for Idempotency-Key (Using Redis). This will help avoid unintended impacts of retries on POST-endpoints

## Getting Started

### Project Structure

This project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

```
wallet-app/
├── cmd/                        // Contains the main applications for the project.
│   └── rest/                   // Entry point for the REST API.
│       └── main.go
├── internal/                   // Private application and library code.
│   ├── config/                 // Loads and manages application configurations.
│   │   └── config.go
│   ├── db/                     // Database initialization and connection logic.
│   │   └── database.go
│   ├── handler/                // HTTP request handlers.
│   │   ├── wallet.go
│   │   └── mocks/             
│   ├── model/                  // Business data structures and request/response models.
│   │   ├── AmountRequest.go
│   │   ├── transaction.go
│   │   ├── transaction_type.go
│   │   ├── TransferRequest.go
│   │   ├── user.go
│   │   └── wallet.go
│   ├── repo/                   // Repository layer: handles interaction with the database.
│   │   └── wallet.go
│   ├── server/                 // HTTP server setup, middleware, and route registration.
│   │   └── http.go
│   └── service/                // Service layer: contains business logic.
│       ├── wallet.go
│       ├── wallet_test.go
│       └── mocks/           
├── deployments/                // Deployment-related configurations.
│   ├── local.env
│   └── prod.env
├── migrations/                 // Database migration scripts.
│   ├── ddl/                    // Data Definition Language (schema creation).
│   │   └── 001_Initialisation.sql
│   └── dml/                    // Data Manipulation Language (sample data insertion).
│       └── 001_Sample_Data.sql
├── pkg/                        // Reusable packages, safe for external import.
│   └── restjson/               // Utility functions for standardized JSON API responses.
│       ├── response.go
│       └── response_test.go
├── .gitignore                  // Specifies files to ignore in version control.
├── .mockery.yaml               // Configuration for the mockery tool (mock generation).
├── Dockerfile                  // Instructions to build the Docker container image.
├── docker-compose.yml          // Setup Postgres and Redis containers for local development
├── go.mod                    
├── go.sum                     
└── Makefile                    // Makefile for common development tasks.
```

### Prerequisites

*   Go (version >= 1.24)
*   Docker (for local development and containerization)

### Local Development

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd wallet-app
    ```

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```
3. **Run Postgres Server (refer to docker-compose.yaml for username/password:**
    ```bash
    docker compose up -d
    ```

4. **Run the application:**
    ```bash
    go run ./cmd/rest/main.go
    ```

## Development Workflow
### Follow Go best practices

Go Code Review Comment: https://go.dev/wiki/CodeReviewComments
Common Go Mistakes: https://100go.co/


### First time setup Database
Setup database schema under ./migrations/ddl/001_Initialisation.sql 
Adding sample data for user, wallet and transactions under ./migrations/dml/001_Sample_Data.sql

### Test
Run tests locally. Currently only some tests are available (refer to internal/service/wallet_test.go)
```bash
make local-test
```

### Mockery
Generate mock file. Ensure your machine has mockery installed.

Mockery office website: https://github.com/vektra/mockery. TL'DR

Install mockery:
```bash
go install github.com/vektra/mockery/v2@v2.53.3
```

After installation, the binary file resides in the GOPATH/bin directory. You haven't add GOPATH to your $PATH, run this command
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

To generate mock file, run this command without any additional flags - thank to `.mockery.yaml`
```bash
mockery
```
