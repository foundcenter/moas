FROM golang:1.8.0
COPY bin/moas /go/bin/moas

RUN groupadd -r gouser && useradd -r -g gouser gouser
USER gouser

CMD moas
