FROM golang:alpine
EXPOSE 1235
WORKDIR /source
COPY udpserver.go /source
RUN go build udpserver.go
CMD ./udpserver 1235