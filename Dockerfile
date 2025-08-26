FROM docker.io/library/golang:1.25.0-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o kryptonim-app ./cmd/kryptonim

FROM docker.io/library/alpine:3.22.1

WORKDIR /app

COPY --from=builder /build/kryptonim-app .

EXPOSE 8080

ENTRYPOINT ["./kryptonim-app"]
