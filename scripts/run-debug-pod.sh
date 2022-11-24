#!/usr/bin/env bash

function join_container_to_ns() {
 sleep 5
 ip netns exec $(basename "$(podman inspect --type container --format '{{ .NetworkSettings.SandboxKey}}' inetmock)") iptables -t nat -A PREROUTING -p tcp -i eth0 -j REDIRECT
}

if [[ $(podman ps -f 'name=inetmock' --format '{{ .ID }}') != "" ]];
then
    podman rm -f inetmock
fi

go build -gcflags "all=-N -l" -o out/inetmock ./cmd/inetmock

# join_container_to_ns &

podman run \
    --rm \
    -ti \
    --ip 10.10.1.1 \
    --cap-add=CAP_NET_RAW \
    --cap-add=CAP_NET_ADMIN \
    --cap-add=CAP_SYS_ADMIN \
    --ulimit memlock=33554432:33554432 \
    --security-opt=seccomp=unconfined \
    --replace \
    --network=libvirt \
    -v /sys:/sys:ro \
    -v "$(pwd):/work" \
    -p 2345:2345 \
    -p 8080:80 \
    -p 9010:9010 \
    --name inetmock \
    inetmock-debug $@