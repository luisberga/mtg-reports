FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod go.sum ./
COPY config.yaml ./
RUN go mod download

COPY . .

RUN go build -o api ./cmd/api

EXPOSE 8080

CMD ["./api"]
