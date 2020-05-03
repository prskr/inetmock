package format

import (
	"encoding/json"
)

type jsonWriter struct {
	encoder *json.Encoder
}

func (j *jsonWriter) Write(in interface{}) error {
	return j.encoder.Encode(in)
}
