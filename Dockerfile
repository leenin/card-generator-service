FROM golang:1.8

RUN mkdir -p /go/src/app

WORKDIR /go/src/app
COPY . /go/src/app

RUN go install -v

CMD ["app"]