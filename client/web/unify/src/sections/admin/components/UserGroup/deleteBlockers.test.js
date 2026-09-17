import { describe, it, expect, vi } from 'vitest'
import { userGroupDeleteBlockers } from './deleteBlockers'

const group = (userGroupID, ...parents) => ({
  userGroupID,
  config: { path: parents.map(selfID => ({ selfID, name: '' })) },
})

function api({ members = [], pages = [[]] } = {}) {
  return {
    userGroupMemberList: vi.fn(async () => ({ set: members })),
    userGroupList: vi.fn(async ({ pageCursor }) => {
      const i = pageCursor ? Number(pageCursor) : 0
      return { set: pages[i], filter: { nextPage: i + 1 < pages.length ? String(i + 1) : '' } }
    }),
  }
}

describe('userGroupDeleteBlockers', () => {
  it('counts members and the groups that report to the group', async () => {
    const $SystemAPI = api({
      members: ['U1', 'U2'],
      pages: [[group('G1'), group('G2', 'G1'), group('G3', 'G2', 'G1'), group('G4', 'G9')]],
    })

    expect(await userGroupDeleteBlockers($SystemAPI, 'G1')).toEqual({ members: 2, childGroups: 2 })
    expect($SystemAPI.userGroupMemberList).toHaveBeenCalledWith({ userGroupID: 'G1' })
  })

  it('reads every page of groups', async () => {
    const $SystemAPI = api({
      pages: [[group('G2', 'G1')], [group('G3', 'G1')], [group('G4', 'G5')]],
    })

    expect(await userGroupDeleteBlockers($SystemAPI, 'G1')).toEqual({ members: 0, childGroups: 2 })
    expect($SystemAPI.userGroupList).toHaveBeenCalledTimes(3)
  })

  it('finds nothing for an empty group', async () => {
    expect(await userGroupDeleteBlockers(api({ pages: [[group('G1')]] }), 'G1')).toEqual({
      members: 0,
      childGroups: 0,
    })
  })
})
