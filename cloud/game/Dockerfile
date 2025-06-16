# Первая стадия: сборка
FROM golang:1.24.4-alpine3.22 AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY game/go.mod game/go.sum ./
COPY go-shared/go.mod go-shared/go.sum /go-shared/
RUN go mod download

COPY game ./
COPY go-shared /go-shared

RUN go build -o game ./

FROM scratch

COPY --from=builder /app/game /game

EXPOSE 8082

CMD ["/game"]
