FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM ubuntu:26.04

ARG DEBIAN_FRONTEND=noninteractive

WORKDIR /app

COPY --from=builder /app/bin/yukichan-bot-linux-amd64-v1.0.0 /app/yukichanbot

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        git \
        inkscape \
        python-is-python3 \
        python3-pip \
        python3-setuptools && \
    PIP_ROOT_USER_ACTION=ignore pip install --no-cache-dir python-chess --break-system-packages && \
    git clone --depth 1 https://github.com/dn1z/pgn2gif.git /tmp/pgn2gif && \
    cd /tmp/pgn2gif && \
    python setup.py install && \
    apt-get purge -y --auto-remove git python3-pip python3-setuptools && \
    rm -rf /root/.cache /var/lib/apt/lists/* /tmp/pgn2gif

ENTRYPOINT ["/app/yukichanbot"]
