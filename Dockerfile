FROM registry.dataos.io/datafoundry/golang:1.6.2

RUN mkdir /hadoop_demo
WORKDIR /hadoop_demo
COPY main.go /hadoop_demo/
COPY vendor/github.com /go/src/github.com

RUN go build

CMD ["tail","-f","/dev/null"]
