FROM golang:alpine
WORKDIR /source
COPY main/main.go /source
RUN go build main.go
CMD ./main