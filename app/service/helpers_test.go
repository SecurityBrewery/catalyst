package service

import (
	"encoding/json"
	"testing"
)

func Test_generateID(t *testing.T) {
	t.Parallel()

	id1 := generateID("p")

	id2 := generateID("p")

	if id1 == id2 {
		t.Errorf("expected unique ids, got %s and %s", id1, id2)
	}

	if len(id1) == 0 || len(id2) == 0 {
		t.Error("expected non empty ids")
	}

	if id1[:2] != "p-" || id2[:2] != "p-" {
		t.Errorf("expected ids to start with prefix")
	}
}

func Test_toString(t *testing.T) {
	t.Parallel()

	str := "hello"

	if toString(&str, "default") != "hello" {
		t.Errorf("unexpected toString result")
	}

	if toString(nil, "default") != "default" {
		t.Errorf("expected default value")
	}
}

func Test_toNullString(t *testing.T) {
	t.Parallel()

	str := "data"

	ns := toNullString(&str)

	if !ns.Valid || ns.String != "data" {
		t.Errorf("unexpected result: %+v", ns)
	}

	ns = toNullString(nil)

	if ns.Valid {
		t.Errorf("expected invalid NullString")
	}
}

func Test_toInt64(t *testing.T) {
	t.Parallel()

	val := 5

	if toInt64(&val, 10) != 5 {
		t.Errorf("unexpected result")
	}

	if toInt64(nil, 10) != 10 {
		t.Errorf("expected default value")
	}
}

func Test_toNullBool(t *testing.T) {
	t.Parallel()

	b := true

	nb := toNullBool(&b)

	if !nb.Valid || !nb.Bool {
		t.Errorf("unexpected result: %+v", nb)
	}

	nb = toNullBool(nil)

	if nb.Valid {
		t.Errorf("expected invalid NullBool")
	}
}

func Test_marshal(t *testing.T) {
	t.Parallel()

	data := map[string]interface{}{"a": 1}

	out := marshal(data)

	var res map[string]interface{}

	err := json.Unmarshal([]byte(out), &res)
	if err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	v, ok := res["a"].(float64)
	if !ok || v != 1 {
		t.Errorf("unexpected marshal result: %v", res)
	}
}

func Test_marshalPointer(t *testing.T) {
	t.Parallel()

	data := map[string]interface{}{"b": 2}

	out := marshalPointer(&data)

	var res map[string]interface{}

	err := json.Unmarshal([]byte(out), &res)
	if err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	v, ok := res["b"].(float64)
	if !ok || v != 2 {
		t.Errorf("unexpected marshalPointer result: %v", res)
	}

	if marshalPointer(nil) != "{}" {
		t.Errorf("expected empty object for nil input")
	}
}

func Test_unmarshal(t *testing.T) {
	t.Parallel()

	jsonStr := `{"c":3}`

	m := unmarshal(jsonStr)

	v, ok := m["c"].(float64)
	if !ok || v != 3 {
		t.Errorf("unexpected unmarshal result: %v", m)
	}

	if unmarshal("invalid") != nil {
		t.Errorf("expected nil for invalid json")
	}
}
