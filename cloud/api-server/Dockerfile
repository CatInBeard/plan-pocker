FROM golang:1.24.4-alpine3.22 AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY api-server/go.mod api-server/go.sum ./
COPY go-shared/go.mod go-shared/go.sum /go-shared/
RUN go mod download

COPY api-server ./
COPY go-shared /go-shared

RUN go build -o api-server ./

FROM scratch

COPY --from=builder /app/api-server /api-server

EXPOSE 8083

CMD ["/api-server"]
