FROM golang:1.11


COPY . /go/src/api-aggregator
WORKDIR /go/src/api-aggregator

ENV GO111MODULE=on

RUN go build

EXPOSE 8080

CMD ./api-aggregator