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
        make protoc gcc musl-dev && \
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
    go get -u github.com/abice/go-enum && \
    go install github.com/golang/protobuf/protoc-gen-go

COPY ./ ./

# Build binaries
RUN make CONTAINER=yes

# Runtime layer

FROM alpine:3.12

# Create appuser and group.
ARG USER=inetmock
ARG GROUP=inetmock
ARG USER_ID=10001
ARG GROUP_ID=10001

COPY --from=build /etc/group /etc/passwd /etc/
COPY --from=build --chown=$USER:$GROUP /work/inetmock /work/imctl /usr/lib/inetmock/bin/
COPY --chown=$USER:$GROUP ./assets/fakeFiles/ /var/lib/inetmock/fakeFiles/
COPY config-container.yaml /etc/inetmock/config.yaml

RUN mkdir -p /var/run/inetmock /var/lib/inetmock/certs /usr/lib/inetmock && \
    chown -R $USER:$GROUP /var/run/inetmock /var/lib/inetmock /usr/lib/inetmock && \
    apk add -U --no-cache libcap

RUN ln -s /usr/lib/inetmock/bin/inetmock /usr/bin/inetmock && \
    ln -s /usr/lib/inetmock/bin/imctl /usr/bin/imctl && \
    setcap 'cap_net_bind_service=+ep' /usr/lib/inetmock/bin/inetmock

HEALTHCHECK --interval=5s --timeout=1s \
    CMD imctl --socket-path /var/run/inetmock/inetmock.sock health container

USER $USER

VOLUME [ "/var/lib/inetmock/ca", "/var/lib/inetmock/certs" ]

ENTRYPOINT ["inetmock"]
