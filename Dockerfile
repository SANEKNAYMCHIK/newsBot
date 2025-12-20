FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

RUN mkdir -p ssl
COPY --from=builder /app/ssl/* ./ssl/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .env

EXPOSE 8443
EXPOSE 8080

CMD ["./main"]