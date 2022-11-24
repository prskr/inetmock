#include <stdbool.h>
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
#include <linux/if_ether.h>
#include <linux/bpf.h>
#include <linux/in.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <linux/udp.h>

enum bpf_func_id___x {
    BPF_FUNC_snprintf___x = 42 /* avoid zero */
};
#define advanced_formatting_available  (bpf_core_enum_value_exists(enum bpf_func_id___x, BPF_FUNC_snprintf___x))