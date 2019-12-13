FROM golang:alpine as builder

WORKDIR /app/
WORKDIR $GOPATH/src/github.com/kamackay/file-server

RUN apk upgrade --update --no-cache && \
        apk add --no-cache \
            git \
            gcc \
            dep \
            curl \
            linux-headers \
            brotli

ADD ./Gopkg.toml Gopkg.lock ./dummy/dummy.go ./

RUN curl https://raw.githubusercontent.com/google/brotli/master/go/cbrotli/BUILD > $GOPATH/BUILD && \
    dep ensure && \
    rm dummy.go

ADD ./ ./

RUN go build -o server ./*.go && cp ./server /app/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server /app/

CMD ["./server"]
