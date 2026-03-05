ARG REGISTRY=ghcr.io
ARG IMAGE_NAMESPACE

FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/web/main.go

FROM ${REGISTRY}/${IMAGE_NAMESPACE}/silsilah-backend-base:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

RUN mkdir -p storage/public

EXPOSE 8080

CMD ["./main"]
