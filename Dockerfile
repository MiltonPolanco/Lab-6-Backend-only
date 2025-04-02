FROM golang:1.20-alpine

WORKDIR /app

# Copiamos archivos de módulos y descargamos dependencias
COPY go.mod ./
RUN go mod download

# Copiamos el resto del código
COPY . .

# Compilamos la aplicación
RUN go build -o seriestracker .

EXPOSE 8080

CMD ["./seriestracker"]
