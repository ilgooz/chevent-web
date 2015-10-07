FROM golang

WORKDIR /go/src/github.com/codeui/chevent-web/api

ADD . .

RUN go get ./... && go install

EXPOSE 80
