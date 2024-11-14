FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /calc-api ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /calc-api .
COPY .env .

EXPOSE 8080

CMD ["./calc-api"]
