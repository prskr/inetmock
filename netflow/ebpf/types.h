#define FW_DEFAULT_PACKET_POLICY 0
#define FW_EMIT_EVENT 1

#define NAT_CURRENT_EPOCH 0

enum transport_proto {
    L4_PROTO_UNKNOWN,
    L4_PROTO_TCP,
    L4_PROTO_UDP
};

struct l4_meta {
    __u16 src_port;
    __u16 dst_port;

    enum {
        CONN_STATE_UNKNOWN,
        CONN_STATE_OPEN,
        CONN_STATE_CLOSING,
        CONN_STATE_FORCE_CLOSE
    } conn_state;

    enum transport_proto transport_proto;
};

struct observed_packet {
    __u32 sourceIp;
    __u32 destIp;
    __u16 sourcePort;
    __u16 destPort;
    enum transport_proto transportProto;
};

struct conn_ident {
    __u32 ip;
    __u16 port;
    enum transport_proto transportProto;
};

struct conn_meta {
    __u32 ip;
    __u16 port;
    enum transport_proto transportProto;
    __u32 lastObserved;
};

struct firewall_rule {
    enum xdp_action policy;
    bool monitorTraffic;
};

struct nat_rule {
    __u32 targetIp;
};