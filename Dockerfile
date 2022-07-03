FROM        golang:1.18-alpine3.16 as base

WORKDIR     /go/lbc-api
RUN         apk -u add build-base bash

FROM        base as dev
RUN         wget https://github.com/cortesi/modd/releases/download/v0.8/modd-0.8-linux64.tgz
RUN         tar -xzf modd-0.8-linux64.tgz
RUN         mv ./modd-0.8-linux64/modd /usr/local/bin
ENTRYPOINT  [ "modd" ]
CMD         [ "-f", "configuration/modd/modd.conf" ]

#ADD         .   /go/lbc-api