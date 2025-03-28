FROM golang:1.18-alpine AS builder
WORKDIR /app

# descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# copiar codigo fuente
COPY . .

# compilar
RUN go build -o main .

# crear imagen
FROM alpine:latest
WORKDIR /app

# copiar el stage anterior
COPY --from=builder /app/main .

# exponer puerto 8080
EXPOSE 8080

# comando para ejecutar app
CMD ["./main"]
