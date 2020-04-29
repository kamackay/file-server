FROM registry.access.redhat.com/ubi8:latest as builder

RUN yum update -y

WORKDIR /app/
WORKDIR $GOPATH/src/gitlab.com/kamackay/filer

RUN yum install -y golang

ADD ./go.mod ./

RUN go mod download && \
    go mod verify

ADD ./ ./

RUN go build -o application.file ./*.go && cp ./application.file /app/

FROM registry.access.redhat.com/ubi8-minimal:latest

WORKDIR /files

COPY --from=builder /app/application.file /server

CMD ["/server"]
