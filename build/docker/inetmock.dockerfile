# Runtime layer
FROM docker.io/alpine:3.14

# Create appuser and group.
ARG USER=inetmock
ARG GROUP=inetmock
ARG USER_ID=10001
ARG GROUP_ID=10001

RUN addgroup -S -g "${GROUP_ID}" "${GROUP}" && \
    adduser \
        --disabled-password \
        --gecos "" \
        --home "/nonexistent" \
        --shell "/sbin/nologin" \
        --no-create-home \
        -G "${GROUP}" \
        --uid "${USER_ID}" \
        "${USER}"

COPY --chown=$USER:$GROUP inetmock imctl /usr/lib/inetmock/bin/
COPY --chown=$USER:$GROUP assets/fakeFiles /var/lib/inetmock/fakeFiles/
COPY --chown=$USER:$GROUP assets/demoCA /var/lib/inetmock/ca
COPY config-container.yaml /etc/inetmock/config.yaml

RUN mkdir -p /var/run/inetmock /var/lib/inetmock/ca /var/lib/inetmock/certs /var/lib/inetmock/data /usr/lib/inetmock && \
    chown -R $USER:$GROUP /var/run/inetmock /var/lib/inetmock /usr/lib/inetmock && \
    apk add -U --no-cache libcap

RUN ln -s /usr/lib/inetmock/bin/inetmock /usr/bin/inetmock && \
    ln -s /usr/lib/inetmock/bin/imctl /usr/bin/imctl && \
    setcap 'cap_net_raw,cap_net_bind_service=eip' /usr/lib/inetmock/bin/inetmock

HEALTHCHECK --interval=5s --timeout=1s \
    CMD imctl --socket-path /var/run/inetmock/inetmock.sock health container

USER $USER

VOLUME [ "/var/lib/inetmock/data" ]

ENTRYPOINT ["inetmock"]
