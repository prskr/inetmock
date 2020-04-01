package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var (
	testCaCrt []byte
	testCaKey []byte
)

func init() {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "*-inetmock")
	if err != nil {
		panic(fmt.Sprintf("failed to create temp dir %v", err))
	}

	options := &tlsOptions{
		ecdsaCurve: "P256",
		rootCaCert: cert{
			publicKeyPath:  filepath.Join(tmpDir, "localhost.pem"),
			privateKeyPath: filepath.Join(tmpDir, "localhost.key"),
		},
		validity: validity{
			ca: certValidity{
				notBeforeRelative: time.Hour * 24 * 30,
				notAfterRelative:  time.Hour * 24 * 30,
			},
		},
	}

	certStore := certStore{
		options: options,
		keyProvider: func() (key interface{}, err error) {
			return privateKeyForCurve(options)
		},
	}

	defer func() {
		_ = os.Remove(tmpDir)
	}()

	_, _, err = certStore.generateCaCert()

	testCaCrt, _ = ioutil.ReadFile(options.rootCaCert.publicKeyPath)
	testCaKey, _ = ioutil.ReadFile(options.rootCaCert.privateKeyPath)

	if err != nil {
		panic(fmt.Sprintf("failed to generate CA cert %v", err))
	}
}
