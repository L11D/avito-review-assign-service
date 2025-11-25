FROM golang:1.24.4-alpine AS builder

WORKDIR /app   

COPY go.mod go.sum ./

COPY pkg/ ./pkg/
COPY e2e_tests/tests/ ./tests/

CMD ["go", "test", "-v", "./..."]