FROM docker.io/library/golang:1.20-alpine as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o inetmock ./cmd/inetmock

FROM alpine:3.17

WORKDIR /app

COPY --from=build /app/inetmock ./
COPY --from=build /app/config-container.yaml /etc/inetmock/config.yaml
COPY --from=build /app/assets/fakeFiles /var/lib/inetmock/fakeFiles
COPY --from=build /app/assets/demoCA /var/lib/inetmock/ca

RUN mkdir -p /var/run/inetmock

CMD /app/inetmock serve