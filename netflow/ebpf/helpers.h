#pragma once

#include "common.h"
#include "types.h"

/* 0x3FFF mask to check for fragment offset field */
#define IP_FRAGMENTED 65343
#define CONFIG_OPTIONS_COUNT 4

#ifndef memcpy
#define memcpy(dest, src, n) __builtin_memcpy((dest), (src), (n))
#endif

#ifndef __section
# define __section(NAME)                  \
    __attribute__((section(NAME), used))
#endif

static inline unsigned short checksum(unsigned short *buf, int bufsz) {
    unsigned long sum = 0;

    while (bufsz > 1) {
        sum += *buf;
        buf++;
        bufsz -= 2;
    }

    if (bufsz == 1) {
        sum += *(unsigned char *) buf;
    }

    sum = (sum & 0xffff) + (sum >> 16);
    sum = (sum & 0xffff) + (sum >> 16);

    return ~sum;
}

static inline struct tcphdr *extract_tcp_meta(struct l4_meta *meta, void *iph, void *data_end) {
    struct tcphdr *hdr = iph + sizeof(struct iphdr);
    if ((void *) hdr + sizeof(struct tcphdr) > data_end) {
        return NULL;
    }

    meta->src_port = bpf_ntohs(hdr->source);
    meta->dst_port = bpf_ntohs(hdr->dest);
    meta->transport_proto = L4_PROTO_TCP;

    if (hdr->fin) {
        meta->conn_state = CONN_STATE_CLOSING;
    } else if (hdr->rst) {
        meta->conn_state = CONN_STATE_FORCE_CLOSE;
    } else {
        meta->conn_state = CONN_STATE_OPEN;
    }

    return hdr;
}

static inline struct udphdr *extract_udp_meta(struct l4_meta *meta, void *iph, void *data_end) {
    struct udphdr *hdr = iph + sizeof(struct iphdr);
    if ((void *) hdr + sizeof(struct udphdr) > data_end) {
        return NULL;
    }

    meta->src_port = bpf_ntohs(hdr->source);
    meta->dst_port = bpf_ntohs(hdr->dest);
    meta->conn_state = CONN_STATE_UNKNOWN;
    meta->transport_proto = L4_PROTO_UDP;

    return hdr;
}

static inline bool extract_meta(struct l4_meta *meta, struct iphdr *iph, void *data_end) {
    __u32 ip_proto = iph->protocol;
    switch (ip_proto) {
        case IPPROTO_TCP:
            return extract_tcp_meta(meta, (void *) iph, data_end);
        case IPPROTO_UDP:
            return extract_udp_meta(meta, (void *) iph, data_end);
        default:
            return false;
    }
}