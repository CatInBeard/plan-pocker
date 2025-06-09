FROM golang:1.24.4-alpine3.22

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/mitranim/gow@latest

COPY . .

EXPOSE 8081

CMD ["gow", "run", "."]
