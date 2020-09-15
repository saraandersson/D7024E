FROM golang:alpine
COPY main.go /home/go/src/main/
WORKDIR /home/go/src/main/
RUN go build main.go
CMD ./main