package rbac

import (
	"errors"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/id"
)

type (
	orgTree struct {
		root *groupNode

		// branchIndex provides all of the available sub-trees in the given orgTree
		//
		// This is primarily meant to speed up the process of determining what
		// permissions should be included in access evaluation
		//
		// @todo benchmark with/without to see if it makes sense
		branchIndex map[id.ID]*groupNode

		// memberGroupIndex holds what user group a given member belongs to
		memberGroupIndex map[id.ID]*groupNode

		// groupMemberIndex holds what members a specific user group holds
		groupMemberIndex map[id.ID][]id.ID
	}

	groupNode struct {
		id       id.ID
		handle   string
		selfID   id.ID
		children []*groupNode

		roles []id.ID
	}

	GroupMembers struct {
		group   *groupNode
		members []id.ID
	}
)

// ConvUserGroup creates internal structures for the org tree
func ConvUserGroup(id id.ID, handle string, selfID id.ID, members, roles []id.ID) GroupMembers {
	return GroupMembers{
		group: &groupNode{
			id:     id,
			handle: handle,
			selfID: selfID,
			roles:  roles,
		},

		members: members,
	}
}

// OrgTree prepares the organization tree subservice
func OrgTree(gm ...GroupMembers) (svc *orgTree, err error) {
	svc = &orgTree{}
	err = svc.Rebuild(gm...)
	return
}

// Rebuild rebuilds the current organization tree dropping all currently indexed data
func (svc *orgTree) Rebuild(gm ...GroupMembers) (err error) {
	svc.root = nil
	svc.branchIndex = nil
	svc.memberGroupIndex = nil
	svc.groupMemberIndex = nil

	gg := make([]*groupNode, 0, len(gm))
	for _, m := range gm {
		gg = append(gg, m.group)
	}

	svc.root, svc.branchIndex, err = buildOrgTree(gg...)
	if err != nil {
		return
	}

	for _, m := range gm {
		err = svc.AssignGroupMembers(m.group.id, m.members...)
		if err != nil {
			return
		}
	}

	return
}

func (svc *orgTree) AddNode(id id.ID, handle string, selfID id.ID) (err error) {
	err = svc.addNode(&groupNode{
		id:     id,
		handle: handle,
		selfID: selfID,
	})
	if err != nil {
		return
	}

	return err
}

func (svc *orgTree) addNode(node *groupNode) (err error) {
	if svc.root == nil {
		svc.root = node

		if svc.branchIndex == nil {
			svc.branchIndex = make(map[id.ID]*groupNode)
		}
		svc.branchIndex[node.id] = node

		if svc.groupMemberIndex == nil {
			svc.groupMemberIndex = make(map[id.ID][]id.ID)
		}
		svc.groupMemberIndex[node.id] = []id.ID{}

		return
	}

	i, n := svc.findNode(node.selfID)
	if i < 0 {
		return fmt.Errorf("cannot insert node %v (%s): parent node %v not found", node.id.Value(), node.handle, node.selfID.Value())
	}

	n.children = append(n.children, node)

	svc.branchIndex[node.id] = node
	svc.groupMemberIndex[node.id] = []id.ID{}

	return
}

func (svc *orgTree) UpdateNode(idx id.ID, handle string, selfID id.ID) (err error) {
	i, n := svc.findNode(idx)
	if i < 0 {
		return fmt.Errorf("cannot update node %v (%s): not indexed", idx.Value(), handle)
	}

	oldSelfID := n.selfID

	n.handle = handle
	n.selfID = selfID

	if svc.branchIndex == nil {
		svc.branchIndex = make(map[id.ID]*groupNode)
	}

	if !oldSelfID.Equal(id.MustNumID(0)) {
		i = svc.isInSlice(svc.branchIndex[oldSelfID].children, idx)
		svc.branchIndex[oldSelfID].children = append(svc.branchIndex[oldSelfID].children[:i], svc.branchIndex[oldSelfID].children[i+1:]...)
	}

	if !selfID.Equal(id.MustNumID(0)) {
		svc.branchIndex[selfID].children = append(svc.branchIndex[selfID].children, n)
	}

	return
}

func (svc *orgTree) RemoveNode(id id.ID) (err error) {
	i, n := svc.findNode(id)
	if i < 0 {
		return fmt.Errorf("cannot delete node %v: not indexed", id.Value())
	}

	if len(svc.groupMemberIndex[id]) > 0 {
		return fmt.Errorf("cannot remove node %v: node has members", id)
	}

	if len(n.children) > 0 {
		return fmt.Errorf("cannot remove node %v: node is above a node ", id)
	}

	for _, n := range svc.root.inline() {
		i := svc.isInSlice(n.children, id)
		if i < 0 {
			continue
		}

		n.children = append(n.children[:i], n.children[i+1:]...)
	}

	delete(svc.branchIndex, id)
	delete(svc.groupMemberIndex, id)

	return
}

// GroupMembers returns all of the members belonging to the specific group
func (svc *orgTree) GroupMembers(group id.ID) (members []id.ID) {
	if svc.memberGroupIndex == nil {
		return
	}

	return svc.groupMemberIndex[group]
}

// IsAbove checks if parentUser is above the childUser
func (svc *orgTree) IsAbove(childUser, parentUser id.ID) bool {
	childGroup := svc.memberGroupIndex[childUser]
	parentGroup := svc.memberGroupIndex[parentUser]

	for _, c := range parentGroup.inline()[1:] {
		if c.id.Equal(childGroup.id) {
			return true
		}
	}

	return false
}

// MemberBranch returns the branch belonging to the user
func (svc *orgTree) MemberBranch(user id.ID) (group []*groupNode, err error) {
	if svc == nil {
		return
	}

	if svc.memberGroupIndex == nil {
		return
	}

	aux, ok := svc.memberGroupIndex[user]
	if !ok {
		return nil, fmt.Errorf("user %s is not a member of any group", user.Value())
	}

	ss, ok := svc.branchIndex[aux.id]
	if !ok {
		return nil, fmt.Errorf("group %s not found in subtree index", aux.id.Value())
	}

	return ss.inline(), nil
}

func (svc *orgTree) AddGroupRole(group id.ID, roles ...id.ID) (err error) {
	i, n := svc.findNode(group)
	if i < 0 {
		return fmt.Errorf("node %v not found", group)
	}

	for _, r := range roles {
		if id.InSlice(r, n.roles...) {
			continue
		}

		n.roles = append(n.roles, r)
	}

	return
}

func (svc *orgTree) RemoveGroupRole(group id.ID, roles ...id.ID) (err error) {
	i, n := svc.findNode(group)
	if i < 0 {
		return fmt.Errorf("node %v not found", group)
	}

	for _, r := range roles {
		n.roles = id.RemoveFromSlice(r, n.roles...)
	}

	return
}

// AssignGroupMembers assigns the members to the specified group
func (svc *orgTree) AssignGroupMembers(group id.ID, members ...id.ID) (err error) {
	if group.IsZero() {
		return
	}

	if svc.memberGroupIndex == nil {
		svc.memberGroupIndex = make(map[id.ID]*groupNode)
	}

	if svc.groupMemberIndex == nil {
		svc.groupMemberIndex = make(map[id.ID][]id.ID)
	}

	g := svc.branchIndex[group]
	if g == nil {
		return fmt.Errorf("group %s not found", group.Value())
	}

	for _, u := range members {
		oldGroup := svc.memberGroupIndex[u]

		// Cleanup previous state
		if oldGroup != nil {
			svc.groupMemberIndex[oldGroup.id] = id.RemoveFromSlice(u, svc.groupMemberIndex[oldGroup.id]...)
		}

		// Handle new state
		svc.memberGroupIndex[u] = g

		if !id.InSlice(u, svc.groupMemberIndex[group]...) {
			svc.groupMemberIndex[group] = append(svc.groupMemberIndex[group], u)
		}

	}

	return
}

// RemoveGroupMembers removes the members from the given user group
func (svc *orgTree) RemoveGroupMembers(group id.ID, members ...id.ID) (err error) {
	if svc == nil {
		return
	}

	if svc.memberGroupIndex == nil {
		return
	}

	for _, m := range members {
		g, ok := svc.memberGroupIndex[m]
		if !ok {
			err = fmt.Errorf("user %s is not a member of any group", m.Value())
			return
		}

		if !g.id.Equal(group) {
			err = fmt.Errorf("user %s is not a member of group %s", m.Value(), group.Value())
			return
		}

		// Cleanup
		delete(svc.memberGroupIndex, m)
		svc.groupMemberIndex[group] = id.RemoveFromSlice(m, svc.groupMemberIndex[group]...)
	}

	return
}

func buildOrgTree(gg ...*groupNode) (root *groupNode, index map[id.ID]*groupNode, err error) {
	mapped := make(map[id.ID]*groupNode, len(gg))

	for _, g := range gg {
		mapped[g.id] = g

		if g.selfID.IsZero() {
			if root != nil {
				return nil, nil, errors.New("multiple root nodes detected")
			}

			root = g
		}
	}

	for _, g := range gg {
		if g.selfID.IsZero() {
			continue
		}

		parent, ok := mapped[g.selfID]
		if !ok {
			return nil, nil, fmt.Errorf("node %s parent %s not found", g.id.Value(), g.selfID.Value())
		}

		parent.children = append(parent.children, g)
	}

	index = make(map[id.ID]*groupNode, len(gg))
	for _, g := range gg {
		index[g.id] = g
	}

	return root, index, nil
}

// inline converts the subtree in a slice using BFS
func (root *groupNode) inline() []*groupNode {
	if root == nil {
		return nil
	}

	out := make([]*groupNode, 0)
	fifo := make([]*groupNode, 0)

	fifo = append(fifo, root)
	i := 0
	for {
		n := fifo[i]
		fifo = append(fifo, n.children...)

		out = append(out, n)

		if i+1 >= len(fifo) {
			break
		}

		i++
	}

	return out
}

func (svc *orgTree) findNode(id id.ID) (i int, n *groupNode) {
	nn := svc.root.inline()

	i = 0
	found := false
	for i = range nn {
		n = nn[i]

		if n.id.Equal(id) {
			found = true
			break
		}
	}

	if !found {
		n = nil
		i = -1
	}

	return
}

func (svc *orgTree) isInSlice(ss []*groupNode, id id.ID) int {
	for i, s := range ss {
		if s.id.Equal(id) {
			return i
		}
	}

	return -1
}
