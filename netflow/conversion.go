package netflow

import (
	"encoding/binary"
	"net/netip"
)

func ipAddr2int(addr netip.Addr) uint32 {
	b := addr.AsSlice()
	reverse(b)
	return binary.BigEndian.Uint32(b)
}

func reverse(input []byte) {
	for i := 0; i < len(input)/2; i++ {
		input[i], input[len(input)-1-i] = input[len(input)-1-i], input[i]
	}
}
