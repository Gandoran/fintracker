FROM golang:1.26.1-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main ./cmd/server/main.go 

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/web/views ./web/views
COPY --from=builder /app/config.yaml .

EXPOSE 8080
CMD ["./main"]