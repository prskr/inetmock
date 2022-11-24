package netflow

import (
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"net/netip"
)

var (
	_ encoding.BinaryMarshaler    = ConnIdent{}
	_ BinaryCollectionUnmarshaler = (*ConnIdent)(nil)

	_ encoding.BinaryMarshaler    = ConnMeta{}
	_ BinaryCollectionUnmarshaler = (*ConnMeta)(nil)
)

type ConnIdent struct {
	Addr      netip.Addr
	Port      uint16
	Transport Protocol
}

func (ConnIdent) BinarySize() int {
	return connIdentBinaryLength
}

func (t ConnIdent) MarshalBinary() (data []byte, err error) {
	data = make([]byte, connIdentBinaryLength)
	err = t.MarshalBinaryTo(data)
	return data, err
}

func (t ConnIdent) MarshalBinaryTo(data []byte) (err error) {
	if l := len(data); l != connIdentBinaryLength {
		return fmt.Errorf("ConnIdent does not marshal into %d bytes", l)
	}

	if addrData := t.Addr.AsSlice(); len(addrData) > net.IPv4len {
		return errors.New("only IPv4 addresses are support for now")
	} else {
		copy(data[:4], addrData)
	}

	binary.LittleEndian.PutUint16(data[4:6], t.Port)
	binary.LittleEndian.PutUint32(data[8:], uint32(t.Transport))

	return nil
}

func (t *ConnIdent) UnmarshalBinary(data []byte) error {
	if len(data) < connIdentBinaryLength {
		return errors.New("data length does not match")
	}
	ipBytes := make([]byte, net.IPv4len)
	copy(ipBytes, data[:4])

	if a, ok := netip.AddrFromSlice(ipBytes); ok {
		t.Addr = a
	}

	t.Port = binary.LittleEndian.Uint16(data[4:6])
	t.Transport = Protocol(binary.LittleEndian.Uint32(data[8:]))
	return nil
}

type ConnMeta struct {
	Addr         netip.Addr
	Port         uint16
	Transport    Protocol
	LastObserved uint32
}

func (ConnMeta) BinarySize() int {
	return connMetaBinaryLength
}

func (t ConnMeta) MarshalBinary() (data []byte, err error) {
	data = make([]byte, connMetaBinaryLength)
	err = t.MarshalBinaryTo(data)
	return data, err
}

func (t ConnMeta) MarshalBinaryTo(data []byte) (err error) {
	if l := len(data); l != connMetaBinaryLength {
		return fmt.Errorf("ConnMeta does not marshal to %d bytes", l)
	}
	copy(data[:4], t.Addr.AsSlice())

	binary.LittleEndian.PutUint16(data[4:6], t.Port)
	binary.LittleEndian.PutUint32(data[8:12], uint32(t.Transport))
	binary.LittleEndian.PutUint32(data[12:], t.LastObserved)

	return nil
}

func (t *ConnMeta) UnmarshalBinary(data []byte) error {
	if len(data) < connMetaBinaryLength {
		return errors.New("data length does not match")
	}
	ipBytes := make([]byte, net.IPv4len)
	copy(ipBytes, data[:4])

	if a, ok := netip.AddrFromSlice(ipBytes); ok {
		t.Addr = a
	}

	t.Port = binary.LittleEndian.Uint16(data[4:6])
	t.Transport = Protocol(binary.LittleEndian.Uint32(data[8:12]))
	t.LastObserved = binary.LittleEndian.Uint32(data[12:])
	return nil
}
