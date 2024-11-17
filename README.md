# Next-OMS
## Overview
The Next-OMS service is designed to handle order management and facilitate user authentication and authorization.

## Interactive API Documentation
Explore the API interactively using one of the following tools:

- Swagger UI: http://localhost:8080/docs/swagger
- Redoc: http://localhost:8080/docs/redoc
- Rapidoc: http://localhost:8080/docs/rapidoc
## How to Run the Service
### Local Development
1. Install dependencies:
```bash
go mod vendor
```
2. Start the service:
```bash
go run main.go serve
```
### Docker Environment
Start the service in a Docker container:
```bash
make dev
```