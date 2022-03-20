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
	"sort"
)

type Set struct {
	hash map[any]nothing
}

type nothing struct{}

func NewSet(initial ...any) *Set {
	s := &Set{make(map[any]nothing)}

	for _, v := range initial {
		s.Insert(v)
	}

	return s
}

func (s *Set) Difference(set *Set) *Set {
	n := make(map[any]nothing)

	for k := range s.hash {
		if _, exists := set.hash[k]; !exists {
			n[k] = nothing{}
		}
	}

	return &Set{n}
}

func (s *Set) Has(element any) bool {
	_, exists := s.hash[element]

	return exists
}

func (s *Set) Insert(element any) {
	s.hash[element] = nothing{}
}

func (s *Set) Intersection(set *Set) *Set {
	n := make(map[any]nothing)

	for k := range s.hash {
		if _, exists := set.hash[k]; exists {
			n[k] = nothing{}
		}
	}

	return &Set{n}
}

func (s *Set) Len() int {
	return len(s.hash)
}

func (s *Set) ProperSubsetOf(set *Set) bool {
	return s.SubsetOf(set) && s.Len() < set.Len()
}

func (s *Set) Remove(element any) {
	delete(s.hash, element)
}

func (s *Set) Minus(set *Set) *Set {
	n := make(map[any]nothing)
	for k := range s.hash {
		n[k] = nothing{}
	}

	for _, v := range set.Values() {
		delete(n, v)
	}

	return &Set{n}
}

func (s *Set) SubsetOf(set *Set) bool {
	if s.Len() > set.Len() {
		return false
	}
	for k := range s.hash {
		if _, exists := set.hash[k]; !exists {
			return false
		}
	}

	return true
}

func (s *Set) Union(set *Set) *Set {
	n := make(map[any]nothing)

	for k := range s.hash {
		n[k] = nothing{}
	}
	for k := range set.hash {
		n[k] = nothing{}
	}

	return &Set{n}
}

func (s *Set) Values() []any {
	values := []any{}

	for k := range s.hash {
		values = append(values, k)
	}
	sort.Slice(values, func(i, j int) bool { return lt(values[i], values[j]) })

	return values
}
