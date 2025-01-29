FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN make build

FROM ubuntu:24.04

WORKDIR /app

COPY --from=builder /app/bin/yukichan-bot-linux-amd64-v1.0.0 /app/yukichanbot

RUN apt-get update && apt-get install -y ca-certificates \
    git python-is-python3 python3-pip inkscape
RUN pip install python-chess --break-system-packages && \
    git clone https://github.com/dn1z/pgn2gif.git /tmp/pgn2gif && \
    cd /tmp/pgn2gif && python setup.py install

ENTRYPOINT ["/app/yukichanbot"]

