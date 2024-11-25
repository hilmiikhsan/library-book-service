FROM golang:1.22.8-alpine

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod tidy

COPY . .

COPY .env .

RUN go build -o library-book-service

RUN chmod +x library-book-service

EXPOSE 9093

EXPOSE 6003

CMD ["./library-book-service"]
