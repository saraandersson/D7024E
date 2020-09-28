FROM golang
WORKDIR /go/src/cli
RUN go get github.com/golang/protobuf/proto
RUN go get -u github.com/golang/protobuf/protoc-gen-go
COPY cli /go/src/cli
COPY d7024e /go/src/d7024e
COPY protobuf /go/src/protobuf
RUN go build cli.go
CMD ./cli