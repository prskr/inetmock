package netutils

import (
	"encoding/binary"
	"net"
	"unsafe"
)

func Uint32ToIP(i uint32) net.IP {
	bytes := (*[4]byte)(unsafe.Pointer(&i))[:]
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

func IPToInt32(ip net.IP) uint32 {
	v4 := ip.To4()
	result := binary.BigEndian.Uint32(v4)
	return result
}

func IPAddressesToBytes(addresses []net.IP) (result [][]byte) {
	for i := range addresses {
		result = append(result, addresses[i])
	}
	return
}

func BytesToIPAddresses(input [][]byte) (result []net.IP) {
	result = make([]net.IP, 0, len(input))

	for idx := range input {
		result = append(result, net.IP(input[idx]))
	}

	return result
}
