FROM golang:1.7

## APPD
RUN mkdir -p /opt/appdynamics/src
ADD golang-sdk-x64-linux-4.5.2.0.tar /opt/appdynamics/src
ENV GOPATH $GOPATH:/opt/appdynamics
ENV LD_LIBRARY_PATH /opt/appdynamics/src/appdynamics/lib

# ENV	SERVICE_USER=myuser \
# 	SERVICE_UID=10001 \
# 	SERVICE_GROUP=mygroup \
# 	SERVICE_GID=10001 \
# 	GIT_SSL_NO_VERIFY=1

ENV GIT_SSL_NO_VERIFY=1
ENV MONGO_HOST mytestdb:27017
ENV HATEAOS user
ENV USER_DATABASE mongodb
#
# RUN	addgroup --group ${SERVICE_GID} ${SERVICE_GROUP}
# RUN	adduser --group "${SERVICE_NAME} user" -D -H -G ${SERVICE_GROUP} -s /sbin/nologin -u ${SERVICE_UID} ${SERVICE_USER}

RUN mkdir /app
COPY . /go/src/github.com/microservices-demo/user/

RUN go get -v github.com/Masterminds/glide
# RUN cd /go/src/github.com/microservices-demo/user
# NOTE: Disable version of github.com/gorilla/mux in glide.lock
RUN cd /go/src/github.com/microservices-demo/user && glide install && go install
# RUN go install

ENTRYPOINT user
EXPOSE 8084


# ## Single Container for AppD
# RUN	chmod +x /app/main && \
# 	# chown -R ${SERVICE_USER}:${SERVICE_GROUP} /app /images && \
# 	setcap 'cap_net_bind_service=+ep' /app/main
#
# # USER ${SERVICE_USER}cd /
#
# CMD ["/app/main", "-port=80"]
#
# EXPOSE 80
#
# # CMD ["bash"]
