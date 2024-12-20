[![Workflow for Go Standard Action](https://github.com/AllanM007/simpler-test/actions/workflows/test-build-deploy.yml/badge.svg)](https://github.com/AllanM007/simpler-test/actions/workflows/test-build-deploy.yml) [![codecov](https://codecov.io/github/AllanM007/simpler-test/graph/badge.svg?token=W5ZXQ6HFO0)](https://codecov.io/github/AllanM007/simpler-test) [![Go Report Card](https://goreportcard.com/badge/github.com/AllanM007/simpler-test)](https://goreportcard.com/report/github.com/AllanM007/simpler-test)

### Backend Engineer Take Home Test @Simpler

- This is a simple CRUD REST API written in Go for a product resource microservice built using golang(gin,gorm) and using postgres as a database.

### Prerequisites

- Go 1.21+
- Docker
- Docker Compose

### Installation

1. Clone the repository

```bash
git clone https://github.com/AllanM007/simpler-test
```

2. Navigate to the directory

```bash
cd simpler-test
```

3. Build and run the Docker containers

```bash
docker compose up -d
```

### API Documentation

- This API is documented using Swagger and can be accessed at:

```
http://localhost:8080/api/swagger/index.html
```

### Pagination

- This API utilizes <strong>Offset</strong> api pagination in the products endpoint by passing <strong>page=?&limit=?</strong> parameters to the `products` endpoint.

### Tests
- The project has unit and integration tests utilizing testcontainers which can be run using:
```
go test -v ./...
```

### Endpoints

- `POST /api/v1/products`: Create a new product.
- `GET /api/v1/products`: Get all products.
- `GET /api/v1/products/:id`: Get a single product.
- `PUT /api/v1/products/:id`: Update a product.
- `DELETE /api/v1/products/:id`: Delete a product.
- `PUT /api/v1/products/:id/sale`: Product Sale.

### CI/CD

- This API has a ci/cd pipeline implemented using <strong>GitHub Actions</strong> that builds, tests and deploys the service to a container registry:  [Docker Hub](https://hub.docker.com/r/mwarangu/simpler-test) .
