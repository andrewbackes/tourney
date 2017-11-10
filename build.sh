#!/bin/bash
cd "$(dirname "$0")"
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
COPY build/server /server
CMD ["/server"]
EOF

docker build -t andrewbackes/tourney -f build/Dockerfile .

rm build/Dockerfile-builder
rm build/Dockerfile