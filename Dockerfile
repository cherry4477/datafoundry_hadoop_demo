FROM registry.dataos.io/datafoundry/golang:1.6.2

RUN mkdir /hadoop_demo
WORKDIR /hadoop_demo
COPY main.go /hadoop_demo/
COPY vendor/github.com /go/src/github.com

RUN go build
RUN "36.110.131.65 hadoop-1.jcloud.local" >> /etc/hosts
RUN "36.110.132.55 hadoop-2.jcloud.local" >> /etc/hosts

CMD ["tail","-f","/dev/null"]
