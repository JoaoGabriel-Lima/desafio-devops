FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app/main ./cmd/app

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main . 

EXPOSE 8080
CMD ["./main"]
