FROM golang:alpine
WORKDIR /server
COPY udpserver.go /server
RUN go build udpserver.go
CMD ./udpserver