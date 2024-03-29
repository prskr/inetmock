ARG BASE_IMAGE=code.icb4dc0.de/inetmock/ci-images/go-ci

FROM ${BASE_IMAGE}

RUN go install github.com/go-delve/delve/cmd/dlv@latest && \
    echo 'dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec $@' >> /usr/local/bin/exec.sh && \
    echo 'dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient $@' >> /usr/local/bin/debug.sh && \
    chmod +x /usr/local/bin/*.sh

WORKDIR /work

ENV GOPROXY=https://goproxy.io,direct

EXPOSE 2345