package pointer

func Pointer[T any](v T) *T {
	return &v
}

func Dereference[T any](v *T) T {
	if v == nil {
		var zero T

		return zero
	}

	return *v
}
