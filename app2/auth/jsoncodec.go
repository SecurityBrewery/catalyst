package auth

import (
	"encoding/json"
	"time"
)

type codecData struct {
	Deadline time.Time      `json:"deadline"`
	Values   map[string]any `json:"values"`
}

type JSONCodec struct{}

func (J *JSONCodec) Encode(deadline time.Time, values map[string]any) ([]byte, error) {
	return json.Marshal(codecData{
		Deadline: deadline,
		Values:   values,
	})
}

func (J *JSONCodec) Decode(bytes []byte) (deadline time.Time, values map[string]any, err error) {
	var data codecData
	if err = json.Unmarshal(bytes, &data); err != nil {
		return time.Time{}, nil, err
	}

	return data.Deadline, data.Values, nil
}
