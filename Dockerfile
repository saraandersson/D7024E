FROM golang
RUN go get github.com/golang/protobuf/proto
RUN go get -u github.com/golang/protobuf/protoc-gen-go
WORKDIR /go/src/main
COPY protobuf /go/src/protobuf
RUN cd ..
RUN cd protobuf
RUN protoc -I ./ --go_out=./ ./message.proto
RUN cd ..
COPY d7024e /go/src/d7024e
COPY main /go/src/main
RUN go build main.go
CMD ./main
#CMD ./main -port=8080 -bootstrap_ip=127.0.0.1 -bootstrap_port=8000