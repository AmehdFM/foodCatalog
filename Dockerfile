# Etapa 1: Build
FROM golang:1.21-alpine AS builder

# Instalar dependencias del sistema
RUN apk add --no-cache git

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar todo el código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seed cmd/seed/main.go

# Etapa 2: Runtime
FROM alpine:latest

# Instalar certificados SSL para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar los binarios desde la etapa de build
COPY --from=builder /app/server .
COPY --from=builder /app/seed .

# Copiar el archivo data.json (opcional, para seed inicial)
COPY --from=builder /app/data.json .

# Exponer el puerto
EXPOSE 8080

# Comando por defecto
CMD ["./server"]