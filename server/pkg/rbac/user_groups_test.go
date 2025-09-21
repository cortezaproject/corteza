package rbac

import (
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/stretchr/testify/require"
)

func TestGroupMembership(t *testing.T) {
	svc, err := testPrepState(t)

	require.NoError(t, err)
	require.NotNil(t, svc)

	mm := svc.GroupMembers(id.MustNumID(1))
	require.Len(t, mm, 2)
	require.Contains(t, mm, id.MustNumID(101))
	require.Contains(t, mm, id.MustNumID(102))

	mm = svc.GroupMembers(id.MustNumID(2))
	require.Len(t, mm, 3)
	require.Contains(t, mm, id.MustNumID(121))
	require.Contains(t, mm, id.MustNumID(122))
	require.Contains(t, mm, id.MustNumID(123))
}

func TestMemberBranch(t *testing.T) {
	svc, err := testPrepState(t)
	require.NoError(t, err)
	require.NotNil(t, svc)

	mm, err := svc.MemberBranch(id.MustNumID(102))
	require.NoError(t, err)
	require.Len(t, mm, 4)
	require.Equal(t, id.MustNumID(1), mm[0].id)
	require.Equal(t, id.MustNumID(2), mm[1].id)
	require.Equal(t, id.MustNumID(3), mm[2].id)
	require.Equal(t, id.MustNumID(4), mm[3].id)
}

func TestIsAbove(t *testing.T) {
	svc, err := testPrepState(t)
	require.NoError(t, err)

	t.Run("yes, child", func(t *testing.T) {
		require.True(t, svc.IsAbove(id.MustNumID(121), id.MustNumID(101)))
	})

	t.Run("yes, grandchild", func(t *testing.T) {
		require.True(t, svc.IsAbove(id.MustNumID(141), id.MustNumID(101)))
	})

	t.Run("no, same lvl", func(t *testing.T) {
		require.False(t, svc.IsAbove(id.MustNumID(101), id.MustNumID(102)))
	})

	t.Run("no, below", func(t *testing.T) {
		require.False(t, svc.IsAbove(id.MustNumID(101), id.MustNumID(121)))
	})
}

func TestAddGroupRole(t *testing.T) {
	t.Run("add role", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root", id.MustNumID(0))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", id.MustNumID(101))
		require.NoError(t, err)

		err = svc.AddGroupRole(id.MustNumID(101), id.MustNumID(991))
		require.NoError(t, err)

		_, n := svc.findNode(id.MustNumID(101))
		require.NotNil(t, n)

		require.Len(t, n.roles, 1)
		require.Equal(t, id.MustNumID(991), n.roles[0])
	})

	t.Run("add duplicate role", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root", id.MustNumID(0))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", id.MustNumID(101))
		require.NoError(t, err)

		err = svc.AddGroupRole(id.MustNumID(101), id.MustNumID(991), id.MustNumID(991))
		require.NoError(t, err)

		_, n := svc.findNode(id.MustNumID(101))
		require.NotNil(t, n)

		require.Len(t, n.roles, 1)
		require.Equal(t, id.MustNumID(991), n.roles[0])
	})

	t.Run("append role", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root", id.MustNumID(0))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", id.MustNumID(101))
		require.NoError(t, err)

		err = svc.AddGroupRole(id.MustNumID(101), id.MustNumID(991))
		require.NoError(t, err)

		err = svc.AddGroupRole(id.MustNumID(101), id.MustNumID(992), id.MustNumID(993))
		require.NoError(t, err)

		_, n := svc.findNode(id.MustNumID(101))
		require.NotNil(t, n)

		require.Len(t, n.roles, 3)
		require.Equal(t, id.MustNumID(991), n.roles[0])
		require.Equal(t, id.MustNumID(992), n.roles[1])
		require.Equal(t, id.MustNumID(993), n.roles[2])
	})
}

func TestRemoveGroupRole(t *testing.T) {
	t.Run("remove role", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root", id.MustNumID(0))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", id.MustNumID(101))
		require.NoError(t, err)

		err = svc.AddGroupRole(id.MustNumID(101), id.MustNumID(991))
		require.NoError(t, err)

		err = svc.RemoveGroupRole(id.MustNumID(101), id.MustNumID(991))
		require.NoError(t, err)

		_, n := svc.findNode(id.MustNumID(101))
		require.NotNil(t, n)

		require.Len(t, n.roles, 0)
	})

	t.Run("remove roles", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root", id.MustNumID(0))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", id.MustNumID(101))
		require.NoError(t, err)

		err = svc.AddGroupRole(id.MustNumID(101), id.MustNumID(991), id.MustNumID(992), id.MustNumID(993))
		require.NoError(t, err)

		err = svc.RemoveGroupRole(id.MustNumID(101), id.MustNumID(991), id.MustNumID(993))
		require.NoError(t, err)

		_, n := svc.findNode(id.MustNumID(101))
		require.NotNil(t, n)

		require.Len(t, n.roles, 1)
		require.Equal(t, id.MustNumID(992), n.roles[0])
	})
}

func TestAddNode(t *testing.T) {
	t.Run("empty svc", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(1), "root", id.MustNumID(0))
		require.NoError(t, err)

		ii := svc.root.inline()

		require.Len(t, ii, 1)
		require.Equal(t, id.MustNumID(1), ii[0].id)
	})

	t.Run("add children", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root", id.MustNumID(0))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", id.MustNumID(101))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(202), "c2", id.MustNumID(101))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(301), "c3", id.MustNumID(201))
		require.NoError(t, err)

		// ---

		checkInline(t, svc.root, id.MustNumID(101), id.MustNumID(201), id.MustNumID(202), id.MustNumID(301))

		checkInline(t, svc.branchIndex[id.MustNumID(101)], id.MustNumID(101), id.MustNumID(201), id.MustNumID(202), id.MustNumID(301))
		checkInline(t, svc.branchIndex[id.MustNumID(201)], id.MustNumID(201), id.MustNumID(301))
		checkInline(t, svc.branchIndex[id.MustNumID(202)], id.MustNumID(202))
		checkInline(t, svc.branchIndex[id.MustNumID(301)], id.MustNumID(301))
	})
}

func TestUpdateNode(t *testing.T) {
	t.Run("update without move", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root", id.MustNumID(0))
		require.NoError(t, err)

		err = svc.UpdateNode(id.MustNumID(101), "root edited", id.MustNumID(0))
		require.NoError(t, err)

		require.Equal(t, id.MustNumID(101), svc.root.id)
		require.Equal(t, "root edited", svc.root.handle)
		require.Equal(t, id.MustNumID(0), svc.root.selfID)
	})

	t.Run("change parent", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root", id.MustNumID(0))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", id.MustNumID(101))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(202), "c2", id.MustNumID(101))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(301), "c3", id.MustNumID(201))
		require.NoError(t, err)

		err = svc.UpdateNode(id.MustNumID(301), "c3 edited", id.MustNumID(202))
		require.NoError(t, err)

		checkInline(t, svc.branchIndex[id.MustNumID(201)], id.MustNumID(201))
		checkInline(t, svc.branchIndex[id.MustNumID(202)], id.MustNumID(202), id.MustNumID(301))
	})

	t.Run("update non existing node", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root", id.MustNumID(0))
		require.NoError(t, err)

		err = svc.UpdateNode(id.MustNumID(999), "root edited", id.MustNumID(0))
		require.Error(t, err)
	})
}

func TestRemoveNode(t *testing.T) {
	t.Run("remove leaf node", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root", id.MustNumID(0))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", id.MustNumID(101))
		require.NoError(t, err)

		err = svc.RemoveNode(id.MustNumID(201))
		require.NoError(t, err)

		checkInline(t, svc.branchIndex[id.MustNumID(101)], id.MustNumID(101))
		checkInline(t, svc.branchIndex[id.MustNumID(201)])
	})

	t.Run("can not remove non-lief nodes", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root", id.MustNumID(0))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", id.MustNumID(101))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(301), "c2", id.MustNumID(201))
		require.NoError(t, err)

		err = svc.RemoveNode(id.MustNumID(201))
		require.Error(t, err)
	})

	t.Run("can not remove non-existing nodes", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root", id.MustNumID(0))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", id.MustNumID(101))
		require.NoError(t, err)

		err = svc.RemoveNode(id.MustNumID(999))
		require.Error(t, err)
	})
}

func checkInline(t *testing.T, node *groupNode, ii ...id.ID) {
	nn := node.inline()
	require.Len(t, nn, len(ii))

	for i, n := range nn {
		require.Equal(t, ii[i], n.id)
	}
}

func TestAssignGroupMembers(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		svc, err := mkOrgTree(
			GroupMembers{
				group: &groupNode{
					id:     id.MustNumID(1),
					handle: "root",
				},
				members: nil,
			},
		)
		require.NoError(t, err)
		require.NotNil(t, svc)

		err = svc.AssignGroupMembers(id.MustNumID(1), id.MustNumID(101), id.MustNumID(102))
		require.NoError(t, err)

		members := svc.GroupMembers(id.MustNumID(1))
		require.Len(t, members, 2)

		require.Equal(t, id.MustNumID(101), members[0])
		require.Equal(t, id.MustNumID(102), members[1])
	})

	t.Run("update", func(t *testing.T) {
		svc, err := mkOrgTree(
			GroupMembers{
				group: &groupNode{
					id:     id.MustNumID(1),
					handle: "root",
				},
				members: []id.ID{id.MustNumID(991), id.MustNumID(102)},
			},

			GroupMembers{
				group: &groupNode{
					id:     id.MustNumID(2),
					handle: "subgroup",
					selfID: id.MustNumID(1),
				},
			},
		)
		require.NoError(t, err)
		require.NotNil(t, svc)

		err = svc.AssignGroupMembers(id.MustNumID(2), id.MustNumID(991))
		require.NoError(t, err)

		members := svc.GroupMembers(id.MustNumID(1))
		require.Len(t, members, 1)

		require.Equal(t, id.MustNumID(102), members[0])

		members = svc.GroupMembers(id.MustNumID(2))
		require.Len(t, members, 1)

		require.Equal(t, id.MustNumID(991), members[0])
	})
}

func TestRemoveGroupMembers(t *testing.T) {
	svc, err := mkOrgTree(
		GroupMembers{
			group: &groupNode{
				id:     id.MustNumID(1),
				handle: "root",
			},
			members: []id.ID{id.MustNumID(101), id.MustNumID(102), id.MustNumID(103)},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, svc)

	err = svc.RemoveGroupMembers(id.MustNumID(1), id.MustNumID(102))
	require.NoError(t, err)

	members := svc.GroupMembers(id.MustNumID(1))
	require.Len(t, members, 2)

	require.Equal(t, id.MustNumID(101), members[0])
	require.Equal(t, id.MustNumID(103), members[1])
}

func TestRebuild(t *testing.T) {
	svc, err := testPrepState(t)
	require.NoError(t, err)
	require.NotNil(t, svc)

	err = svc.Rebuild(GroupMembers{
		group: &groupNode{
			id:     id.MustNumID(999),
			handle: "newRoot",
		},
		members: []id.ID{id.MustNumID(991), id.MustNumID(992), id.MustNumID(993)},
	})

	require.NoError(t, err)

	members := svc.GroupMembers(id.MustNumID(1))
	require.Len(t, members, 0)

	members = svc.GroupMembers(id.MustNumID(999))
	require.Len(t, members, 3)
}

func TestBuildOrgTree(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		out, _, err := buildOrgTree()
		require.NoError(t, err)
		require.Nil(t, out)
	})

	t.Run("just root", func(t *testing.T) {
		out, _, err := buildOrgTree(&groupNode{id: id.MustNumID(1), handle: "root", selfID: id.MustNumID(0)})
		require.NoError(t, err)

		nodes := out.inline()
		require.Len(t, nodes, 1)
		require.Equal(t, "root", nodes[0].handle)
	})

	t.Run("multiple root error", func(t *testing.T) {
		_, _, err := buildOrgTree(
			&groupNode{id: id.MustNumID(1), handle: "root1", selfID: id.MustNumID(0)},
			&groupNode{id: id.MustNumID(2), handle: "root2", selfID: id.MustNumID(0)},
		)
		require.Error(t, err)
		require.EqualError(t, err, "multiple root nodes detected")
	})

	t.Run("missing parent error", func(t *testing.T) {
		_, _, err := buildOrgTree(
			&groupNode{id: id.MustNumID(1), handle: "root1", selfID: id.MustNumID(0)},
			&groupNode{id: id.MustNumID(2), handle: "child1", selfID: id.MustNumID(99)},
		)
		require.Error(t, err)
		require.EqualError(t, err, `node "2" parent "99" not found`)
	})

	t.Run("build tree", func(t *testing.T) {
		out, _, err := buildOrgTree(
			&groupNode{id: id.MustNumID(1), handle: "root", selfID: id.MustNumID(0)},
			&groupNode{id: id.MustNumID(2), handle: "child1", selfID: id.MustNumID(1)},
			&groupNode{id: id.MustNumID(3), handle: "child2", selfID: id.MustNumID(1)},
			&groupNode{id: id.MustNumID(4), handle: "child1_grandchild1", selfID: id.MustNumID(2)},
			&groupNode{id: id.MustNumID(5), handle: "child1_grandchild2", selfID: id.MustNumID(2)},
			&groupNode{id: id.MustNumID(6), handle: "child2_grandchild1", selfID: id.MustNumID(3)},
			&groupNode{id: id.MustNumID(7), handle: "child1_grandchild2_ggrandchild1", selfID: id.MustNumID(5)},
		)
		require.NoError(t, err)

		nn := out.inline()
		require.Len(t, nn, 7)
		require.Equal(t, "root", nn[0].handle)
		require.Equal(t, "child1", nn[1].handle)
		require.Equal(t, "child2", nn[2].handle)
		require.Equal(t, "child1_grandchild1", nn[3].handle)
		require.Equal(t, "child1_grandchild2", nn[4].handle)
		require.Equal(t, "child2_grandchild1", nn[5].handle)
		require.Equal(t, "child1_grandchild2_ggrandchild1", nn[6].handle)
	})
}

func TestFindNode(t *testing.T) {
	svc, err := testPrepState(t)
	require.NoError(t, err)

	t.Run("first", func(t *testing.T) {
		i, n := svc.findNode(id.MustNumID(1))
		require.Equal(t, 0, i)
		require.NotNil(t, n)
	})

	t.Run("last", func(t *testing.T) {
		i, n := svc.findNode(id.MustNumID(4))
		require.Equal(t, 3, i)
		require.NotNil(t, n)
	})

	t.Run("middle", func(t *testing.T) {
		i, n := svc.findNode(id.MustNumID(3))
		require.Equal(t, 2, i)
		require.NotNil(t, n)
	})

	t.Run("not existing", func(t *testing.T) {
		i, n := svc.findNode(id.MustNumID(999))
		require.Equal(t, -1, i)
		require.Nil(t, n)
	})
}

func testPrepState(t *testing.T) (svc *orgTree, err error) {
	return mkOrgTree(GroupMembers{
		group: &groupNode{
			id:     id.MustNumID(1),
			handle: "root",
		},
		members: []id.ID{id.MustNumID(101), id.MustNumID(102)},
	}, GroupMembers{
		group: &groupNode{
			id:     id.MustNumID(2),
			handle: "subgroup1",
			selfID: id.MustNumID(1),
		},
		members: []id.ID{id.MustNumID(121), id.MustNumID(122), id.MustNumID(123)},
	}, GroupMembers{
		group: &groupNode{
			id:     id.MustNumID(3),
			handle: "subgroup2",
			selfID: id.MustNumID(1),
		},
		members: []id.ID{id.MustNumID(131)},
	}, GroupMembers{
		group: &groupNode{
			id:     id.MustNumID(4),
			handle: "subgroup1_subgroup1",
			selfID: id.MustNumID(2),
		},
		members: []id.ID{id.MustNumID(141), id.MustNumID(142)},
	})
}
