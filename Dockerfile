FROM fedora:latest as builder

RUN dnf update -y

WORKDIR /app/
WORKDIR $GOPATH/src/gitlab.com/kamackay/filer

RUN dnf install -y golang

ADD ./go.mod ./

RUN go mod download && \
    go mod verify

ADD ./ ./

RUN go build -o application.file ./*.go && cp ./application.file /app/

FROM fedora:latest

WORKDIR /files

COPY --from=builder /app/application.file /server

CMD ["/server"]
