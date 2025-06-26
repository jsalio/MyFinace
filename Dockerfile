# Etapa de construcción
FROM golang:1.23.9-alpine AS builder

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar los archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar el código fuente
COPY . .

# Construir la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o financial-app

# Etapa de producción
FROM alpine:latest

WORKDIR /app

# Copiar el binario desde la etapa de construcción
COPY --from=builder /app/financial-app .
COPY --from=builder /app/.env* ./

# Puerto expuesto
EXPOSE 8085

# Comando para ejecutar la aplicación
CMD ["./financial-app"]
