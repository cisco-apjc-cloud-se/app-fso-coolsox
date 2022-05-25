FROM centos:centos8

# Install Packages
RUN yum install -y \
    java-1.8.0-openjdk \
    net-tools \
    zip

# Add Environment variables
ENV JAVA_HOME /usr/lib/jvm/jre-1.8.0-openjdk

# Add user
RUN useradd -ms /bin/bash appd
WORKDIR /home/appd

# Copy AppDynamics agent onto container
COPY ["files/db-agent-21.9.0.2521.zip", "/home/appd/"]

# Unzip AppDynamcis Agent
RUN /bin/mkdir -p /home/appd/db-agent && /usr/bin/unzip /home/appd/db-agent-21.9.0.2521.zip -d /home/appd/db-agent/ && /bin/chown -R appd:appd /home/appd

# Change to appd user
USER appd

# Start the Agent
CMD ["/home/appd/db-agent/start-dbagent"]
