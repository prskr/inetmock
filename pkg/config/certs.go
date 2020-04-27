package config

import "time"

type CurveType string

type File struct {
	PublicKeyPath  string
	PrivateKeyPath string
}

type ValidityDuration struct {
	NotBeforeRelative time.Duration
	NotAfterRelative  time.Duration
}

type ValidityByPurpose struct {
	CA     ValidityDuration
	Server ValidityDuration
}

type CertOptions struct {
	RootCACert    File
	CertCachePath string
	Curve         CurveType
	Validity      ValidityByPurpose
}
