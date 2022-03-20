// Adapted from https://github.com/badgerodon/collections under the MIT License
// Original License:
//
// Copyright (c) 2012 Caleb Doxsey
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package caql

import (
	"testing"
)

func Test(t *testing.T) {
	t.Parallel()

	s := NewSet()

	s.Insert(5)

	if s.Len() != 1 {
		t.Errorf("Length should be 1")
	}

	if !s.Has(5) {
		t.Errorf("Membership test failed")
	}

	s.Remove(5)

	if s.Len() != 0 {
		t.Errorf("Length should be 0")
	}

	if s.Has(5) {
		t.Errorf("The set should be empty")
	}

	// Difference
	s1 := NewSet(1, 2, 3, 4, 5, 6)
	s2 := NewSet(4, 5, 6)
	s3 := s1.Difference(s2)

	if s3.Len() != 3 {
		t.Errorf("Length should be 3")
	}

	if !(s3.Has(1) && s3.Has(2) && s3.Has(3)) {
		t.Errorf("Set should only contain 1, 2, 3")
	}

	// Intersection
	s3 = s1.Intersection(s2)
	if s3.Len() != 3 {
		t.Errorf("Length should be 3 after intersection")
	}

	if !(s3.Has(4) && s3.Has(5) && s3.Has(6)) {
		t.Errorf("Set should contain 4, 5, 6")
	}

	// Union
	s4 := NewSet(7, 8, 9)
	s3 = s2.Union(s4)

	if s3.Len() != 6 {
		t.Errorf("Length should be 6 after union")
	}

	if !(s3.Has(7)) {
		t.Errorf("Set should contain 4, 5, 6, 7, 8, 9")
	}

	// Subset
	if !s1.SubsetOf(s1) {
		t.Errorf("set should be a subset of itself")
	}
	// Proper Subset
	if s1.ProperSubsetOf(s1) {
		t.Errorf("set should not be a subset of itself")
	}
}
