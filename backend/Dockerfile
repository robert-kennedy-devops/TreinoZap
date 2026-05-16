FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/api     ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/migrate ./cmd/migrate
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/worker  ./cmd/worker

FROM alpine:3.19

RUN apk add --no-cache ca-certificates wget

WORKDIR /app

COPY --from=builder /app/bin/api     .
COPY --from=builder /app/bin/migrate .
COPY --from=builder /app/bin/worker  .
COPY --from=builder /app/migrations  ./migrations

RUN mkdir -p /app/whatsapp-store

EXPOSE 8080

CMD ["sh", "-c", "./migrate up && ./api"]
