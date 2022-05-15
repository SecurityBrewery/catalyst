package pointer

import "time"

func String(v string) *string {
	return &v
}

func Int64(v int64) *int64 {
	return &v
}

func Bool(v bool) *bool {
	return &v
}

func Time(v time.Time) *time.Time {
	return &v
}
