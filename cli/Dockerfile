FROM golang:alpine
WORKDIR /go/src/cli
COPY cli /go/src/cli
COPY d7024e /go/src/d7024e
RUN go build cli.go
CMD ./cli