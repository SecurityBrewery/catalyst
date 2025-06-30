package pointer

import "testing"

func TestPointer(t *testing.T) {
	t.Parallel()

	v := 42

	ptr := Pointer(v)
	if ptr == nil {
		t.Fatal("Pointer returned nil")
	}

	if *ptr != v {
		t.Errorf("Pointer value = %v, want %v", *ptr, v)
	}
}

func TestDereference(t *testing.T) {
	t.Parallel()

	v := 42
	ptr := &v

	if Dereference(ptr) != v {
		t.Errorf("Dereference(ptr) = %v, want %v", Dereference(ptr), v)
	}

	var nilPtr *int
	if Dereference(nilPtr) != 0 {
		t.Errorf("Dereference(nil) = %v, want 0", Dereference(nilPtr))
	}
}
