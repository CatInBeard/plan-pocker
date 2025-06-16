FROM golang:1.24.4-alpine3.22 AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY ./websocket-server/go.mod ./websocket-server/go.sum ./
COPY ./go-shared/go.mod ./go-shared/go.sum /go-shared/
RUN go mod download

COPY ./websocket-server ./
COPY go-shared /go-shared

RUN go build -o websocket-server .

FROM scratch

COPY --from=builder /app/websocket-server /websocket-server

EXPOSE 8081

CMD ["/websocket-server"]
