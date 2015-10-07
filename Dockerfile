FROM golang

RUN curl -sSL https://get.docker.com/ | sh

WORKDIR /go/src/github.com/codeui/chevent-web

ADD . .

RUN cd updater && go get ./... && go install

VOLUME ["/var/run/docker.sock"]
