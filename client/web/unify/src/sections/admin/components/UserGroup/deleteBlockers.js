// What keeps a user group from being deleted: members who are not deleted, and
// groups that report to it. The server refuses the delete on the same counts.
export async function userGroupDeleteBlockers($SystemAPI, userGroupID) {
  const members = await $SystemAPI.userGroupMemberList({ userGroupID })

  let childGroups = 0
  let pageCursor
  do {
    const { set = [], filter = {} } = await $SystemAPI.userGroupList({ limit: 100, pageCursor })
    childGroups += set.filter(g =>
      (g.config?.path || []).some(p => p.selfID === userGroupID),
    ).length
    pageCursor = filter.nextPage
  } while (pageCursor)

  return { members: (members?.set || members || []).length, childGroups }
}
