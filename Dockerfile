FROM registry.dataos.io/datafoundry/golang:1.6.2

RUN mkdir /hadoop_demo
WORKDIR /hadoop_demo
COPY main.go /hadoop_demo/
COPY vendor/github.com /go/src/github.com

RUN go build
COPY krb5.conf /etc/
RUN apt-get update   && \
    apt-get install -y --no-install-recommends  krb5-user krb5-config

COPY start.sh /hadoop_demo/

CMD ./start.sh
