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
1. Run `docker-compose up --build`
2. There will be already 9 users created in the database with the following ids:
   - `11111111-1111-1111-1111-111111111111`
   - `22222222-2222-2222-2222-222222222222`
   - `33333333-3333-3333-3333-333333333333`
   - `44444444-4444-4444-4444-444444444444`
   - `55555555-5555-5555-5555-555555555555`
   - `66666666-6666-6666-6666-666666666666`
   - `77777777-7777-7777-7777-777777777777`
   - `88888888-8888-8888-8888-888888888888`
   - `99999999-9999-9999-9999-999999999999`
3. The API Handler will be available at `http://localhost:8080`
4. The Game Results validator will be running in the background.
5. The Game Results API endpoint can be reached at `http://localhost:8000/api/v1/users/11111111-1111-1111-1111-111111111111/game_results`
6. To send a post request to the API endpoint, you can use the following curl command:
   ```bash
   curl -X POST http://localhost:8000/api/v1/users/11111111-1111-1111-1111-111111111111/game_results -H 'Content-Type: application/json' -H "Source-Type: game" -d "{\"state\": \"win\", \"amount\": \"10.15\", \"transactionId\": \"12\"}" 
   ```

## Testing
### Local tests:
#### Unit tests:
Unit tests are all written with the standard go testing library.
- Type `make unit-test` to run them.

#### Tests coverage:
- Type `make coverage-report` to generate an HTML report with the tests coverage.
- Type `make coverage-total` to check the total tests coverage.

## Future improvements
- Use Store procedures to calculate the balance.
- Allow the Validator to be configurable for pause interval and total of games to validate.
- (pending)

## Missing features
- 90% tests code coverage.

## Considerations
- The mention of 'balance should be corrected..' in the requirements, obligates the system to change the balance twice, which represents a risk. 
- We could avoid updating the balance twice if changing the design slightly.
