FROM golang:1.7-alpine
ENV sourcesdir /go/src/github.com/microservices-demo/user/
ENV MONGO_HOST mytestdb:27017
ENV HATEAOS user
ENV USER_DATABASE mongodb
ENV GIT_SSL_NO_VERIFY=1

## APPD
RUN mkdir -p /opt/appdynamics/src
ADD golang-sdk-x64-linux-4.5.2.0.tar /opt/appdynamics/src
ENV GOPATH $GOPATH:/opt/appdynamics
ENV LD_LIBRARY_PATH /opt/appdynamics/src/appdynamics/lib

COPY . ${sourcesdir}
RUN apk update
RUN apk add git

## APPD GCC
RUN apk add build-base libc6-compat

RUN go get -v github.com/Masterminds/glide && cd ${sourcesdir} && glide install && go install

ENTRYPOINT user
EXPOSE 8084
