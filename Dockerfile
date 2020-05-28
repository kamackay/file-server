FROM registry.access.redhat.com/ubi8:latest as go

RUN yum update -y

WORKDIR /app/
WORKDIR $GOPATH/src/gitlab.com/kamackay/filer

RUN yum install -y golang git

ADD ./go.mod ./

RUN go mod download && \
    go mod verify

COPY ./ ./

RUN go build -o application.file ./*.go && cp ./application.file /app/

FROM node:alpine as react

WORKDIR /app

COPY ./ui/package.json ./

RUN yarn

COPY ./ui/ ./

RUN yarn build

FROM registry.access.redhat.com/ubi8:latest

WORKDIR /temp
# To create the temp folder inside of the image

WORKDIR /files

COPY --from=go /app/application.file /server

COPY --from=react /app/build /ui

COPY ./mime.types /etc

CMD /server
