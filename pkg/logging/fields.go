package logging

import (
	"net"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type IPArray []net.IP

func (ss IPArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range ss {
		arr.AppendString(ss[i].String())
	}
	return nil
}

func IP(key string, ip net.IP) zap.Field {
	return zap.Stringer(key, ip)
}

func IPs(key string, ips []net.IP) zap.Field {
	return zap.Array(key, IPArray(ips))
}
