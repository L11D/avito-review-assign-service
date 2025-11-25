FROM golang:1.24.4-alpine AS builder

WORKDIR /app   

COPY go.mod go.sum ./
RUN go mod download

COPY .golangci.yml ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN golangci-lint run

RUN go build -o main ./cmd/app/main.go

FROM alpine:latest

COPY migrations/ ./migrations/
COPY --from=builder /app/main .

CMD [ "./main" ]