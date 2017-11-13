#!/bin/bash -x
cd "$(dirname "$0")"/..
mkdir -p build

cat <<EOF > build/Dockerfile-builder
FROM golang:1.9

RUN mkdir -p /gopath

ENV GOPATH /gopath

RUN go get github.com/andrewbackes/chess
RUN go get github.com/gorilla/handlers
RUN go get github.com/gorilla/mux
RUN go get github.com/sirupsen/logrus
RUN go get github.com/x-cray/logrus-prefixed-formatter
RUN go get gopkg.in/mgo.v2/bson
EOF

# debian
docker build -t tourney-builder -f build/Dockerfile-builder .
docker run -v "$(pwd)/:/gopath/src/github.com/andrewbackes/tourney" tourney-builder bash -c \
    "go build -o /gopath/src/github.com/andrewbackes/tourney/build/server /gopath/src/github.com/andrewbackes/tourney/cmd/server/main.go"

cat <<EOF > build/Dockerfile
FROM debian:stretch
COPY build/server /tourney/server
COPY cmd/server/books/2700draw.bin /tourney/2700draw.bin
WORKDIR /tourney/
CMD ["/tourney/server"]
EOF

TAG=$(git rev-parse --short HEAD)
docker build -t andrewbackes/tourney:${TAG} -f build/Dockerfile .

rm build/Dockerfile-builder
rm build/Dockerfile