FROM golang:1.26.2 AS builder 

WORKDIR /app 

COPY go.mod go.sum ./

RUN go mod download 

COPY . . 

RUN CGO_ENABLED=0 go build -o shortener ./cmd/shortener

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/shortener ./shortener

CMD ["./shortener"]