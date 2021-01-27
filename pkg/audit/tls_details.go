package audit

import "crypto/tls"

var (
	tlsToEntity = map[uint16]TLSVersion{
		tls.VersionSSL30: TLSVersion_SSLv30,
		tls.VersionTLS10: TLSVersion_TLS10,
		tls.VersionTLS11: TLSVersion_TLS11,
		tls.VersionTLS12: TLSVersion_TLS12,
		tls.VersionTLS13: TLSVersion_TLS13,
	}
)

type TLSDetails struct {
	Version     string
	CipherSuite string
	ServerName  string
}

func TLSVersionToEntity(version uint16) TLSVersion {
	if v, known := tlsToEntity[version]; known {
		return v
	}
	return TLSVersion_SSLv30
}

func NewTLSDetailsFromProto(entity *TLSDetailsEntity) *TLSDetails {
	if entity == nil {
		return nil
	}

	return &TLSDetails{
		Version:     entity.GetVersion().String(),
		CipherSuite: entity.GetCipherSuite(),
		ServerName:  entity.GetServerName(),
	}
}

func (d TLSDetails) ProtoMessage() *TLSDetailsEntity {
	return &TLSDetailsEntity{
		Version:     TLSVersion(TLSVersion_value[d.Version]),
		CipherSuite: d.CipherSuite,
		ServerName:  d.ServerName,
	}
}
