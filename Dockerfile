FROM golang:1.24.5-alpine AS base

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o kryptonim-app

EXPOSE 8080

CMD ["/build/kryptonim-app"]