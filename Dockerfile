FROM ubuntu:latest as go

ENV DEBIAN_FRONTEND noninteractive

WORKDIR /app/
WORKDIR $GOPATH/src/gitlab.com/kamackay/filer

RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y software-properties-common && \
    add-apt-repository -y ppa:longsleep/golang-backports && \
    apt-get install -y git apt-utils golang-go ca-certificates

ADD ./go.mod ./

RUN go mod download && \
    go mod verify

COPY ./ ./

RUN go build -tags=jsoniter -o application.file ./index.go && cp ./application.file /app/

FROM node:alpine as react

WORKDIR /app

COPY ./ui/package.json ./ui/yarn.lock ./

RUN yarn

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

COPY --from=react /app/build /ui/
#COPY --from=react /app/build /ui/public

COPY ./mime.types /etc

FROM scratch
COPY --from=stage / /
RUN ls -alh /ui/public

CMD /server
