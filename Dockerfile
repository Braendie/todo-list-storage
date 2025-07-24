FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o ./bin/storage ./cmd/migrator/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/bin/storage .
COPY --from=builder /app/config/local.yaml ./config.yaml

ENV CONFIG_PATH=/root/config.yaml

CMD ["./storage"]