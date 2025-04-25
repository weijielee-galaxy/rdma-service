FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive

ARG OFED=MLNX_OFED_LINUX-23.10-1.1.9.0-ubuntu20.04-x86_64
ENV OFED_TAR ${OFED}.tgz
# COPY ${OFED_TAR} /tmp

ARG GO_VER=1.22.6
ARG GO_TAR=go${GO_VER}.linux-amd64.tar.gz


RUN apt-get update && \
    apt-get -y install wget apt-utils vim tmux libcap2 iproute2 iperf iperf3 python3 tcpdump lldpad net-tools git build-essential && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN wget -O /tmp/${OFED_TAR} https://content.mellanox.com/ofed/MLNX_OFED-23.10-1.1.9.0/${OFED_TAR}

WORKDIR /tmp
RUN set -ex && \
    tar zxvf ${OFED}.tgz &&\
    cd ${OFED} && \
    ./mlnxofedinstall --hpc --user-space-only --force && \
    mkdir -p /data && \
    cp /tmp/${OFED}.tgz /data && \
    cd /tmp && \
    rm -rf /tmp/${OFED}

RUN mkdir -p /data && \
    wget -O /tmp/${GO_TAR} https://go.dev/dl/${GO_TAR} && \
    tar -C /usr/local -xzf /tmp/${GO_TAR} && \
    cp /tmp/${GO_TAR} /data && \
    rm -rf /tmp/${GO_TAR}

ENV PATH=$PATH:/usr/local/go/bin/

COPY . /app
RUN cd /app && \
    go mod tidy && \
    make build && \
    cp ./bin/rdma-service /usr/local/bin/rdma-service

WORKDIR /home/