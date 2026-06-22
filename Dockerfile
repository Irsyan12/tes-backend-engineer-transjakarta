FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/fleet-api ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/fleet-api .
COPY .env .

EXPOSE 3000

CMD ["./fleet-api"]
