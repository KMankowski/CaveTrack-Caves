## THIS IS A DRAFT. Needs updating when deploying to staging or prod.

FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/server

FROM alpine:3.22.1

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 50051
CMD ["./server"]
