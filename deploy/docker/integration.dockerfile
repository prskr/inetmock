ARG BASE_IMAGE

FROM ${BASE_IMAGE}

COPY assets/fakeFiles /var/lib/inetmock/fakeFiles/
COPY assets/demoCA /var/lib/inetmock/ca
COPY testdata/config-integration.yaml /etc/inetmock/config.yaml

VOLUME /var/lib/inetmock/data

USER root

CMD ["serve", "--config=/etc/inetmock/config.yaml"]
