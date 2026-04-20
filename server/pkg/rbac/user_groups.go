package rbac

import (
	"errors"
	"fmt"

	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/stack"
	"go.uber.org/zap"
)

type (
	orgTree struct {
		logger *zap.Logger

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
		id     id.ID
		handle string

		// Nodes this one is parent of
		children []groupNodeConnection
		// Nodes this one is child of
		//
		// Used to simplify lookup
		parents []groupNodeConnection

		// Slice of paths this node was provided (this node -> parents)
		paths []GroupNodePath

		roles []id.ID
	}

	GroupNodePath struct {
		SelfID id.ID
		Name   string
	}

	groupNodeConnection struct {
		node *groupNode
		name string
	}

	GroupMembers struct {
		group   *groupNode
		members []id.ID
	}
)

// ConvUserGroup creates internal structures for the org tree
func ConvUserGroup(id id.ID, handle string, members, roles []id.ID, paths []GroupNodePath) GroupMembers {
	return GroupMembers{
		group: &groupNode{
			id:     id,
			handle: handle,
			roles:  roles,
			paths:  paths,
		},

		members: members,
	}
}

// OrgTree prepares the organization tree subservice
func OrgTree(log *zap.Logger, gm ...GroupMembers) (svc *orgTree, err error) {
	svc, err = mkOrgTree(gm...)
	if err != nil {
		return
	}

	svc.logger = log

	err = svc.Rebuild(gm...)
	return
}

func mkOrgTree(gm ...GroupMembers) (svc *orgTree, err error) {
	svc = &orgTree{
		logger: zap.NewNop(),
	}

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

func (svc *orgTree) AddNode(id id.ID, handle string, paths ...GroupNodePath) (err error) {
	svc.logger.Debug("adding node", zap.Any("id", id.Value()), zap.String("handle", handle), zap.Any("paths", formatPaths(paths)))

	n := &groupNode{
		id:     id,
		handle: handle,
		paths:  paths,
	}

	err = svc.addNode(n)
	if err != nil {
		return
	}

	return
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

	for _, p := range node.paths {
		i, n := svc.findNode(p.SelfID)
		if i < 0 {
			return fmt.Errorf("cannot add node path: parent node %v not found", p.SelfID)
		}

		n.children = append(n.children, groupNodeConnection{
			node: node,
			name: p.Name,
		})

		node.parents = append(node.parents, groupNodeConnection{
			node: n,
			name: p.Name,
		})
	}

	svc.branchIndex[node.id] = node
	svc.groupMemberIndex[node.id] = []id.ID{}

	return
}

func (svc *orgTree) UpdateNode(idx id.ID, handle string, paths ...GroupNodePath) (err error) {
	svc.logger.Debug("updating node", zap.Any("id", idx.Value()), zap.String("handle", handle), zap.Any("paths", formatPaths(paths)))

	i, n := svc.findNode(idx)
	if i < 0 {
		return fmt.Errorf("cannot update node %v (%s): not indexed", idx.Value(), handle)
	}

	oldPaths := n.paths

	n.handle = handle
	n.paths = paths

	if svc.branchIndex == nil {
		svc.branchIndex = make(map[id.ID]*groupNode)
	}

	for _, op := range oldPaths {
		i = svc.isInConnections(svc.branchIndex[op.SelfID].children, idx)
		svc.branchIndex[op.SelfID].children = append(svc.branchIndex[op.SelfID].children[:i], svc.branchIndex[op.SelfID].children[i+1:]...)

		i = svc.isInConnections(svc.branchIndex[idx].parents, op.SelfID)
		svc.branchIndex[idx].parents = append(svc.branchIndex[idx].parents[:i], svc.branchIndex[idx].parents[i+1:]...)
	}

	for _, p := range paths {
		svc.branchIndex[p.SelfID].children = append(svc.branchIndex[p.SelfID].children, groupNodeConnection{
			node: n,
			name: p.Name,
		})

		svc.branchIndex[idx].parents = append(svc.branchIndex[idx].parents, groupNodeConnection{
			node: svc.branchIndex[p.SelfID],
			name: p.Name,
		})
	}

	return
}

func (svc *orgTree) RemoveNode(id id.ID) (err error) {
	svc.logger.Debug("removing node", zap.Any("id", id.Value()))

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
		i := svc.isInConnections(n.children, id)
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
func (svc *orgTree) IsAbove(childUser, parentUser id.ID, paths ...string) bool {
	childGroup := svc.memberGroupIndex[childUser]
	parentGroup := svc.memberGroupIndex[parentUser]

	return svc.isAbove(parentGroup, childGroup, paths...)
}

// isAbove checks if parent is above child based on the hierarchy
//
// The check goes bottom up as the number of paths will (in most cases) always be lower.
func (svc *orgTree) isAbove(parent, child *groupNode, paths ...string) (ok bool) {
	checkConnection := func(c groupNodeConnection) bool {
		if len(paths) == 0 {
			return true
		}

		if c.name == "" {
			return true
		}

		for _, pth := range paths {
			if c.name == pth {
				return true
			}
		}

		return false
	}

	stack := stack.Stack[*groupNode]{}
	for _, p := range child.parents {
		if !checkConnection(p) {
			continue
		}

		stack.Push(p.node)
	}

	for !stack.Empty() {
		n := stack.Pop()
		if n.id.Equal(parent.id) {
			return true
		}

		for _, c := range n.parents {
			if !checkConnection(c) {
				continue
			}

			stack.Push(c.node)
		}
	}

	return false
}

// MemberBranch returns the branch belonging to the user
func (svc *orgTree) MemberBranch(user id.ID, paths ...string) (group []*groupNode, err error) {
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

	return ss.inline(paths...), nil
}

func (svc *orgTree) AddGroupRole(group id.ID, roles ...id.ID) (err error) {
	svc.logger.Debug("adding group roles", zap.Any("group", group.Value()), zap.Any("roles", id.StringifySlice(roles...)))

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
	svc.logger.Debug("removing group roles", zap.Any("group", group.Value()), zap.Any("roles", id.StringifySlice(roles...)))

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

	svc.logger.Debug("assigning group members", zap.Any("group", group.Value()), zap.Any("roles", id.StringifySlice(members...)))

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

	svc.logger.Debug("assigning group members", zap.Any("group", group.Value()), zap.Any("roles", id.StringifySlice(members...)))

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

		if len(g.paths) == 0 {
			if root != nil {
				return nil, nil, errors.New("multiple root nodes detected")
			}

			root = g
		}
	}

	for _, g := range gg {
		if len(g.paths) == 0 {
			continue
		}

		for _, p := range g.paths {
			parent, ok := mapped[p.SelfID]
			if !ok {
				return nil, nil, fmt.Errorf("node %s parent %s not found", g.id.Value(), p.SelfID.Value())
			}

			parent.children = append(parent.children, groupNodeConnection{
				name: p.Name,
				node: g,
			})

			g.parents = append(g.parents, groupNodeConnection{
				node: parent,
				name: p.Name,
			})
		}
	}

	index = make(map[id.ID]*groupNode, len(gg))
	for _, g := range gg {
		index[g.id] = g
	}

	return root, index, nil
}

// inline converts the subtree in a slice using BFS
func (root *groupNode) inline(paths ...string) []*groupNode {
	if root == nil {
		return nil
	}

	// Helper to check if a connection should be included based on path
	//
	// - It we don't specify any path, all nodes are used
	// - If a connection name is empty, the node is included regardless of defined paths
	// - If a connection defines a name, it must occur in the requested paths slice
	inclConnection := func(c groupNodeConnection) bool {
		if len(paths) == 0 {
			return true
		}

		if c.name == "" {
			return true
		}

		for _, pth := range paths {
			if c.name == pth {
				return true
			}
		}

		return false
	}

	// outSeen helps us filter out duplicates
	outSeen := make(map[id.ID]bool, 4)
	out := make([]*groupNode, 0)
	fifo := make([]*groupNode, 0)

	fifo = append(fifo, root)
	i := 0
	for {
		n := fifo[i]
		if outSeen[n.id] {
			if i+1 >= len(fifo) {
				break
			}

			i++

			continue
		}

		for _, c := range n.children {
			if !inclConnection(c) {
				continue
			}

			fifo = append(fifo, c.node)
		}

		out = append(out, n)
		outSeen[n.id] = true

		if i+1 >= len(fifo) {
			break
		}

		i++
	}

	return out
}

func formatPaths(paths []GroupNodePath) []map[string]any {
	out := make([]map[string]any, 0, len(paths))

	for _, p := range paths {
		out = append(out, map[string]any{
			"selfID": p.SelfID,
			"name":   p.Name,
		})
	}

	return out
}

func (n *groupNode) format() map[string]any {
	out := map[string]any{
		"id":     n.id.Value(),
		"handle": n.handle,
		"paths":  formatPaths(n.paths),
		"roles":  id.StringifySlice(n.roles...),
	}

	cc := make([]id.ID, 0, len(n.children))
	for _, c := range n.children {
		cc = append(cc, c.node.id)
	}

	out["children"] = id.StringifySlice(cc...)

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

func (svc *orgTree) isInConnections(ss []groupNodeConnection, id id.ID) int {
	for i, s := range ss {
		if s.node.id.Equal(id) {
			return i
		}
	}

	return -1
}
