FROM golang:alpine
WORKDIR /go/src/cli
COPY cli /go/src/cli
RUN go build cli.go
CMD ./cli