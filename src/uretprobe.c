// go:build ignore

#include "common.h"

#include <bpf/bpf_endian.h>
#include <bpf/bpf_tracing.h>

char __license[] SEC("license") = "Dual MIT/GPL";

struct event {
    u32 pid;
    u32 uid;
    u8 line[496];
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 256 * 1024 /* 256 KB */);
} events SEC(".maps");

// Force emitting struct event into the ELF.
const struct event* unused __attribute__((unused));

SEC("uretprobe/zsh_readline")
int uretprobe_zsh_readline(struct pt_regs* ctx)
{
    struct event event;

    event.pid = bpf_get_current_pid_tgid();
    event.uid = bpf_get_current_uid_gid() & 0xFFFFFFFF; // lower 32 bits contain uid
    bpf_probe_read(&event.line, sizeof(event.line), (void*)PT_REGS_RC(ctx));

    bpf_ringbuf_output(&events, &event, sizeof(event), 0);

    return 0;
}
