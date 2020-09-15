FROM golang:alpine
COPY main/main /home/go/src/main/main
WORKDIR /home/go/src/main/main
RUN go build main.go
CMD ./main