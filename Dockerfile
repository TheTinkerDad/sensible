FROM golang:latest AS builder 

ENV GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/github.com/TheTinkerDad/sensible

RUN go get github.com/GeertJohan/go.rice \
	github.com/GeertJohan/go.rice/rice \
	github.com/mattn/go-sqlite3

COPY main.go .
COPY data ./data

RUN go build

RUN rice append --exec sensible 

FROM python:3

RUN pip3 install ansible

COPY --from=builder /go/src/github.com/TheTinkerDad/sensible/sensible .

EXPOSE 8080

ENTRYPOINT ["/sensible"]  

