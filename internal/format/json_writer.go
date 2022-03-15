package format

import (
	"encoding/json"
)

type jsonWriter struct {
	encoder *json.Encoder
}

func (j *jsonWriter) Write(in any) error {
	return j.encoder.Encode(in)
}
