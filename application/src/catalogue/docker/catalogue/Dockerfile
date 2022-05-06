FROM golang:1.7

## APPD
# ENV GOPATH /go
# ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
# RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" /"opt/appdynamics/src" && chmod -R 777 "$GOPATH"

RUN mkdir -p /opt/appdynamics/src
ADD golang-sdk-x64-linux-4.5.2.0.tar /opt/appdynamics/src
ENV GOPATH $GOPATH:/opt/appdynamics
ENV LD_LIBRARY_PATH /opt/appdynamics/src/appdynamics/lib
# RUN ldd --version # not in alpine?
# RUN uname -a

# ENV	SERVICE_USER=myuser \
# 	SERVICE_UID=10001 \
# 	SERVICE_GROUP=mygroup \
# 	SERVICE_GID=10001 \
# 	GIT_SSL_NO_VERIFY=1
ENV GIT_SSL_NO_VERIFY=1
#
# RUN	addgroup --group ${SERVICE_GID} ${SERVICE_GROUP}
# RUN	adduser --group "${SERVICE_NAME} user" -D -H -G ${SERVICE_GROUP} -s /sbin/nologin -u ${SERVICE_UID} ${SERVICE_USER}

RUN mkdir /app
COPY . /go/src/github.com/microservices-demo/catalogue/
COPY images/ /go/images/

RUN go get -u github.com/FiloSottile/gvt
# RUN go install github.com/FiloSottile/gvt@latest
RUN cd /go/src/github.com/microservices-demo/catalogue \
## Restore Go dependencies
&& gvt restore \
## Update Gorilla/Mux to use Use()
&& gvt update github.com/gorilla/mux

# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main github.com/microservices-demo/catalogue/cmd/cataloguesvc
RUN go build -a -o /app/main github.com/microservices-demo/catalogue/cmd/cataloguesvc

## Single Container for AppD
RUN	chmod +x /app/main && \
	# chown -R ${SERVICE_USER}:${SERVICE_GROUP} /app /images && \
	setcap 'cap_net_bind_service=+ep' /app/main

# USER ${SERVICE_USER}cd /

CMD ["/app/main", "-port=80"]

EXPOSE 80

# CMD ["bash"]
