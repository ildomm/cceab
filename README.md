[![tests](https://github.com/ildomm/cceab/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/ildomm/cceab/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ildomm/cceab?cache=v1)](https://goreportcard.com/report/github.com/ildomm/cceab)
![coverage](https://raw.githubusercontent.com/ildomm/cceab/badges/.badges/main/coverage.svg)

# cceab
User's account balance simulator.

Main features:
- Persist user's games results
- Validate user's account balance based on the games results

## Architecture

The application consists of two parts:
### 1. API Handler
- The API Handler is responsible for handling the HTTP requests in regards to games results.

#### API endpoints
- `GET /api/v1/health` - Returns the health of the service.
- `POST /api/v1/device/{id}/game_results` - Persists the game results of the user with the given id.

### 2. Game Results validator
- Background job that validate the user's account balance based on the games results.

### Database schema

```mermaid
---
title: Tables.
---
erDiagram
   users }|..|{ game_results : "One-to-Many"
           
```

## Build process
- Type `make build` to generate the binaries in the `bin` folder.

### API Handler
- The entrypoint is in `cmd/api_handler/main.go`

### Game Results validator
- The entrypoint is in `cmd/validator/main.go`

## Development

### Environment variables
Required environment variable for both application components:
- `DATABASE_URL` - the Postgres database URL, eg: `postgres://user:password@host:port/database`

## Deployment
Steps to deploy the application using Docker:
1. (pending)

## Testing
### Local tests:
#### Unit tests:
Unit tests are all written with the standard go testing library.
- Type `make unit-test` to run them.

#### Tests coverage:
- Type `make coverage-report` to generate an HTML report with the tests coverage.
- Type `make coverage-total` to check the total tests coverage.

## Future improvements
- (pending)

## Missing features
- 90% tests code coverage.

## Considerations
- The mention of 'balance should be corrected..' in the requirements, obligates the system to change the balance twice, which represents a risk. 
- We could avoid updating the balance twice if changing the design slightly.
