# Fulfilment Order History

This repository contains the source code for the Order History Service, a Go-based microservice responsible for managing and providing access to customer order history data. It's designed to be deployed on Google Cloud Run.

## Overview

The Order History Service provides RESTful APIs for interacting with order history data:

*   Retrieving individual order details.
*   Listing orders for a specific customer.
*   Filtering and sorting orders based on various criteria (e.g., date, status).

This service is built using Go and is optimized for deployment on Google Cloud Run, leveraging its serverless and autoscaling capabilities.

## Getting Started

### Project Structure

This project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

```
order-history/
├── cmd/
│   └── http/
│       └── main.go             // HTTP entry point
├── internal/
│   ├── handler/
│   │   ├── http/               // Request handlers
│   │   │   ├── summary.go      
│   │   │   ├── return.go
│   │   │   └── router.go 
│   ├── config/
│   │   └── config.go           // Loads and manages application configurations
│   ├── db/           
│   │   └── db.go               // Database initialization & migrations
│   ├── presenter/              // Optional: Handle the mapping business data structure & 
│   │   ├── summary/            //           business input/ouput data
│   │   │   ├── request.go
│   │   │   ├── response.go
│   │   │   └── mapper.go
│   │   ├── return/
│   │   │   ├── request.go
│   │   │   ├── response.go
│   │   │   └── mapper.go
│   ├── model/                  // Business data structure
│   │   ├── summary.go        
│   │   └── return.go     
│   ├── service/                // Business logic
│   │   ├── summary/            
│   │   │   └── get.go
│   │   └── return/             
│   │       └── get.go
│   │       └── create.go
│   ├── repo/                   // Handle the interaction with external datasource
│   │   ├── summary.go          
│   │   └── return.go           
├── deployments/                // Handle the deployment configuration
│   ├── http/
│   │   ├── preprod.env
│   │   ├── prod.env
│   │   └── Dockerfile
├── pkg/                        // Reusable packages
├── .gitignore                  // Specifies files to ignore in version control
├── go.mod                      // Go module definition
└── go.sum                      // Go module checksums

```

### Prerequisites

*   Go (version >= 1.24)
*   Google Cloud SDK (gcloud)
*   Docker (for local development and containerization)

### Local Development

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd order-history
    ```

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Run the application:**
    ```bash
    go run .cmd/http/main.go
    ```

## Development Workflow
### Follow Go best practices

Go Code Review Comment: https://go.dev/wiki/CodeReviewComments
Common Go Mistakes: https://100go.co/

### Go Lint
Running the linter is mandatory in CI and during local development. Ensure the linter and govulncheck are run before creating a pull request (PR).

Backend Chapter: [Golang Linter](https://ntuclink.atlassian.net/wiki/spaces/BC/pages/2001960989/Golang+Linter+-+Golangci-lint)

### Test
Run tests locally
```bash
make local-test
```

### Mockery
Generate mock file. Ensure your machine has mockery installed.

Mockery office website: https://github.com/vektra/mockery. TL'DR

Instal mockery:
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

## Contact

For any questions or issues, please contact: <dpd-fulfilment@ntucenterprise.sg>