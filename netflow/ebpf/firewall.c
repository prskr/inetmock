#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include "helpers.h"

char LICENSE[] SEC("license") = "Dual MIT/GPL";

volatile const enum xdp_action DEFAULT_POLICY;
volatile const enum xdp_action DEFAULT_POLICY = XDP_DROP;

volatile const bool EMIT_UNMATCHED;
volatile const bool EMIT_UNMATCHED = false;

struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
} observed_packets SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(key_size, sizeof(struct conn_ident));
    __uint(value_size, sizeof(struct firewall_rule));
    __uint(max_entries, 1024);
} firewall_rules SEC(".maps");

struct packet_result {
    enum xdp_action result;
    bool emit;
    struct observed_packet pkt;
};

static inline bool handle_packet(struct xdp_md *ctx, struct packet_result *result) {
    void *data = (void *) (long) ctx->data;
    void *data_end = (void *) (long) ctx->data_end;
    struct ethhdr *eth = data;
    __u16 proto;

    if (data + sizeof(struct ethhdr) > data_end) {
        result->result = XDP_DROP;
        return true;
    }

    proto = bpf_ntohs(eth->h_proto);

    /* don't touch ARP requests */
    if (proto == ETH_P_ARP) {
        result->result = XDP_PASS;
        return true;
    }

    if (proto != ETH_P_IP && proto != ETH_P_IPV6) {
        return true;
    }

    struct iphdr *iph = data + sizeof(struct ethhdr);
    if ((void *) iph + sizeof(struct iphdr) > data_end) {
        result->result = XDP_DROP;
        return true;
    }

    /* do not support fragmented packets as L4 headers may be missing */
    if (iph->frag_off & IP_FRAGMENTED) {
        result->result = XDP_DROP;
        return true;
    }

    result->pkt.sourceIp = iph->saddr;
    result->pkt.destIp = iph->daddr;

    struct l4_meta meta;
    __builtin_memset(&meta, 0, sizeof(meta));

    if (!extract_meta(&meta, iph, data_end)) {
        return true;
    }

    result->emit = EMIT_UNMATCHED;

    struct conn_ident id;
    __builtin_memset(&id, 0, sizeof(struct conn_ident));

    id.port = meta.dst_port;
    id.transportProto = meta.transport_proto;

    struct firewall_rule *attached_rule;
    attached_rule = bpf_map_lookup_elem(&firewall_rules, &id);
    if (attached_rule) {
        result->emit = attached_rule->monitorTraffic;
        result->result = attached_rule->policy;
    }

    result->pkt.sourcePort = meta.src_port;
    result->pkt.destPort = meta.dst_port;
    result->pkt.transportProto = meta.transport_proto;

    return false;
}

SEC("xdp/perf")
int xdp_ingress_perf(struct xdp_md *ctx) {
    struct packet_result result;
    __builtin_memset(&result, 0, sizeof(result));

    result.result = DEFAULT_POLICY;

    bool error_occurred = handle_packet(ctx, &result);

    if (error_occurred || !result.emit) {
        return result.result;
    }

    long submissionResult = bpf_perf_event_output(ctx, &observed_packets, BPF_F_CURRENT_CPU, &result.pkt, sizeof(struct observed_packet));

    if (submissionResult != 0) {
        bpf_printk("Failed to submit observed packet: %d\n", submissionResult);
    }

    return result.result;
}

SEC("xdp/ring")
int xdp_ingress_ring(struct xdp_md *ctx) {
    struct packet_result result;
    __builtin_memset(&result, 0, sizeof(result));

    result.result = DEFAULT_POLICY;

    bool errorOccurred = handle_packet(ctx, &result);

    if (errorOccurred || !result.emit) {
        return result.result;
    }

    long submissionResult = bpf_ringbuf_output(&observed_packets, &result.pkt, sizeof(result.pkt), 0);

    if (submissionResult != 0) {
        bpf_printk("Failed to submit observed packet: %d\n", submissionResult);
    }

    return result.result;
}

SEC("xdp/mock")
int xdp_mock(struct xdp_md *ctx) {
    return DEFAULT_POLICY;
}