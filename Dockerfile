FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM ubuntu:26.04

ARG DEBIAN_FRONTEND=noninteractive

WORKDIR /app

COPY --from=builder /app/bin/yukichan-bot-linux-amd64-v1.0.0 /app/yukichanbot

RUN apt-get update && apt-get install -y ca-certificates \
    git python-is-python3 python3-pip inkscape python3-setuptools

RUN pip install python-chess pillow \
    -i https://pypi.tuna.tsinghua.edu.cn/simple \
    --break-system-packages

COPY ./vend/pgn2gif /tmp/pgn2gif

RUN cd /tmp/pgn2gif && python setup.py install

ENTRYPOINT ["/app/yukichanbot"]
