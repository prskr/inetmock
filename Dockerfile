FROM golang:1.14-alpine as build

# Create appuser.
ARG USER=inetmock
ARG USER_ID=10001

ENV CGO_ENABLED=0

# Prepare build stage - can be cached
WORKDIR /work
RUN apk add -U --no-cache \
        make protoc gcc musl-dev && \
    adduser \
        --disabled-password \
        --gecos "" \
        --home "/nonexistent" \
        --shell "/sbin/nologin" \
        --no-create-home \
        --uid "${USER_ID}" \
        "${USER}"

# Fetch dependencies
COPY Makefile go.mod go.sum ./
RUN go mod download && \
    go get -u github.com/golang/mock/mockgen@latest && \
    go install github.com/golang/protobuf/protoc-gen-go

COPY ./ ./

# Build binary and plugins
RUN make CONTAINER=yes

# Runtime layer

FROM scratch

ENV INETMOCK_PLUGINS_DIRECTORY=/app/plugins/

WORKDIR /app

COPY --from=build /etc/passwd /etc/group /etc/
COPY --from=build --chown=$USER /work/inetmock ./
COPY --from=build --chown=$USER /work/*.so ./plugins/

USER $USER:$USER

ENTRYPOINT ["/app/inetmock"]
