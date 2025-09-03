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

func TestAssignGroupMembers(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		svc, err := OrgTree(
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
		svc, err := OrgTree(
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
	svc, err := OrgTree(
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

func testPrepState(t *testing.T) (svc *orgTree, err error) {
	return OrgTree(GroupMembers{
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
