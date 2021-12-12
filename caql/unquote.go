// Adapted from https://github.com/golang/go
// Original License:
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the https://go.dev/LICENSE file.

package caql

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

// unquote interprets s as a single-quoted, double-quoted,
// or backquoted string literal, returning the string value
// that s quotes.
func unquote(s string) (string, error) {
	n := len(s)
	if n < 2 {
		return "", strconv.ErrSyntax
	}
	quote := s[0]
	if quote != s[n-1] {
		return "", strconv.ErrSyntax
	}
	s = s[1 : n-1]

	if quote == '`' {
		if strings.ContainsRune(s, '`') {
			return "", strconv.ErrSyntax
		}
		if strings.ContainsRune(s, '\r') {
			// -1 because we know there is at least one \r to remove.
			buf := make([]byte, 0, len(s)-1)
			for i := 0; i < len(s); i++ {
				if s[i] != '\r' {
					buf = append(buf, s[i])
				}
			}
			return string(buf), nil
		}
		return s, nil
	}
	if quote != '"' && quote != '\'' {
		return "", strconv.ErrSyntax
	}
	if strings.ContainsRune(s, '\n') {
		return "", strconv.ErrSyntax
	}

	// Is it trivial? Avoid allocation.
	if !strings.ContainsRune(s, '\\') && !strings.ContainsRune(s, rune(quote)) {
		switch quote {
		case '"', '\'':
			if utf8.ValidString(s) {
				return s, nil
			}
		}
	}

	var runeTmp [utf8.UTFMax]byte
	buf := make([]byte, 0, 3*len(s)/2) // Try to avoid more allocations.
	for len(s) > 0 {
		c, multibyte, ss, err := strconv.UnquoteChar(s, quote)
		if err != nil {
			return "", err
		}
		s = ss
		if c < utf8.RuneSelf || !multibyte {
			buf = append(buf, byte(c))
		} else {
			n := utf8.EncodeRune(runeTmp[:], c)
			buf = append(buf, runeTmp[:n]...)
		}
	}
	return string(buf), nil
}
