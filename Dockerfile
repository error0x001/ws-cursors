FROM golang:1.17.0-alpine3.14 AS builder

WORKDIR /opt/ws-cursors

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY ./go.mod ./go.sum ./
COPY ./internal/  ./internal/
COPY ./cmd/  ./cmd/

# build app binary
RUN go mod tidy && \
    go mod vendor
RUN go build -v \
        -mod=vendor \
        -tags netgo \
        -ldflags '-s' \
        -o /go/bin/ws-cursors \
        ./cmd/ws-cursors/main.go

# run
FROM phusion/baseimage:focal-1.0.0-amd64
WORKDIR /go/bin
COPY --from=builder /go/bin .
COPY ./templates/ ./templates/
ENTRYPOINT ["./ws-cursors"]
