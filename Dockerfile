FROM golang:1.11

RUN mkdir -p /go/src/app

WORKDIR /go/src/app

COPY . /go/src/app

RUN go get -d -v ./...

RUN go install -v ./...

RUN go build -o avito-test-task

CMD ["./avito-test-task"]
