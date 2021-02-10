FROM golang:1.15-alpine as build

WORKDIR /app

COPY ./ ./

RUN go build -o inetmock ./cmd/inetmock

FROM alpine:3.13

WORKDIR /app

COPY --from=build /app/inetmock ./
COPY --from=build /app/config-container.yaml /etc/inetmock/config.yaml
COPY --from=build /app/assets/fakeFiles /var/lib/inetmock/fakeFiles
COPY --from=build /app/assets/demoCA /var/lib/inetmock/ca

RUN mkdir -p /var/run/inetmock

CMD /app/inetmock serve