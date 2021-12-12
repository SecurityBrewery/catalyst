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

import "sort"

type (
	Set struct {
		hash map[interface{}]nothing
	}

	nothing struct{}
)

// Create a new set
func New(initial ...interface{}) *Set {
	s := &Set{make(map[interface{}]nothing)}

	for _, v := range initial {
		s.Insert(v)
	}

	return s
}

// Find the difference between two sets
func (s *Set) Difference(set *Set) *Set {
	n := make(map[interface{}]nothing)

	for k := range s.hash {
		if _, exists := set.hash[k]; !exists {
			n[k] = nothing{}
		}
	}

	return &Set{n}
}

// Call f for each item in the set
func (s *Set) Do(f func(interface{})) {
	for k := range s.hash {
		f(k)
	}
}

// Test to see whether or not the element is in the set
func (s *Set) Has(element interface{}) bool {
	_, exists := s.hash[element]
	return exists
}

// Add an element to the set
func (s *Set) Insert(element interface{}) {
	s.hash[element] = nothing{}
}

// Find the intersection of two sets
func (s *Set) Intersection(set *Set) *Set {
	n := make(map[interface{}]nothing)

	for k := range s.hash {
		if _, exists := set.hash[k]; exists {
			n[k] = nothing{}
		}
	}

	return &Set{n}
}

// Return the number of items in the set
func (s *Set) Len() int {
	return len(s.hash)
}

// Test whether or not this set is a proper subset of "set"
func (s *Set) ProperSubsetOf(set *Set) bool {
	return s.SubsetOf(set) && s.Len() < set.Len()
}

// Remove an element from the set
func (s *Set) Remove(element interface{}) {
	delete(s.hash, element)
}

func (s *Set) Minus(set *Set) *Set {
	n := make(map[interface{}]nothing)
	for k := range s.hash {
		n[k] = nothing{}
	}

	for _, v := range set.Values() {
		delete(n, v)
	}

	return &Set{n}
}

// Test whether or not this set is a subset of "set"
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

// Find the union of two sets
func (s *Set) Union(set *Set) *Set {
	n := make(map[interface{}]nothing)

	for k := range s.hash {
		n[k] = nothing{}
	}
	for k := range set.hash {
		n[k] = nothing{}
	}

	return &Set{n}
}

func (s *Set) Values() []interface{} {
	values := []interface{}{}

	for k := range s.hash {
		values = append(values, k)
	}
	sort.Slice(values, func(i, j int) bool { return lt(values[i], values[j]) })

	return values
}
