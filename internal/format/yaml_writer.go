package format

import (
	"gopkg.in/yaml.v3"
)

type yamlWriter struct {
	encoder *yaml.Encoder
}

func (y *yamlWriter) Write(in any) (err error) {
	return y.encoder.Encode(in)
}
