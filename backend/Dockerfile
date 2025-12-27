FROM golang:alpine3.22 AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY handlers ./handlers
COPY setup ./setup
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main-app

FROM alpine:3.20
COPY --from=builder /app/main-app .


EXPOSE 8085

CMD ["./main-app"]