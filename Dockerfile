FROM golang:alpine
WORKDIR /go/src/main
COPY d7024e /go/src/d7024e
COPY main /go/src/main
RUN go build main.go
CMD ./main 
#-port=8080 -bootstrap_ip=127.0.0.1 -bootstrap_port=8000