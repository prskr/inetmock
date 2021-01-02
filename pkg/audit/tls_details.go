package audit

import (
	"crypto/tls"
)

var (
	tlsToEntity = map[uint16]TLSVersion{
		tls.VersionSSL30: TLSVersion_SSLv30,
		tls.VersionTLS10: TLSVersion_TLS10,
		tls.VersionTLS11: TLSVersion_TLS11,
		tls.VersionTLS12: TLSVersion_TLS12,
		tls.VersionTLS13: TLSVersion_TLS13,
	}
	entityToTls = map[TLSVersion]uint16{
		TLSVersion_SSLv30: tls.VersionSSL30,
		TLSVersion_TLS10:  tls.VersionTLS10,
		TLSVersion_TLS11:  tls.VersionTLS11,
		TLSVersion_TLS12:  tls.VersionTLS12,
		TLSVersion_TLS13:  tls.VersionTLS13,
	}
	cipherSuiteIDLookup = func(name string) uint16 {
		for _, cs := range tls.CipherSuites() {
			if cs.Name == name {
				return cs.ID
			}
		}
		return 0
	}
)

type TLSDetails struct {
	Version     uint16
	CipherSuite uint16
	ServerName  string
}

func NewTLSDetailsFromProto(entity *TLSDetailsEntity) *TLSDetails {
	if entity == nil {
		return nil
	}

	return &TLSDetails{
		Version:     entityToTls[entity.GetVersion()],
		CipherSuite: cipherSuiteIDLookup(entity.GetCipherSuite()),
		ServerName:  entity.GetServerName(),
	}
}

func (d TLSDetails) ProtoMessage() *TLSDetailsEntity {
	return &TLSDetailsEntity{
		Version:     tlsToEntity[d.Version],
		CipherSuite: tls.CipherSuiteName(d.CipherSuite),
		ServerName:  d.ServerName,
	}
}
