FROM golang:1.14-buster as build

# Create appuser.
ARG USER=inetmock
ARG USER_ID=10001

ENV CGO_ENABLED=0

# Prepare build stage - can be cached
WORKDIR /work
RUN apt-get update && \
    apt-get install -y --no-install-recommends make gcc && \
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
RUN go mod download

COPY ./ ./

# Build binary and plugins
RUN make CONTAINER=yes

# Runtime layer

FROM scratch

WORKDIR /app

COPY --from=build /etc/passwd /etc/group /etc/
COPY --from=build --chown=$USER /work/inetmock ./
COPY --from=build --chown=$USER /work/plugins/ ./plugins/

USER $USER:$USER

ENTRYPOINT ["/app/inetmock"]