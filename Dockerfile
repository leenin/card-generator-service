FROM golang:1.11-alpine

RUN mkdir -p /go/src/github.com/leenin/card-generator-service

WORKDIR /go/src/github.com/leenin/card-generator-service

COPY . .

RUN go install -v

CMD ["card-generator-service"]

EXPOSE 8000