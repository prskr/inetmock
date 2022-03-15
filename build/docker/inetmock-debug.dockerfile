FROM registry.gitlab.com/inetmock/ci-image/go

ENV GOPROXY=""
ENV GOSUMDB=""

RUN go install github.com/go-delve/delve/cmd/dlv@latest && \
    echo 'dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec $@' >> /usr/local/bin/exec.sh && \
    echo 'dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient $@' >> /usr/local/bin/debug.sh && \
    chmod +x /usr/local/bin/*.sh

WORKDIR /work

EXPOSE 2345