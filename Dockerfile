FROM golang:1.9

RUN mkdir -p /go/src/github.com/jwma/jump-jump
ADD $PWD/app /go/src/github.com/jwma/jump-jump/app
WORKDIR /go/src/github.com/jwma/jump-jump/app

RUN go get -v -d ./
RUN go build -v -o app