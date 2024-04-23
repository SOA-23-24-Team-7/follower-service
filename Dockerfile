FROM golang:alpine AS builder
WORKDIR /app
COPY . .
EXPOSE 8095
ENTRYPOINT ["go", "run", "main.go"]