FROM golang:alpine
EXPOSE 8001
WORKDIR /server
COPY udpserver.go /server
RUN go build udpserver.go
CMD ./udpserver