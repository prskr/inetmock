FROM golang:1.14-alpine as build

# Create appuser and group.
ARG USER=inetmock
ARG GROUP=inetmock
ARG USER_ID=10001
ARG GROUP_ID=10001

ENV CGO_ENABLED=0

# Prepare build stage - can be cached
WORKDIR /work
RUN apk add -U --no-cache \
        make protoc gcc musl-dev libcap && \
    addgroup -S -g "${GROUP_ID}" "${GROUP}" && \
    adduser \
        --disabled-password \
        --gecos "" \
        --home "/nonexistent" \
        --shell "/sbin/nologin" \
        --no-create-home \
        -G "${GROUP}" \
        --uid "${USER_ID}" \
        "${USER}"

# Fetch dependencies
COPY Makefile go.mod go.sum ./
RUN go mod download && \
    go get -u github.com/golang/mock/mockgen@latest && \
    go install github.com/golang/protobuf/protoc-gen-go

COPY ./ ./

# Build binaries
RUN make CONTAINER=yes && \
    mkdir -p /usr/lib/inetmock/bin/ && \
    chown $USER:$GROUP inetmock imctl && \
    mv inetmock imctl /usr/lib/inetmock/bin/ && \
    setcap 'cap_net_bind_service=+ep' /usr/lib/inetmock/bin/inetmock

# Runtime layer

FROM alpine:3.12

# Create appuser and group.
ARG USER=inetmock
ARG GROUP=inetmock
ARG USER_ID=10001
ARG GROUP_ID=10001

COPY --from=build /etc/group /etc/passwd /etc/
COPY --from=build /usr/lib/inetmock/bin /usr/lib/inetmock/bin
COPY config-container.yaml /etc/inetmock/config.yaml

RUN mkdir -p /var/run/inetmock /var/lib/inetmock/certs /usr/lib/inetmock && \
    chown -R $USER:$GROUP /var/run/inetmock /var/lib/inetmock /usr/lib/inetmock

RUN ln -s /usr/lib/inetmock/bin/inetmock /usr/bin/inetmock && \
    ln -s /usr/lib/inetmock/bin/imctl /usr/bin/imctl

USER $USER

VOLUME [ "/var/lib/inetmock/ca", "/var/lib/inetmock/certs" ]

ENTRYPOINT ["inetmock"]
