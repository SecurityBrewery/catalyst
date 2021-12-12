// Adapted from https://github.com/golang/go
// Original License:
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the https://go.dev/LICENSE file.

package caql

import "testing"

type MatchTest struct {
	pattern, s string
	match      bool
	err        error
}

var matchTests = []MatchTest{
	{"abc", "abc", true, nil},
	{"%", "abc", true, nil},
	{"%c", "abc", true, nil},
	{"a%", "a", true, nil},
	{"a%", "abc", true, nil},
	{"a%", "ab/c", false, nil},
	{"a%/b", "abc/b", true, nil},
	{"a%/b", "a/c/b", false, nil},
	{"a%b%c%d%e%/f", "axbxcxdxe/f", true, nil},
	{"a%b%c%d%e%/f", "axbxcxdxexxx/f", true, nil},
	{"a%b%c%d%e%/f", "axbxcxdxe/xxx/f", false, nil},
	{"a%b%c%d%e%/f", "axbxcxdxexxx/fff", false, nil},
	{"a%b_c%x", "abxbbxdbxebxczzx", true, nil},
	{"a%b_c%x", "abxbbxdbxebxczzy", false, nil},
	{"a\\%b", "a%b", true, nil},
	{"a\\%b", "ab", false, nil},
	{"a_b", "a☺b", true, nil},
	{"a___b", "a☺b", false, nil},
	{"a_b", "a/b", false, nil},
	{"a%b", "a/b", false, nil},
	{"\\", "a", false, ErrBadPattern},
	{"%x", "xxx", true, nil},
}

func TestMatch(t *testing.T) {
	for _, tt := range matchTests {
		ok, err := match(tt.pattern, tt.s)
		if ok != tt.match || err != tt.err {
			t.Errorf("match(%#q, %#q) = %v, %v want %v, %v", tt.pattern, tt.s, ok, err, tt.match, tt.err)
		}
	}
}
