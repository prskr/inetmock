package audit

import (
	"crypto/tls"

	v1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
)

var (
	tlsToEntity = map[uint16]v1.TLSVersion{
		tls.VersionTLS10: v1.TLSVersion_TLS_VERSION_TLS10,
		tls.VersionTLS11: v1.TLSVersion_TLS_VERSION_TLS11,
		tls.VersionTLS12: v1.TLSVersion_TLS_VERSION_TLS12,
		tls.VersionTLS13: v1.TLSVersion_TLS_VERSION_TLS13,
	}
)

type TLSDetails struct {
	Version     string
	CipherSuite string
	ServerName  string
}

func TLSVersionToEntity(version uint16) v1.TLSVersion {
	if v, known := tlsToEntity[version]; known {
		return v
	}
	return v1.TLSVersion_TLS_VERSION_UNSPECIFIED
}

func NewTLSDetailsFromProto(entity *v1.TLSDetailsEntity) *TLSDetails {
	if entity == nil {
		return nil
	}

	return &TLSDetails{
		Version:     entity.GetVersion().String(),
		CipherSuite: entity.GetCipherSuite(),
		ServerName:  entity.GetServerName(),
	}
}

func (d TLSDetails) ProtoMessage() *v1.TLSDetailsEntity {
	return &v1.TLSDetailsEntity{
		Version:     v1.TLSVersion(v1.TLSVersion_value[d.Version]),
		CipherSuite: d.CipherSuite,
		ServerName:  d.ServerName,
	}
}
