FROM golang:1.24.5-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o kryptonim-app

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /build/kryptonim-app .

EXPOSE 8080

CMD ["./kryptonim-app"]