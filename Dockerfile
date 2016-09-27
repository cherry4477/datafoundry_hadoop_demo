FROM registry.dataos.io/datafoundry/golang:1.6.2

RUN mkdir /hadoop_demo
WORKDIR /hadoop_demo
COPY main.go /hadoop_demo/
RUN go build

CMD ["tail","-f","/dev/null"]
