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

- The API is documented using Swagger and can be accessed at:

```
http://localhost:8080/swagger/index.html
```

### Pagination

- The API utilizes <strong>Offset</strong> api pagination in the products endpoint by passing <strong>page=&limit=</strong> parameters to the `products` endpoint.

### Tests
- The project has a couple of tests which can be run using:
```
go test -v ./...
```

### Endpoints

- `POST /api/v1/create-product`: Create a new product.
- `GET /api/v1/products`: Get all products.
- `GET /api/v1/product/:id`: Get a single product.
- `PUT /api/v1/update-product/:id`: Update a product.
- `DELETE /api/v1/delete-product/:id`: Delete a product.
- `POST /api/v1/product-sale`: Product Sale.
