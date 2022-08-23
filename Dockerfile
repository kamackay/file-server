FROM ubuntu:latest as go

ENV DEBIAN_FRONTEND noninteractive

WORKDIR /app/
WORKDIR $GOPATH/src/gitlab.com/kamackay/filer

RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y git apt-utils golang-go

ADD ./go.mod ./

RUN go mod download && \
    go mod verify

COPY ./ ./

RUN go build -tags=jsoniter -o application.file ./index.go && cp ./application.file /app/

FROM node:16-alpine as react

WORKDIR /app

COPY ./ui/package.json ./ui/yarn.lock ./

RUN yarn --frozen-lockfile

COPY ./ui/ ./

RUN yarn build

FROM ubuntu:latest as stage

ENV DEBIAN_FRONTEND noninteractive

# Install Video Conversion Libraries
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y --reinstall ffmpeg apt-utils ca-certificates vim

# To create the temp folder inside of the image
WORKDIR /temp

WORKDIR /files

COPY --from=go /app/application.file /server

COPY --from=react /app/build /ui

COPY ./mime.types /etc

# FROM scratch
# COPY --from=stage / /

CMD /server
