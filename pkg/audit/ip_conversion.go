package audit

import (
	"encoding/binary"
	"math/big"
	"net"
)

func ipv4ToUint32(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func ipv6ToBytes(ip net.IP) uint64 {
	ipv6 := big.NewInt(0)
	ipv6.SetBytes(ip)
	return ipv6.Uint64()
}

func uint32ToIP(i uint32) (ip net.IP) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, i)
	ip = buf
	ip = ip.To4()
	return
}

func uint64ToIP(i uint64) (ip net.IP) {
	ip = big.NewInt(int64(i)).FillBytes(make([]byte, 16))
	return
}
