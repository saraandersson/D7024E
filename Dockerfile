FROM golang:alpine
WORKDIR /home/go/src/main
COPY main/main.go /home/go/src/main
RUN go build main.go
CMD ./main