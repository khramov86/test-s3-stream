FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o s3-downloader .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/s3-downloader .

EXPOSE 8080

CMD ["./s3-downloader"]