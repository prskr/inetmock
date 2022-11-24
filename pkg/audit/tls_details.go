package audit

import (
	"crypto/tls"

	auditv1 "inetmock.icb4dc0.de/inetmock/pkg/audit/v1"
)

var tlsToEntity = map[uint16]auditv1.TLSVersion{
	tls.VersionTLS10: auditv1.TLSVersion_TLS_VERSION_TLS10,
	tls.VersionTLS11: auditv1.TLSVersion_TLS_VERSION_TLS11,
	tls.VersionTLS12: auditv1.TLSVersion_TLS_VERSION_TLS12,
	tls.VersionTLS13: auditv1.TLSVersion_TLS_VERSION_TLS13,
}

type TLSDetails struct {
	Version     string
	CipherSuite string
	ServerName  string
}

func TLSVersionToEntity(version uint16) auditv1.TLSVersion {
	if v, known := tlsToEntity[version]; known {
		return v
	}
	return auditv1.TLSVersion_TLS_VERSION_UNSPECIFIED
}

func NewTLSDetailsFromProto(entity *auditv1.TLSDetailsEntity) *TLSDetails {
	if entity == nil {
		return nil
	}

	return &TLSDetails{
		Version:     entity.GetVersion().String(),
		CipherSuite: entity.GetCipherSuite(),
		ServerName:  entity.GetServerName(),
	}
}

func (d TLSDetails) ProtoMessage() *auditv1.TLSDetailsEntity {
	return &auditv1.TLSDetailsEntity{
		Version:     auditv1.TLSVersion(auditv1.TLSVersion_value[d.Version]),
		CipherSuite: d.CipherSuite,
		ServerName:  d.ServerName,
	}
}
