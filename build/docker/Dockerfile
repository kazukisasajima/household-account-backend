FROM golang:1.23.1-alpine3.19

RUN apk --no-cache add curl

WORKDIR /go/src/web

COPY ./ .

RUN go mod download

ENV GO111MODULE=on

RUN go build -C ./cmd/server/

# ENTRYPOINT ./cmd/server/server
ENTRYPOINT ["./cmd/server/server"]