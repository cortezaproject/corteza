package envoy

import (
	"context"
	"testing"
)

var (
	ops []string
)

type (
	NodeTest struct {
		Id        string
		relations NodeRelation
		exOrder   *[]string
	}
)

func (n *NodeTest) ID() string {
	return n.Id
}

func (n *NodeTest) Resource() string {
	return "test"
}

func (n *NodeTest) Relations() NodeRelation {
	return n.relations
}

func (n *NodeTest) Run(ctx context.Context) error {
	ops = append(ops, "R:"+n.ID())
	return nil
}

func (n *NodeTest) ResolveConflict(ctx context.Context) error {
	ops = append(ops, "C:"+n.ID())
	return nil
}

func TestGraphStructure_relations(t *testing.T) {
	var tests = []struct {
		name     string
		nodes    []*NodeTest
		node     string
		children []string
		parents  []string
	}{
		{
			name: "N1 -> N2",
			nodes: []*NodeTest{
				{Id: "N1", relations: NodeRelation{"test": []string{"N2"}}},
				{Id: "N2", relations: NodeRelation{}},
			},
			node:     "N1",
			children: []string{"N2"},
			parents:  []string{},
		},
		{
			name: "N1 -> N2",
			nodes: []*NodeTest{
				{Id: "N1", relations: NodeRelation{"test": []string{"N2"}}},
				{Id: "N2", relations: NodeRelation{}},
			},
			node:     "N2",
			children: []string{},
			parents:  []string{"N1"},
		},
		{
			name: "N1 -> N1 :: cycle to self",
			nodes: []*NodeTest{
				{Id: "N1", relations: NodeRelation{"test": []string{"N1"}}},
			},
			node:     "N1",
			children: []string{"N1"},
			parents:  []string{"N1"},
		},
		{
			name: "N1 -> N1 -> N2 <- N2",
			nodes: []*NodeTest{
				{Id: "N1", relations: NodeRelation{"test": []string{"N1", "N2"}}},
				{Id: "N2", relations: NodeRelation{"test": []string{"N2"}}},
			},
			node:     "N1",
			children: []string{"N1", "N2"},
			parents:  []string{"N1"},
		},
		{
			name: "N1 -> N1 -> N2 <- N2",
			nodes: []*NodeTest{
				{Id: "N1", relations: NodeRelation{"test": []string{"N1", "N2"}}},
				{Id: "N2", relations: NodeRelation{"test": []string{"N2"}}},
			},
			node:     "N2",
			children: []string{"N2"},
			parents:  []string{"N1", "N2"},
		},
	}

	for _, test := range tests {
		g := Graph{}
		for _, n := range test.nodes {
			g.Add(n)
		}

		cc := g.Children(g.findNode("test", test.node)[0])
		pp := g.Parents(g.findNode("test", test.node)[0])

		if len(cc) != len(test.children) {
			t.Errorf("[%s] node child missmatch; list range doesnt match; exp=%d got=%d", test.name, len(test.children), len(cc))
			return
		}
		for i, c := range cc {
			if c.ID() != test.children[i] {
				t.Errorf("[%s] node child missmatch; exp=%s got=%s pos=%d", test.name, test.children[i], c.ID(), i)
				return
			}
		}

		if len(pp) != len(test.parents) {
			t.Errorf("[%s] node parent missmatch; list range doesnt match; exp=%d got=%d", test.name, len(test.parents), len(pp))
			return
		}
		for i, p := range pp {
			if p.ID() != test.parents[i] {
				t.Errorf("[%s] node parent missmatch; exp=%s got=%s pos=%d", test.name, test.parents[i], p.ID(), i)
				return
			}
		}
	}
}

func TestGraphStructure_inversion(t *testing.T) {
	var tests = []struct {
		name     string
		nodes    []*NodeTest
		node     string
		children []string
		parents  []string
	}{
		{
			name: "N1 -> N2",
			nodes: []*NodeTest{
				{Id: "N1", relations: NodeRelation{"test": []string{"N2"}}},
				{Id: "N2", relations: NodeRelation{}},
			},
			node:     "N1",
			children: []string{},
			parents:  []string{"N2"},
		},
		{
			name: "N1 -> N2",
			nodes: []*NodeTest{
				{Id: "N1", relations: NodeRelation{"test": []string{"N2"}}},
				{Id: "N2", relations: NodeRelation{}},
			},
			node:     "N2",
			children: []string{"N1"},
			parents:  []string{},
		},
		{
			name: "N1 -> N1 :: cycle to self",
			nodes: []*NodeTest{
				{Id: "N1", relations: NodeRelation{"test": []string{"N1"}}},
			},
			node:     "N1",
			children: []string{"N1"},
			parents:  []string{"N1"},
		},
	}

	for _, test := range tests {
		g := Graph{}
		for _, n := range test.nodes {
			g.Add(n)
		}

		g.Invert()

		cc := g.Children(g.findNode("test", test.node)[0])
		pp := g.Parents(g.findNode("test", test.node)[0])

		if len(cc) != len(test.children) {
			t.Errorf("[%s] node child missmatch; list range doesnt match; exp=%d got=%d", test.name, len(test.children), len(cc))
			return
		}
		for i, c := range cc {
			if c.ID() != test.children[i] {
				t.Errorf("[%s] node child missmatch; exp=%s got=%s pos=%d", test.name, test.children[i], c.ID(), i)
				return
			}
		}

		if len(pp) != len(test.parents) {
			t.Errorf("[%s] node parent missmatch; list range doesnt match; exp=%d got=%d", test.name, len(test.parents), len(pp))
			return
		}
		for i, p := range pp {
			if p.ID() != test.parents[i] {
				t.Errorf("[%s] node parent missmatch; exp=%s got=%s pos=%d", test.name, test.parents[i], p.ID(), i)
				return
			}
		}
	}
}

func TestGraphStructure_execution(t *testing.T) {
	var tests = []struct {
		name  string
		nodes []*NodeTest
		ops   []string
	}{
		{
			name: "N1 -> N2",
			nodes: []*NodeTest{
				{Id: "N1", relations: NodeRelation{"test": []string{"N2"}}},
				{Id: "N2", relations: NodeRelation{}},
			},
			ops: []string{"R:N2", "R:N1"},
		},
		{
			name: "N1 -> N1 -> N2",
			nodes: []*NodeTest{
				{Id: "N1", relations: NodeRelation{"test": []string{"N2", "N1"}}},
				{Id: "N2", relations: NodeRelation{}},
			},
			ops: []string{"R:N2", "C:N1", "R:N1"},
		},
		{
			name: "N2 -> N1 -> N1",
			nodes: []*NodeTest{
				{Id: "N1", relations: NodeRelation{"test": []string{"N1"}}},
				{Id: "N2", relations: NodeRelation{"test": []string{"N1"}}},
			},
			ops: []string{"C:N1", "R:N1", "R:N2"},
		},
		{
			name: "N1 -> N1 <-> N2 <-> N3 <- N1",
			nodes: []*NodeTest{
				{Id: "N1", relations: NodeRelation{"test": []string{"N2", "N1", "N3"}}},
				{Id: "N2", relations: NodeRelation{"test": []string{"N1", "N3"}}},
				{Id: "N3", relations: NodeRelation{"test": []string{"N2"}}},
			},
			ops: []string{"C:N1", "C:N2", "R:N3", "R:N1", "R:N2"},
		},
	}

	for _, test := range tests {
		ops = make([]string, 0)
		g := Graph{}
		for _, n := range test.nodes {
			g.Add(n)
		}

		g.Invert()
		g.Run(context.Background())

		for i, o := range test.ops {
			if ops[i] != o {
				t.Errorf("[%s] operation missmatch; exp=%s got=%s", test.name, o, ops[i])
			}
		}

		// if len(cc) != len(test.children) {
		// 	t.Errorf("[%s] node child missmatch; list range doesnt match; exp=%d got=%d", test.name, len(test.children), len(cc))
		// 	return
		// }
		// for i, c := range cc {
		// 	if c.ID() != test.children[i] {
		// 		t.Errorf("[%s] node child missmatch; exp=%s got=%s pos=%d", test.name, test.children[i], c.ID(), i)
		// 		return
		// 	}
		// }

		// if len(pp) != len(test.parents) {
		// 	t.Errorf("[%s] node parent missmatch; list range doesnt match; exp=%d got=%d", test.name, len(test.parents), len(pp))
		// 	return
		// }
		// for i, p := range pp {
		// 	if p.ID() != test.parents[i] {
		// 		t.Errorf("[%s] node parent missmatch; exp=%s got=%s pos=%d", test.name, test.parents[i], p.ID(), i)
		// 		return
		// 	}
		// }
	}
}
