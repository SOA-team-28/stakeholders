FROM golang:alpine AS builder
WORKDIR /app
COPY . .
EXPOSE 50051
ENTRYPOINT ["go", "run", "main.go"]
