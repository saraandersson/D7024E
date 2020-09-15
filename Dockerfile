FROM golang:alpine
WORKDIR /go/src/main
COPY d7024e /go/src/d7024e
COPY main /go/src/main
RUN go build main.go
CMD ./main