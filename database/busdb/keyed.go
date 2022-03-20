package busdb

import "encoding/json"

type Keyed[T any] struct {
	Key string
	Doc *T
}

func (p *Keyed[T]) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(p.Doc)
	if err != nil {
		panic(err)
	}

	var m map[string]any
	err = json.Unmarshal(b, &m)
	if err != nil {
		panic(err)
	}

	m["_key"] = p.Key

	return json.Marshal(m)
}
