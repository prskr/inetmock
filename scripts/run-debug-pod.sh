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

join_container_to_ns &

podman run \
    --rm \
    -ti \
    --ip 10.10.1.1 \
    --cap-add=CAP_NET_RAW \
    --cap-add=CAP_NET_BIND_SERVICE \
    --cap-add=CAP_NET_ADMIN \
    --security-opt=seccomp=unconfined \
    --replace \
    --network=libvirt \
    -v "$(pwd):/work" \
    --name inetmock \
    inetmock-debug