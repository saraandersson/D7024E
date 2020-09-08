FROM golang:alpine
WORKDIR /source
COPY udpserver.go /source
RUN go build udpserver.go
CMD ./udpserver 1234