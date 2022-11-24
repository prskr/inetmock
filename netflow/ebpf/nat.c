#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/pkt_cls.h>
#include <linux/ip.h>

#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>
#include "helpers.h"

char LICENSE[] SEC("license") = "Dual MIT/GPL";

volatile const __u32 INTERFACE_IP;
volatile const __u32 INTERFACE_IP = 16777343;

struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __type(key, __u32);
    __type(value, __u32);
    __uint(max_entries, CONFIG_OPTIONS_COUNT);
} nat_config SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(key_size, sizeof(struct conn_ident));
    __uint(value_size, sizeof(struct conn_meta));
    __uint(max_entries, 1024);
} conn_track SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(key_size, sizeof(struct conn_ident));
    __uint(value_size, sizeof(struct conn_meta));
    __uint(max_entries, 1024);
} nat_translations SEC(".maps");

#ifdef TESTING

static int capture_fake_packets() {
#pragma clang loop unroll(full)
    for (int i = 0; i < 5; i++) {
        submit_packet(16777343, 16777343, 12345, 135, TCP);
        submit_packet(16777343, 16777343, 32458, 1053, UDP);
        submit_packet(16777343, 16777343, 32789, 9090, TCP);
        submit_packet(838860810, 16777226, 20987, 3128, TCP);
        submit_packet(2058528960, 28485824, 11873, 3128, TCP);
    }

    return 0;
}

#endif

SEC("classifier/ingress")
int ingress(struct __sk_buff *skb) {
    void *data = (void *) (long) skb->data;
    void *data_end = (void *) (long) skb->data_end;
    struct ethhdr *eth = data;
    __u32 eth_proto;

    if (data + sizeof(struct ethhdr) > data_end) {
        return TC_ACT_SHOT;
    }

    eth_proto = bpf_ntohs(eth->h_proto);
    if (eth_proto != ETH_P_IP && eth_proto != ETH_P_IPV6) {
        return TC_ACT_OK;
    }

    struct iphdr *iph = data + sizeof(*eth);
    if ((void *) iph + sizeof(*iph) > data_end) {
        return TC_ACT_OK;
    }

    /* do not support fragmented packets as L4 headers may be missing */
    if (iph->frag_off & IP_FRAGMENTED) {
        return TC_ACT_OK;
    }

    if(iph->daddr == INTERFACE_IP) {
        return TC_ACT_OK;
    }

    struct l4_meta meta;
    __builtin_memset(&meta, 0, sizeof(meta));
    if (!extract_meta(&meta, iph, data_end)) {
        return TC_ACT_OK;
    }

    struct conn_ident id;
    __builtin_memset(&id, 0, sizeof(struct conn_ident));

    id.ip = iph->daddr;
    id.port = meta.dst_port;
    id.transportProto = meta.transport_proto;

    struct nat_rule *rule;
    rule = bpf_map_lookup_elem(&nat_translations, &id);

    if (!rule) {
        id.ip = 0;
        rule = bpf_map_lookup_elem(&nat_translations, &id);
    }

    if (!rule) {
        return TC_ACT_OK;
    }

    __u32 *epoch;
    int config_key = NAT_CURRENT_EPOCH;
    epoch = bpf_map_lookup_elem(&nat_config, &config_key);

    if (!epoch) {
        bpf_printk("No epoch set\n");
        return TC_ACT_OK;
    }

    struct conn_ident src;
    __builtin_memset(&src, 0, sizeof(src));
    src.ip = iph->saddr;
    src.port = meta.src_port;
    src.transportProto = meta.transport_proto;

    struct conn_meta dst;
    __builtin_memset(&dst, 0, sizeof(dst));

    dst.ip = iph->daddr;
    dst.port = meta.dst_port;
    dst.lastObserved = *epoch;
    dst.transportProto = meta.transport_proto;

    struct conn_meta *known = bpf_map_lookup_elem(&conn_track, &src);
    if (known && meta.conn_state == CONN_STATE_FORCE_CLOSE) {
        bpf_map_delete_elem(&conn_track, &src);
    } else {
        bpf_map_update_elem(&conn_track, &src, &dst, BPF_ANY);
    }

    iph->daddr = rule->targetIp;
    iph->tos = 7 << 2;
    iph->check = 0;
    iph->check = checksum((unsigned short *) iph, sizeof(struct iphdr));

    return TC_ACT_OK;
}

SEC("classifier/egress")
int egress(struct __sk_buff *skb) {
    void *data = (void *) (long) skb->data;
    void *data_end = (void *) (long) skb->data_end;
    struct ethhdr *eth = data;
    __u32 eth_proto;

    if (data + sizeof(struct ethhdr) > data_end) {
        return TC_ACT_OK;
    }

    eth_proto = bpf_ntohs(eth->h_proto);
    if (eth_proto != 0 && eth_proto != ETH_P_IP && eth_proto != ETH_P_IPV6) {
        return TC_ACT_OK;
    }

    struct iphdr *iph = data + sizeof(*eth);
    if ((void *) iph + sizeof(*iph) > data_end) {
        return TC_ACT_OK;
    }

    /* do not support fragmented packets as L4 headers may be missing */
    if (iph->frag_off & IP_FRAGMENTED) {
        return TC_ACT_OK;
    }

    struct l4_meta meta;

    if (!extract_meta(&meta, iph, data_end)) {
        return TC_ACT_OK;
    }

    struct conn_ident dst;
    __builtin_memset(&dst, 0, sizeof(dst));
    dst.ip = iph->daddr;
    dst.port = meta.dst_port;
    dst.transportProto = meta.transport_proto;

    struct conn_meta *orig_src = bpf_map_lookup_elem(&conn_track, &dst);
    if (!orig_src) {
        return TC_ACT_OK;
    }

    if (meta.conn_state == CONN_STATE_CLOSING || meta.conn_state == CONN_STATE_FORCE_CLOSE) {
        bpf_map_delete_elem(&conn_track, &dst);
    }

    iph->saddr = orig_src->ip;

    iph->tos = 7 << 2;
    iph->check = 0;
    iph->check = checksum((unsigned short *) iph, sizeof(struct iphdr));

    return TC_ACT_OK;
}

SEC("classifier/mock")
int nat_mock(struct __sk_buff *skb) {
    return TC_ACT_OK;
}