FROM golang:1.20-alpine
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api-gateway .

CMD ["./api-gateway"]
