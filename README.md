# Ledger Service

A Go-based financial ledger service for tracking customer deposits, purchases, and restaurant commissions.

## Architecture

Clean architecture implementation with three main layers:

- **Core**: Business logic and domain types
- **Infrastructure**: Database, queue, and web implementations  
- **Handlers**: HTTP API endpoints

## Features

- Customer deposit tracking
- Purchase transaction processing with automatic 5% commission calculation
- Balance management for customers and restaurants
- Asynchronous transaction processing with in-memory queue
- MongoDB persistence layer
- RESTful API with OpenAPI documentation
- Structured JSON logging
- Context propagation with timeouts

## Transaction Types

- `DEPOSIT`: Customer adds money to their balance
- `PURCHASE`: Customer buys from restaurant (triggers commission)
- `COMMISSION`: Automatic 5% fee deducted from restaurant balance

## API Endpoints

- `POST /api/customers/{customerId}/transactions/deposits` - Create deposit
- `POST /api/customers/{customerId}/transactions/purchase` - Create purchase
- `GET /api/balances/{userId}` - Get user balance
- `GET /api/customers/{customerId}/transactions` - Get customer transactions
- `GET /api/restaurants/{restaurantId}/transactions` - Get restaurant transactions

## Running

```bash
# Start MongoDB with Docker Compose
docker-compose up -d

# Run the service
go run cmd/ledger/main.go
```

Server starts on port 8081. API documentation available at http://localhost:8081/docs

## Dependencies

- Go 1.24.4
- MongoDB (via docker-compose)
- Huma v2 (REST API framework)
- MongoDB Go Driver