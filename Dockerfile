FROM golang:onbuild

RUN go get github.com/Azure/azure-sdk-for-go/storage && go get github.com/gorilla/mux && go get github.com/kelseyhightower/envconfig

RUN mkdir /app

WORKDIR /app

ADD . /app

EXPOSE 8000

RUN go build ./server.go

CMD ["./server"]
