FROM golang:1.21 as builder

WORKDIR /root
# build benchmark
RUN apt-get install -y git
COPY admin/benchmarker /root/
RUN go build \
    && mv main benchmark
