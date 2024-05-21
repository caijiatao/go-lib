FROM golang:1.19.12 as builder

WORKDIR /go/src

RUN cat /dev/null  > /etc/apt/sources.list.d/debian.sources && \
    echo "deb https://mirrors.tuna.tsinghua.edu.cn/debian/ bookworm main contrib non-free non-free-firmware" > /etc/apt/sources.list && \
    echo "deb https://mirrors.tuna.tsinghua.edu.cn/debian/ bookworm-updates main contrib non-free non-free-firmware" >> /etc/apt/sources.list && \
    echo "deb https://mirrors.tuna.tsinghua.edu.cn/debian/ bookworm-backports main contrib non-free non-free-firmware" >> /etc/apt/sources.list && \
    echo "deb https://security.debian.org/debian-security bookworm-security main contrib non-free non-free-firmware" >> /etc/apt/sources.list

RUN apt-get update --fix-missing && \
    apt-get -y upgrade && \
    apt-get install git && \
    apt-get install openssh-client

COPY . .

RUN cd ./docker-example \
    && git config --global url."ssh://git@gitlab.grandhoo.com/".insteadOf "https://gitlab.grandhoo.com/" \
    && go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w GOPRIVATE=gitlab.grandhoo.com \
    && go env \
    && go mod tidy \
    && go build -o docker-example main.go
