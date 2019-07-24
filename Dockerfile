FROM golang:latest

WORKDIR /go/src/pickabooker

COPY . /go/src/pickabooker

RUN cd /go/src/pickabooker && go build -o server

ENTRYPOINT "./server"

EXPOSE 8080


