# Etapa 1: build
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Copiar mod y sum primero para aprovechar cache
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el código
COPY . .

# Compilar binario
RUN go build -o mi-api main.go

# Etapa 2: contenedor final
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/mi-api .

EXPOSE 8080

CMD ["./mi-api"]