package service

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_marshal(t *testing.T) {
	t.Parallel()

	data := map[string]any{"a": 1}

	out := marshal(data)

	var res map[string]any
	err := json.Unmarshal([]byte(out), &res)
	require.NoError(t, err, "invalid json")

	v, ok := res["a"].(float64)
	assert.True(t, ok, "key 'a' not found or not float64")
	assert.InEpsilon(t, float64(1), v, 0, "unexpected marshal result")
}

func Test_unmarshal(t *testing.T) {
	t.Parallel()

	jsonStr := `{"c":3}`

	m := unmarshal(jsonStr)

	v, ok := m["c"].(float64)
	assert.True(t, ok, "key 'c' not found or not float64")
	assert.InEpsilon(t, float64(3), v, 0, "unexpected unmarshal result")

	assert.Nil(t, unmarshal("invalid"), "expected nil for invalid json")
}
