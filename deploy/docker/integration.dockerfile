FROM docker.io/alpine:3.14

WORKDIR /app

COPY out/inetmock ./
COPY assets/fakeFiles /var/lib/inetmock/fakeFiles/
COPY assets/demoCA /var/lib/inetmock/ca
COPY config-container.yaml /etc/inetmock/config.yaml

RUN mkdir -p /var/run/inetmock /var/lib/inetmock/ca /var/lib/inetmock/certs /var/lib/inetmock/data /usr/lib/inetmock

ENTRYPOINT ["/app/inetmock"]
CMD ["serve"]
