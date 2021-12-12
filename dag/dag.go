// Adapted from https://github.com/philopon/go-toposort under the MIT License
// Original License:
//
// Copyright (c) 2017 Hirotomo Moriwaki
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

package dag

import (
	"errors"
	"sort"
)

type Graph struct {
	nodes []string

	outputs map[string]map[string]struct{}

	// node: number of parents
	inputs map[string]int
}

func NewGraph() *Graph {
	return &Graph{
		nodes:   []string{},
		inputs:  make(map[string]int),
		outputs: make(map[string]map[string]struct{}),
	}
}

func (g *Graph) AddNode(name string) error {
	g.nodes = append(g.nodes, name)

	if _, ok := g.outputs[name]; ok {
		return errors.New("duplicate detected")
	}
	g.outputs[name] = make(map[string]struct{})
	g.inputs[name] = 0
	return nil
}

func (g *Graph) AddNodes(names ...string) error {
	for _, name := range names {
		if err := g.AddNode(name); err != nil {
			return err
		}
	}
	return nil
}

func (g *Graph) AddEdge(from, to string) error {
	m, ok := g.outputs[from]
	if !ok {
		return errors.New("node does not exist")
	}

	m[to] = struct{}{}
	g.inputs[to]++

	return nil
}

func (g *Graph) Toposort() ([]string, error) {
	outputs := map[string]map[string]struct{}{}
	for key, value := range g.outputs {
		outputs[key] = map[string]struct{}{}
		for k, v := range value {
			outputs[key][k] = v
		}
	}

	L := make([]string, 0, len(g.nodes))
	S := make([]string, 0, len(g.nodes))

	sort.Strings(g.nodes)
	for _, n := range g.nodes {
		if g.inputs[n] == 0 {
			S = append(S, n)
		}
	}

	for len(S) > 0 {
		var n string
		n, S = S[0], S[1:]
		L = append(L, n)

		ms := make([]string, len(outputs[n]))
		for _, k := range keys(outputs[n]) {
			m := k
			// i := outputs[n][m]
			// ms[i-1] = m
			ms = append(ms, m)
		}

		for _, m := range ms {
			delete(outputs[n], m)
			g.inputs[m]--

			if g.inputs[m] == 0 {
				S = append(S, m)
			}
		}
	}

	N := 0
	for _, v := range g.inputs {
		N += v
	}

	if N > 0 {
		return L, errors.New("cycle detected")
	}

	return L, nil
}

func keys(m map[string]struct{}) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (g *Graph) GetParents(id string) []string {
	var parents []string
	for node, targets := range g.outputs {
		if _, ok := targets[id]; ok {
			parents = append(parents, node)
		}
	}
	sort.Strings(parents)
	return parents
}

func (g *Graph) GetRoot() (string, error) {
	var roots []string
	for n, parents := range g.inputs {
		if parents == 0 {
			roots = append(roots, n)
		}
	}
	if len(roots) != 1 {
		return "", errors.New("more than one root")
	}
	return roots[0], nil
}
