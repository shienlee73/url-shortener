FROM golang:alpine

WORKDIR /app

ENV GIN_MODE=release

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main main.go

ENTRYPOINT [ "./main" ]

CMD [ "--address", "0.0.0.0", "--port", "8080", "--redis-addr", "redis:6379", "--redis-password", "" ]