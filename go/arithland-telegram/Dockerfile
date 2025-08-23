FROM golang:1.24.6-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./arithland-telegram ./arithland-telegram

RUN go build -o bot ./arithland-telegram/cmd/bot

CMD ["./bot"]
