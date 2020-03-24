FROM golang:alpine as builder

WORKDIR /app/
WORKDIR $GOPATH/src/github.com/kamackay/filer

RUN apk upgrade --update --no-cache && \
        apk add --no-cache \
            git \
            gcc \
            curl \
            linux-headers

ADD ./go.mod ./

RUN go mod download && \
    go mod verify

ADD ./ ./

RUN go build -o application.file ./*.go && cp ./application.file /app/

FROM alpine:latest

COPY --from=builder /app/application.file /server

CMD ["/server"]
