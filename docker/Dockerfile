#FROM ubuntu:18.04

#RUN apt-get update
#RUN apt-get -y upgrade
#RUN apt-get -y install golang
# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab
#WORKDIR /app

#COPY . /home/go/src/D7024E

#WORKDIR /home/go/src/D7024E

#COPY . .
#COPY . /go/src/project/

#RUN echo "PWD is: $PWD"

FROM golang:alpine
RUN mkdir /app 
ADD . /app/
WORKDIR /app
RUN ls
RUN go build -o udpclient.go

#ADD ./D7024E /image

#CMD go run udpclient.go

#RUN go build

#RUN go build -o /D7024E

#RUN go run main.go config.go server.go

#CMD ["./main"]

#CMD ["./D7024E/udpclient"]

