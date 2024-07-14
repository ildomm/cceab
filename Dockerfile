FROM golang:1.22.5-bookworm AS builder

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    git ca-certificates && \
    update-ca-certificates

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build binaries
RUN go build -o /api_handler ./cmd/api_handler
RUN go build -o /validator ./cmd/validator

FROM golang:1.22.5-bookworm

WORKDIR /app

COPY --from=builder /api_handler .
COPY --from=builder /validator .

# Expose port 8000 for the API handler
EXPOSE 8000

# Set environment variable for the database URL
ENV DATABASE_URL=postgres://postgres:password@postgres:5433/cceab_development?sslmode=disable

# API handler
ENTRYPOINT ["./api_handler"]

# Starts the validator
CMD ["./validator"]
