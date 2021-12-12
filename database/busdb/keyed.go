package busdb

import "encoding/json"

type Keyed struct {
	Key string
	Doc interface{}
}

func (p Keyed) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(p.Doc)
	if err != nil {
		panic(err)
	}

	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		panic(err)
	}

	m["_key"] = p.Key

	return json.Marshal(m)
}
