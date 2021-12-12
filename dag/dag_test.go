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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func index(s []string, v string) int {
	for i, s := range s {
		if s == v {
			return i
		}
	}
	return -1
}

type Edge struct {
	From string
	To   string
}

func TestDuplicatedNode(t *testing.T) {
	graph := NewGraph()
	assert.NoError(t, graph.AddNode("a"))
	assert.Error(t, graph.AddNode("a"))
}

func TestWikipedia(t *testing.T) {
	graph := NewGraph()
	assert.NoError(t, graph.AddNodes("2", "3", "5", "7", "8", "9", "10", "11"))

	edges := []Edge{
		{"7", "8"},
		{"7", "11"},

		{"5", "11"},

		{"3", "8"},
		{"3", "10"},

		{"11", "2"},
		{"11", "9"},
		{"11", "10"},

		{"8", "9"},
	}

	for _, e := range edges {
		assert.NoError(t, graph.AddEdge(e.From, e.To))
	}

	result, err := graph.Toposort()
	if err != nil {
		t.Errorf("closed path detected in no closed pathed graph")
	}

	for _, e := range edges {
		if i, j := index(result, e.From), index(result, e.To); i > j {
			t.Errorf("dependency failed: not satisfy %v(%v) > %v(%v)", e.From, i, e.To, j)
		}
	}
}

func TestCycle(t *testing.T) {
	graph := NewGraph()
	assert.NoError(t, graph.AddNodes("1", "2", "3"))

	assert.NoError(t, graph.AddEdge("1", "2"))
	assert.NoError(t, graph.AddEdge("2", "3"))
	assert.NoError(t, graph.AddEdge("3", "1"))

	_, err := graph.Toposort()
	if err == nil {
		t.Errorf("closed path not detected in closed pathed graph")
	}
}

func TestGraph_GetParents(t *testing.T) {
	type fields struct {
		nodes []string
		edges map[string]string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{"parents 2", fields{nodes: []string{"1", "2", "3"}, edges: map[string]string{"1": "2", "2": "3"}}, args{id: "2"}, []string{"1"}},
		{"parents 3", fields{nodes: []string{"1", "2", "3"}, edges: map[string]string{"1": "3", "2": "3"}}, args{id: "3"}, []string{"1", "2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGraph()
			for _, node := range tt.fields.nodes {
				assert.NoError(t, g.AddNode(node))
			}
			for from, to := range tt.fields.edges {
				assert.NoError(t, g.AddEdge(from, to))
			}

			if got := g.GetParents(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetParents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDAG_AddNode(t *testing.T) {
	dag := NewGraph()

	v := "1"
	assert.NoError(t, dag.AddNode(v))

	assert.Error(t, dag.AddNode(v))
}

func TestDAG_AddEdge(t *testing.T) {
	dag := NewGraph()
	assert.NoError(t, dag.AddNode("0"))
	assert.NoError(t, dag.AddNode("1"))
	assert.NoError(t, dag.AddNode("2"))
	assert.NoError(t, dag.AddNode("3"))

	// add a single edge and inspect the graph
	assert.NoError(t, dag.AddEdge("1", "2"))

	if parents := dag.GetParents("2"); len(parents) != 1 {
		t.Errorf("GetParents(v2) = %d, want 1", len(parents))
	}

	assert.NoError(t, dag.AddEdge("2", "3"))

	_ = dag.AddEdge("0", "1")
}

func TestDAG_GetParents(t *testing.T) {
	dag := NewGraph()
	assert.NoError(t, dag.AddNode("1"))
	assert.NoError(t, dag.AddNode("2"))
	assert.NoError(t, dag.AddNode("3"))
	_ = dag.AddEdge("1", "3")
	_ = dag.AddEdge("2", "3")

	parents := dag.GetParents("3")
	if length := len(parents); length != 2 {
		t.Errorf("GetParents(v3) = %d, want 2", length)
	}
}

func TestDAG_GetDescendants(t *testing.T) {
	dag := NewGraph()
	assert.NoError(t, dag.AddNode("1"))
	assert.NoError(t, dag.AddNode("2"))
	assert.NoError(t, dag.AddNode("3"))
	assert.NoError(t, dag.AddNode("4"))

	assert.NoError(t, dag.AddEdge("1", "2"))
	assert.NoError(t, dag.AddEdge("2", "3"))
	assert.NoError(t, dag.AddEdge("2", "4"))
}

func TestDAG_Topsort(t *testing.T) {
	dag := NewGraph()
	assert.NoError(t, dag.AddNode("1"))
	assert.NoError(t, dag.AddNode("2"))
	assert.NoError(t, dag.AddNode("3"))
	assert.NoError(t, dag.AddNode("4"))

	assert.NoError(t, dag.AddEdge("1", "2"))
	assert.NoError(t, dag.AddEdge("2", "3"))
	assert.NoError(t, dag.AddEdge("2", "4"))

	desc, _ := dag.Toposort()
	assert.Equal(t, desc, []string{"1", "2", "3", "4"})
}

func TestDAG_TopsortStable(t *testing.T) {
	dag := NewGraph()
	assert.NoError(t, dag.AddNode("1"))
	assert.NoError(t, dag.AddNode("2"))
	assert.NoError(t, dag.AddNode("3"))

	assert.NoError(t, dag.AddEdge("1", "2"))
	assert.NoError(t, dag.AddEdge("1", "3"))

	desc, _ := dag.Toposort()
	assert.Equal(t, desc, []string{"1", "2", "3"})
}

func TestDAG_TopsortStable2(t *testing.T) {
	dag := NewGraph()

	assert.NoError(t, dag.AddNodes("block-ioc", "block-iocs", "block-sender", "board", "fetch-iocs", "escalate", "extract-iocs", "mail-available", "search-email-gateway"))
	assert.NoError(t, dag.AddEdge("block-iocs", "block-ioc"))
	assert.NoError(t, dag.AddEdge("block-sender", "extract-iocs"))
	assert.NoError(t, dag.AddEdge("board", "escalate"))
	assert.NoError(t, dag.AddEdge("board", "mail-available"))
	assert.NoError(t, dag.AddEdge("fetch-iocs", "block-iocs"))
	assert.NoError(t, dag.AddEdge("extract-iocs", "fetch-iocs"))
	assert.NoError(t, dag.AddEdge("mail-available", "block-sender"))
	assert.NoError(t, dag.AddEdge("mail-available", "extract-iocs"))
	assert.NoError(t, dag.AddEdge("mail-available", "search-email-gateway"))
	assert.NoError(t, dag.AddEdge("search-email-gateway", "extract-iocs"))

	sorted, err := dag.Toposort()
	assert.NoError(t, err)

	want := []string{"board", "escalate", "mail-available", "block-sender", "search-email-gateway", "extract-iocs", "fetch-iocs", "block-iocs", "block-ioc"}
	assert.Equal(t, want, sorted)
}
