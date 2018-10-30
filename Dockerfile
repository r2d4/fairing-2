FROM jupyter/base-notebook:1145fb1198b2 as base
RUN pip install astunparse


FROM golang:1.10.1 as builder

ENV DOCKER_CREDENTIAL_GCR_VERSION=1.4.3
RUN curl -LO https://github.com/GoogleCloudPlatform/docker-credential-gcr/releases/download/v${DOCKER_CREDENTIAL_GCR_VERSION}/docker-credential-gcr_linux_amd64-${DOCKER_CREDENTIAL_GCR_VERSION}.tar.gz && \
    tar -zxvf docker-credential-gcr_linux_amd64-${DOCKER_CREDENTIAL_GCR_VERSION}.tar.gz && \
    mv docker-credential-gcr /usr/bin/docker-credential-gcr && \
    rm docker-credential-gcr_linux_amd64-${DOCKER_CREDENTIAL_GCR_VERSION}.tar.gz && \
    chmod +x /usr/bin/docker-credential-gcr && \
    docker-credential-gcr configure-docker

WORKDIR /go/src/github.com/kubeflow/fairing
COPY vendor ./vendor
COPY Makefile .

COPY cmd ./cmd
COPY pkg ./pkg

RUN make install

FROM base
ENV JUPYTER_TOKEN=token
COPY --from=builder --chown=jovyan:users /root/.docker/config.json /home/jovyan/.docker/config.json
COPY --from=builder /usr/bin/docker-credential-gcr /usr/bin/docker-credential-gcr

COPY --from=builder /go/bin/fairing /usr/bin/fairing
COPY --chown=jovyan:users fairing-py  /home/jovyan/work/lib/
