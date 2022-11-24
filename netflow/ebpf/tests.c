#include <linux/bpf.h>
#include <linux/pkt_cls.h>
#include <bpf/bpf_helpers.h>

#include "helpers.h"

char LICENSE[] SEC("license") = "Dual MIT/GPL";

struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
} observed_packets SEC(".maps");

SEC("classifier/perf-tests")
int emit_test_events_perf(struct __sk_buff *skb) {
    __u16 ports[6] = {53, 80, 443, 853, 3128, 8080};

    struct observed_packet pkt;
    __builtin_memset(&pkt, 0, sizeof(pkt));

    pkt.transportProto = L4_PROTO_TCP;
    pkt.sourcePort = 23876;
    pkt.sourceIp = 16777343;
    pkt.destIp = 16777343;

#pragma clang loop unroll(full)
    for (int i = 0; i < sizeof ports; i++) {
        pkt.destPort = ports[i];
        bpf_perf_event_output(skb, &observed_packets, BPF_F_CURRENT_CPU, &pkt, sizeof(struct observed_packet));
    }
    return TC_ACT_OK;
}

SEC("classifier/ringbuf-tests")
int emit_test_events_ring_buf(struct __sk_buff *skb) {
    __u16 ports[6] = {53, 80, 443, 853, 3128, 8080};

    struct observed_packet pkt;
    __builtin_memset(&pkt, 0, sizeof(pkt));

    pkt.transportProto = L4_PROTO_TCP;
    pkt.sourcePort = 23876;
    pkt.sourceIp = 16777343;
    pkt.destIp = 16777343;

#pragma clang loop unroll(full)
    for (int i = 0; i < sizeof ports; i++) {
        pkt.destPort = ports[i];
        bpf_ringbuf_output(&observed_packets, &pkt, sizeof(struct observed_packet), 0);
    }
    return TC_ACT_OK;
}