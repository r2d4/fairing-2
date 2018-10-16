FROM golang:1.10.1 as builder
WORKDIR /go/src/github.com/r2d4/notebuilder
COPY . .
RUN make install

FROM jupyter/base-notebook:1145fb1198b2
ENV JUPYTER_TOKEN=token
COPY --from=builder /go/bin/notebuilder .
