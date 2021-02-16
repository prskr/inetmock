package app

import (
	"os"
	"path/filepath"
)

type Data struct {
	PCAP  string
	Audit string
}

func (d *Data) setup() (err error) {
	if d.PCAP, err = ensureDataDir(d.PCAP); err != nil {
		return
	}
	if d.Audit, err = ensureDataDir(d.Audit); err != nil {
		return
	}

	return
}

func ensureDataDir(dataDirPath string) (cleanedPath string, err error) {
	cleanedPath = dataDirPath
	if !filepath.IsAbs(cleanedPath) {
		if cleanedPath, err = filepath.Abs(cleanedPath); err != nil {
			return
		}
	}

	err = os.MkdirAll(cleanedPath, 0640)
	return
}
