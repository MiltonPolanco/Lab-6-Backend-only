FROM golang:1.20-alpine

WORKDIR /app


COPY go.mod ./
RUN go mod download


COPY . .

RUN go build -o seriestracker .

EXPOSE 8080

CMD ["./seriestracker"]
