FROM golang:alpine
WORKDIR /home/go/src/main
COPY d7024e /home/go/src/d7024e
COPY main /home/go/src/main
RUN go build main.go
CMD ./main