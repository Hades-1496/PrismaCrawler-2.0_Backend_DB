# Etapa 1: Construcción (Builder)
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Compilamos el binario para arquitectura Linux pura sin dependencias externas enlazadas de C
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o prismacrawler cmd/api/main.go

# Etapa 2: Producción (Contenedor ultraligero solo con el binario)
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/prismacrawler .

# Exponemos el puerto estándar para despliegues
EXPOSE 8000

CMD ["./prismacrawler"]