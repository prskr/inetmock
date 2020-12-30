package audit

import (
	"crypto/tls"

	"google.golang.org/protobuf/proto"
)

type TLSDetails struct {
	Version     uint16
	CipherSuite uint16
	ServerName  string
}

func (d TLSDetails) ProtoMessage() proto.Message {
	var version TLSVersion

	switch d.Version {
	case tls.VersionTLS10:
		version = TLSVersion_TLS10
	case tls.VersionTLS11:
		version = TLSVersion_TLS11
	case tls.VersionTLS12:
		version = TLSVersion_TLS12
	case tls.VersionTLS13:
		version = TLSVersion_TLS13
	default:
		version = TLSVersion_SSLv30
	}

	return &TLSDetailsEntity{
		Version:     version,
		CipherSuite: tls.CipherSuiteName(d.CipherSuite),
		ServerName:  d.ServerName,
	}
}
