FROM jupyter/base-notebook:1145fb1198b2 as base
RUN pip install astunparse

FROM golang:1.10.1 as builder
WORKDIR /go/src/github.com/r2d4/fairing
COPY vendor ./vendor
COPY Makefile .

COPY cmd ./cmd
COPY pkg ./pkg

RUN make install

FROM base
ENV JUPYTER_TOKEN=token
COPY --from=builder /go/bin/fairing /usr/bin/fairing
COPY --chown=jovyan:users hack/*.py  /home/jovyan/work/lib/

