FROM fedora:rawhide as go

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

FROM fedora:rawhide as stage

# Install Video Conversion Libraries
RUN yum update -y && \
    yum -y install https://download1.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm && \
    yum -y install https://download1.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-$(rpm -E %fedora).noarch.rpm && \
    yum install -y which ffmpeg && \
    rm -rf /usr/lib/python3.9 && \
    rm -rf /usr/lib/python3.8 && \
    rm -rf /usr/lib/bin64 && \
    rm -rf /var/cache

# To create the temp folder inside of the image
WORKDIR /temp

WORKDIR /files

COPY --from=go /app/application.file /server

COPY --from=react /app/build /ui

COPY ./mime.types /etc

CMD /server
