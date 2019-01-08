FROM golang:1.11-alpine

RUN mkdir -p /go/src/card-service

WORKDIR /go/src/card-service
COPY . /go/src/card-service

RUN go install -v

CMD ["card-service"]

EXPOSE 8000