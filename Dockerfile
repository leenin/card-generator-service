FROM golang:1.8

WORKDIR /go/src/app
COPY . /go/src/app

RUN go install -v

CMD ["app"]