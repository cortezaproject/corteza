package rbac

import (
	"testing"

	"github.com/crusttech/human/server/pkg/id"
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
	// @todo add some test cases for named paths; they use the same underlying fnc so I'm skipping for now
	t.Run("unnamed paths", func(t *testing.T) {
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
	})
}

func TestIsAbove(t *testing.T) {
	t.Run("unnamed", func(t *testing.T) {
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
	})

	t.Run("named", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root")
		require.NoError(t, err)
		err = svc.AssignGroupMembers(id.MustNumID(101), id.MustNumID(101))
		require.NoError(t, err)

		//

		err = svc.AddNode(id.MustNumID(201), "c1", mkPpn(id.MustNumID(101), "a")...)
		require.NoError(t, err)
		err = svc.AssignGroupMembers(id.MustNumID(201), id.MustNumID(201))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(202), "c2", mkPpn(id.MustNumID(101), "a", id.MustNumID(101), "b")...)
		require.NoError(t, err)
		err = svc.AssignGroupMembers(id.MustNumID(202), id.MustNumID(202))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(203), "c3", mkPp(id.MustNumID(101))...)
		require.NoError(t, err)
		err = svc.AssignGroupMembers(id.MustNumID(203), id.MustNumID(203))
		require.NoError(t, err)

		//

		err = svc.AddNode(id.MustNumID(301), "cc1", mkPpn(id.MustNumID(201), "a")...)
		require.NoError(t, err)
		err = svc.AssignGroupMembers(id.MustNumID(301), id.MustNumID(301))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(302), "cc2", mkPpn(id.MustNumID(202), "a")...)
		require.NoError(t, err)
		err = svc.AssignGroupMembers(id.MustNumID(302), id.MustNumID(302))
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(303), "cc3", mkPpn(id.MustNumID(101), "a")...)
		require.NoError(t, err)
		err = svc.AssignGroupMembers(id.MustNumID(303), id.MustNumID(303))
		require.NoError(t, err)

		// ---

		t.Run("yes, child, labeled path", func(t *testing.T) {
			require.True(t, svc.IsAbove(id.MustNumID(201), id.MustNumID(101), "a"))
		})

		t.Run("yes, child, has labeled path but looking generic", func(t *testing.T) {
			require.True(t, svc.IsAbove(id.MustNumID(201), id.MustNumID(101)))
		})

		t.Run("no, child, wrong labeled path", func(t *testing.T) {
			require.False(t, svc.IsAbove(id.MustNumID(201), id.MustNumID(101), "b"))
		})

		t.Run("yes, grandchild, labeled path", func(t *testing.T) {
			require.True(t, svc.IsAbove(id.MustNumID(302), id.MustNumID(101), "a"))
		})

	})
}

func TestAddGroupRole(t *testing.T) {
	t.Run("add role", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root")
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", mkPp(id.MustNumID(101))...)
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

		err = svc.AddNode(id.MustNumID(101), "root")
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", mkPp(id.MustNumID(101))...)
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

		err = svc.AddNode(id.MustNumID(101), "root")
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", mkPp(id.MustNumID(101))...)
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

		err = svc.AddNode(id.MustNumID(101), "root")
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", mkPp(id.MustNumID(101))...)
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

		err = svc.AddNode(id.MustNumID(101), "root")
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", mkPp(id.MustNumID(101))...)
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

		err = svc.AddNode(id.MustNumID(1), "root")
		require.NoError(t, err)

		ii := svc.root.inline()

		require.Len(t, ii, 1)
		require.Equal(t, id.MustNumID(1), ii[0].id)
	})

	t.Run("add children", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root")
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", mkPp(id.MustNumID(101))...)
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(202), "c2", mkPp(id.MustNumID(101))...)
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(301), "c3", mkPp(id.MustNumID(201))...)
		require.NoError(t, err)

		// ---

		checkInline(t, svc.root, id.MustNumID(101), id.MustNumID(201), id.MustNumID(202), id.MustNumID(301))

		checkInline(t, svc.branchIndex[id.MustNumID(101)], id.MustNumID(101), id.MustNumID(201), id.MustNumID(202), id.MustNumID(301))
		checkInline(t, svc.branchIndex[id.MustNumID(201)], id.MustNumID(201), id.MustNumID(301))
		checkInline(t, svc.branchIndex[id.MustNumID(202)], id.MustNumID(202))
		checkInline(t, svc.branchIndex[id.MustNumID(301)], id.MustNumID(301))
	})

	t.Run("check path management", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		// Add initial stuff
		err = svc.AddNode(id.MustNumID(101), "root")
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", GroupNodePath{SelfID: id.MustNumID(101)})
		require.NoError(t, err)
		err = svc.AddNode(id.MustNumID(202), "c2", GroupNodePath{SelfID: id.MustNumID(101)})
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(301), "cc1", GroupNodePath{SelfID: id.MustNumID(201)})
		require.NoError(t, err)

		require.Len(t, svc.branchIndex[id.MustNumID(101)].parents, 0)
		require.Len(t, svc.branchIndex[id.MustNumID(101)].children, 2)
		require.Equal(t, id.MustNumID(201), svc.branchIndex[id.MustNumID(101)].children[0].node.id)
		require.Equal(t, id.MustNumID(202), svc.branchIndex[id.MustNumID(101)].children[1].node.id)

		require.Len(t, svc.branchIndex[id.MustNumID(201)].parents, 1)
		require.Len(t, svc.branchIndex[id.MustNumID(201)].children, 1)
		require.Equal(t, id.MustNumID(301), svc.branchIndex[id.MustNumID(201)].children[0].node.id)
		require.Equal(t, id.MustNumID(101), svc.branchIndex[id.MustNumID(201)].parents[0].node.id)

		// Change c2 connection

		err = svc.UpdateNode(id.MustNumID(202), "ccc2 edit", GroupNodePath{SelfID: id.MustNumID(301)})
		require.NoError(t, err)

		require.Len(t, svc.branchIndex[id.MustNumID(101)].parents, 0)
		require.Len(t, svc.branchIndex[id.MustNumID(101)].children, 1)
		require.Equal(t, id.MustNumID(201), svc.branchIndex[id.MustNumID(101)].children[0].node.id)

		require.Len(t, svc.branchIndex[id.MustNumID(301)].parents, 1)
		require.Len(t, svc.branchIndex[id.MustNumID(301)].children, 1)
		require.Equal(t, id.MustNumID(201), svc.branchIndex[id.MustNumID(301)].parents[0].node.id)
		require.Equal(t, id.MustNumID(202), svc.branchIndex[id.MustNumID(301)].children[0].node.id)

		require.Len(t, svc.branchIndex[id.MustNumID(202)].parents, 1)
		require.Len(t, svc.branchIndex[id.MustNumID(202)].children, 0)
		require.Equal(t, id.MustNumID(301), svc.branchIndex[id.MustNumID(202)].parents[0].node.id)

		// Remove c2 node

		err = svc.RemoveNode(id.MustNumID(202))
		require.NoError(t, err)

		require.Len(t, svc.branchIndex[id.MustNumID(301)].parents, 1)
		require.Len(t, svc.branchIndex[id.MustNumID(301)].children, 0)
		require.Equal(t, id.MustNumID(201), svc.branchIndex[id.MustNumID(301)].parents[0].node.id)
	})
}

func TestMultiPath(t *testing.T) {
	svc, err := mkOrgTree()
	require.NoError(t, err)

	err = svc.AddNode(id.MustNumID(101), "root")
	require.NoError(t, err)

	// lvl 2 roots

	err = svc.AddNode(id.MustNumID(201), "c1", mkPp(id.MustNumID(101))...)
	require.NoError(t, err)

	err = svc.AddNode(id.MustNumID(202), "c2", mkPp(id.MustNumID(101))...)
	require.NoError(t, err)

	// Multi path child
	err = svc.AddNode(id.MustNumID(301), "cc1", mkPp(id.MustNumID(201), id.MustNumID(202))...)
	require.NoError(t, err)

	// Check

	gg := svc.branchIndex[id.MustNumID(201)].inline()
	require.Len(t, gg, 2)

	require.Equal(t, id.MustNumID(201), gg[0].id)
	require.Equal(t, id.MustNumID(301), gg[1].id)

	gg = svc.branchIndex[id.MustNumID(202)].inline()
	require.Len(t, gg, 2)

	require.Equal(t, id.MustNumID(202), gg[0].id)
	require.Equal(t, id.MustNumID(301), gg[1].id)
}

func TestNamedPaths(t *testing.T) {
	svc, err := mkOrgTree()
	require.NoError(t, err)

	err = svc.AddNode(id.MustNumID(101), "root")
	require.NoError(t, err)

	//

	err = svc.AddNode(id.MustNumID(201), "c1", mkPpn(id.MustNumID(101), "a")...)
	require.NoError(t, err)

	err = svc.AddNode(id.MustNumID(202), "c2", mkPpn(id.MustNumID(101), "b")...)
	require.NoError(t, err)

	err = svc.AddNode(id.MustNumID(203), "c3", mkPpn(id.MustNumID(101), "")...)
	require.NoError(t, err)

	//

	err = svc.AddNode(id.MustNumID(311), "cc11", mkPpn(id.MustNumID(201), "a", id.MustNumID(201), "b")...)
	require.NoError(t, err)

	err = svc.AddNode(id.MustNumID(312), "cc12", mkPpn(id.MustNumID(201), "")...)
	require.NoError(t, err)

	err = svc.AddNode(id.MustNumID(321), "cc21", mkPpn(id.MustNumID(202), "")...)
	require.NoError(t, err)

	err = svc.AddNode(id.MustNumID(331), "cc31", mkPpn(id.MustNumID(203), "a")...)
	require.NoError(t, err)

	err = svc.AddNode(id.MustNumID(332), "cc32", mkPpn(id.MustNumID(203), "b")...)
	require.NoError(t, err)

	//

	err = svc.AddNode(id.MustNumID(401), "ccc01", mkPpn(id.MustNumID(311), "")...)
	require.NoError(t, err)

	err = svc.AddNode(id.MustNumID(402), "ccc02", mkPpn(id.MustNumID(312), "b")...)
	require.NoError(t, err)

	err = svc.AddNode(id.MustNumID(403), "ccc03", mkPpn(id.MustNumID(332), "")...)
	require.NoError(t, err)

	//
	//
	//

	t.Run("201 any path", func(t *testing.T) {
		bb := svc.branchIndex[id.MustNumID(201)].inline()
		require.Len(t, bb, 5)

		require.Equal(t, id.MustNumID(201), bb[0].id)
		require.Equal(t, id.MustNumID(311), bb[1].id)
		require.Equal(t, id.MustNumID(312), bb[2].id)
		require.Equal(t, id.MustNumID(401), bb[3].id)
		require.Equal(t, id.MustNumID(402), bb[4].id)
	})

	t.Run("201 a path", func(t *testing.T) {
		bb := svc.branchIndex[id.MustNumID(201)].inline("a")
		require.Len(t, bb, 4)

		require.Equal(t, id.MustNumID(201), bb[0].id)
		require.Equal(t, id.MustNumID(311), bb[1].id)
		require.Equal(t, id.MustNumID(312), bb[2].id)
		require.Equal(t, id.MustNumID(401), bb[3].id)
	})

}

func TestUpdateNode(t *testing.T) {
	t.Run("update without move", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root")
		require.NoError(t, err)

		err = svc.UpdateNode(id.MustNumID(101), "root edited")
		require.NoError(t, err)

		require.Equal(t, id.MustNumID(101), svc.root.id)
		require.Equal(t, "root edited", svc.root.handle)
	})

	t.Run("change parent", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root")
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", mkPp(id.MustNumID(101))...)
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(202), "c2", mkPp(id.MustNumID(101))...)
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(301), "c3", mkPp(id.MustNumID(201))...)
		require.NoError(t, err)

		err = svc.UpdateNode(id.MustNumID(301), "c3 edited", mkPp(id.MustNumID(202))...)
		require.NoError(t, err)

		checkInline(t, svc.branchIndex[id.MustNumID(201)], id.MustNumID(201))
		checkInline(t, svc.branchIndex[id.MustNumID(202)], id.MustNumID(202), id.MustNumID(301))
	})

	t.Run("update non existing node", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root")
		require.NoError(t, err)

		err = svc.UpdateNode(id.MustNumID(999), "root edited", mkPp(id.MustNumID(0))...)
		require.Error(t, err)
	})
}

func TestRemoveNode(t *testing.T) {
	t.Run("remove leaf node", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root")
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", mkPp(id.MustNumID(101))...)
		require.NoError(t, err)

		err = svc.RemoveNode(id.MustNumID(201))
		require.NoError(t, err)

		checkInline(t, svc.branchIndex[id.MustNumID(101)], id.MustNumID(101))
		checkInline(t, svc.branchIndex[id.MustNumID(201)])
	})

	t.Run("can not remove non-lief nodes", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root")
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", mkPp(id.MustNumID(101))...)
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(301), "c2", mkPp(id.MustNumID(201))...)
		require.NoError(t, err)

		err = svc.RemoveNode(id.MustNumID(201))
		require.Error(t, err)
	})

	t.Run("can not remove non-existing nodes", func(t *testing.T) {
		svc, err := mkOrgTree()
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(101), "root")
		require.NoError(t, err)

		err = svc.AddNode(id.MustNumID(201), "c1", mkPp(id.MustNumID(101))...)
		require.NoError(t, err)

		err = svc.RemoveNode(id.MustNumID(999))
		require.Error(t, err)
	})
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
					paths:  mkPp(id.MustNumID(1)),
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
		out, _, err := buildOrgTree(&groupNode{id: id.MustNumID(1), handle: "root"})
		require.NoError(t, err)

		nodes := out.inline()
		require.Len(t, nodes, 1)
		require.Equal(t, "root", nodes[0].handle)
	})

	t.Run("multiple root error", func(t *testing.T) {
		_, _, err := buildOrgTree(
			&groupNode{id: id.MustNumID(1), handle: "root1"},
			&groupNode{id: id.MustNumID(2), handle: "root2"},
		)
		require.Error(t, err)
		require.EqualError(t, err, "multiple root nodes detected")
	})

	t.Run("missing parent error", func(t *testing.T) {
		_, _, err := buildOrgTree(
			&groupNode{id: id.MustNumID(1), handle: "root1"},
			&groupNode{id: id.MustNumID(2), handle: "child1", paths: mkPp(id.MustNumID(99))},
		)
		require.Error(t, err)
		require.EqualError(t, err, `node "2" parent "99" not found`)
	})

	t.Run("build tree", func(t *testing.T) {
		out, _, err := buildOrgTree(
			&groupNode{id: id.MustNumID(1), handle: "root"},
			&groupNode{id: id.MustNumID(2), handle: "child1", paths: mkPp(id.MustNumID(1))},
			&groupNode{id: id.MustNumID(3), handle: "child2", paths: mkPp(id.MustNumID(1))},
			&groupNode{id: id.MustNumID(4), handle: "child1_grandchild1", paths: mkPp(id.MustNumID(2))},
			&groupNode{id: id.MustNumID(5), handle: "child1_grandchild2", paths: mkPp(id.MustNumID(2))},
			&groupNode{id: id.MustNumID(6), handle: "child2_grandchild1", paths: mkPp(id.MustNumID(3))},
			&groupNode{id: id.MustNumID(7), handle: "child1_grandchild2_ggrandchild1", paths: mkPp(id.MustNumID(5))},
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
			paths:  mkPp(id.MustNumID(1)),
		},
		members: []id.ID{id.MustNumID(121), id.MustNumID(122), id.MustNumID(123)},
	}, GroupMembers{
		group: &groupNode{
			id:     id.MustNumID(3),
			handle: "subgroup2",
			paths:  mkPp(id.MustNumID(1)),
		},
		members: []id.ID{id.MustNumID(131)},
	}, GroupMembers{
		group: &groupNode{
			id:     id.MustNumID(4),
			handle: "subgroup1_subgroup1",
			paths:  mkPp(id.MustNumID(2)),
		},
		members: []id.ID{id.MustNumID(141), id.MustNumID(142)},
	})
}

func mkPp(ii ...id.ID) []GroupNodePath {
	out := make([]GroupNodePath, len(ii))
	for i, x := range ii {
		out[i] = GroupNodePath{
			SelfID: x,
		}
	}

	return out
}

func mkPpn(ii ...any) []GroupNodePath {
	if len(ii)%2 != 0 {
		panic("number of things not even")
	}

	out := make([]GroupNodePath, 0, len(ii))
	for i := 0; i < len(ii); i += 2 {
		id := ii[i].(id.ID)
		label := ii[i+1].(string)

		out = append(out, GroupNodePath{
			SelfID: id,
			Name:   label,
		})
	}

	return out
}

func checkInline(t *testing.T, node *groupNode, ii ...id.ID) {
	nn := node.inline()
	require.Len(t, nn, len(ii))

	for i, n := range nn {
		require.Equal(t, ii[i], n.id)
	}
}
