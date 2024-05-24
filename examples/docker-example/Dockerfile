FROM golang:1.19.12 as builder

WORKDIR /go/src

COPY . ./docker-exmple

RUN cd ./docker-example \
    && go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env \
    && go mod tidy \
    && go build -o docker-example main.go

ENTRYPOINT ./docker-example