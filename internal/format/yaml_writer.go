package format

import (
	"gopkg.in/yaml.v2"
)

type yamlWriter struct {
	encoder *yaml.Encoder
}

func (y *yamlWriter) Write(in interface{}) (err error) {
	return y.encoder.Encode(in)
}
