FROM golang:1.7
# Originally 1.6

ENV GIT_SSL_NO_VERIFY=1

## APPD
RUN mkdir -p /opt/appdynamics/src
ADD golang-sdk-x64-linux-4.5.2.0.tar /opt/appdynamics/src
ENV GOPATH $GOPATH:/opt/appdynamics
ENV LD_LIBRARY_PATH /opt/appdynamics/src/appdynamics/lib

RUN mkdir /app
COPY . /go/src/github.com/microservices-demo/payment/

# RUN go get -u github.com/FiloSottile/gvt
# RUN cd /go/src/github.com/microservices-demo/payment/ && gvt restore

RUN go get -u github.com/FiloSottile/gvt
# RUN go install github.com/FiloSottile/gvt@latest
RUN cd /go/src/github.com/microservices-demo/payment \
## Restore Go dependencies
&& gvt restore \
## Update Gorilla/Mux to use Use()
&& gvt update github.com/gorilla/mux

# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main github.com/microservices-demo/payment/cmd/paymentsvc
RUN go build -a -o /app/main github.com/microservices-demo/payment/cmd/paymentsvc

CMD ["/app/main", "-port=80"]

#EXPOSE 80

## Single Container for AppD
RUN	chmod +x /app/main && \
	# chown -R ${SERVICE_USER}:${SERVICE_GROUP} /app /images && \
	setcap 'cap_net_bind_service=+ep' /app/main

# USER ${SERVICE_USER}cd /

CMD ["/app/main", "-port=80"]

EXPOSE 80

# CMD ["bash"]
