FROM golang:1.25.5-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o room-service cmd/main.go

FROM alpine:latest
WORKDIR /app

RUN sed -i 's/https/http/' /etc/apk/repositories
RUN apk add --no-cache ca-certificates tzdata && \
    update-ca-certificates

COPY --from=builder /app/room-service .

EXPOSE 8082

CMD ["./room-service"]