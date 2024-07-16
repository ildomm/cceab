
[![tests](https://github.com/ildomm/cceab/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/ildomm/cceab/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ildomm/cceab?cache=v1)](https://goreportcard.com/report/github.com/ildomm/cceab)
![coverage](https://raw.githubusercontent.com/ildomm/cceab/badges/.badges/main/coverage.svg)

# cceab
User Account Balance Simulator

## Main Features
- Persistence of user game results
- Validation of user account balances based on game results

## Architecture
The application is composed of two primary components:

### 1. API Handler
The API Handler manages HTTP requests related to game results.

#### API Endpoints
- `GET /api/v1/health` - Returns the health status of the service.
- `POST /api/v1/users/{id}/game_results` - Persists the game results for the specified user.

### 2. Game Results Validator
A background job that validates user account balance based on game results.

### Database Schema
```mermaid
---
title: Tables.
---
erDiagram
   users }|..|{ game_results : "One-to-Many"
           
```

## Build Process
- Execute `make build` to generate the binaries in the `bin` directory.

### Entry Points
- API Handler: `cmd/api_handler/main.go`
- Game Results Validator: `cmd/validator/main.go`

## Development

### Environment Variables
The following environment variable is required for both application components:
- `DATABASE_URL` - The Postgres database URL, e.g., `postgres://user:password@host:port/database`

## Deployment
To deploy the application using Docker Compose, follow these steps:
- Execute `docker-compose up --build`

### To perform a manual test:
1. The database will be pre-populated with users having the following IDs:
   - `11111111-1111-1111-1111-111111111111`
   - `22222222-2222-2222-2222-222222222222`
   - `33333333-3333-3333-3333-333333333333`
   - `44444444-4444-4444-4444-444444444444`
   - `55555555-5555-5555-5555-555555555555`
   - `66666666-6666-6666-6666-666666666666`
   - `77777777-7777-7777-7777-777777777777`
   - `88888888-8888-8888-8888-888888888888`
   - `99999999-9999-9999-9999-999999999999`
2. The Game Results Validator will run as a background process.
3. Send a POST request to the API endpoint, using the following command:
   ```bash
   curl -X POST http://localhost:8000/api/v1/users/11111111-1111-1111-1111-111111111111/game_results -H 'Content-Type: application/json' -H "Source-Type: game" -d "{\"state\": \"win\", \"amount\": \"10.15\", \"transactionId\": \"12\"}" 
   ```

## Testing

### Local Tests
#### Unit Tests
Unit tests are implemented using the standard Go testing library.
- Execute `make unit-test` to run the tests.

#### Test Coverage
- Execute `make coverage-report` to generate an HTML report of the test coverage.
- Execute `make coverage-total` to view the total test coverage.

## Future Improvements
- Implement stored procedures to calculate account balances.
- Allow configuration of the Validator for pause intervals and the number of games to validate.
- Additional enhancements are pending.

## Missing Features
- Achieving 90% test code coverage.

## Considerations
- The requirement to "correct the balance" necessitates updating the balance twice, introducing a risk.
- This risk can be mitigated by adjusting the design.
