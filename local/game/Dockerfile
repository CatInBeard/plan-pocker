FROM golang:1.24.4-alpine3.22

RUN apk add --no-cache git

WORKDIR /app

COPY game/go.mod game/go.sum ./
COPY go-shared/go.mod go-shared/go.sum /go-shared/
RUN go mod download

RUN go install github.com/mitranim/gow@latest

EXPOSE 8082

CMD ["gow", "run", "."]
