FROM docker.io/alpine:3.14 as builder

RUN touch /tmp/.keep

# Runtime layer
FROM gcr.io/distroless/static:debug-nonroot

USER nonroot:nonroot

COPY --from=builder --chown=nonroot:nonroot /tmp/.keep /var/lib/inetmock/data/pcap/.keep
COPY --from=builder --chown=nonroot:nonroot /tmp/.keep /var/lib/inetmock/data/audit/.keep
COPY --from=builder --chown=nonroot:nonroot /tmp/.keep /var/lib/inetmock/data/certs/.keep
COPY --from=builder --chown=nonroot:nonroot /tmp/.keep /var/run/inetmock/inetmock.sock

COPY --chown=nonroot:nonroot inetmock imctl /usr/lib/inetmock/bin/
COPY --chown=nonroot:nonroot assets/fakeFiles /var/lib/inetmock/fakeFiles/
COPY --chown=nonroot:nonroot assets/demoCA /var/lib/inetmock/ca
COPY config-container.yaml /etc/inetmock/config.yaml

ENTRYPOINT ["/usr/lib/inetmock/bin/inetmock"]
CMD ["serve"]