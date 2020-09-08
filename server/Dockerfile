FROM golang:alpine
EXPOSE 8080
WORKDIR /source
COPY udpserver.go /source
RUN go build udpserver.go
CMD ./udpserver 8080