import { describe, it, expect, beforeAll } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import CInputUserGroup from './CInputUserGroup.vue'

beforeAll(() => {
  // @ts-expect-error assigning a stub over a property jsdom does not define
  window.matchMedia = (q: string) => ({
    matches: false,
    media: q,
    onchange: null,
    addEventListener() {},
    removeEventListener() {},
    addListener() {},
    removeListener() {},
    dispatchEvent: () => false,
  })
})

const GROUPS = [
  { userGroupID: '1', meta: { short: 'Root' } },
  { userGroupID: '2', meta: { short: 'Sales' } },
  { userGroupID: '3', meta: { short: 'Support' } },
]

const api = {
  userGroupListCancellable: () => ({
    response: async () => ({ set: GROUPS }),
    cancel: () => {},
  }),
}

function offered(props: Record<string, unknown>) {
  const w = mount(CInputUserGroup, { props, global: { provide: { $SystemAPI: api } } })
  return flushPromises().then(() =>
    (w.findComponent({ name: 'Select' }).props('options') as typeof GROUPS).map(g => g.userGroupID),
  )
}

describe('CInputUserGroup excludeUserGroups', () => {
  it('offers every group by default', async () => {
    expect(await offered({})).toEqual(['1', '2', '3'])
  })

  it('leaves the excluded groups out', async () => {
    expect(await offered({ excludeUserGroups: ['2'] })).toEqual(['1', '3'])
  })
})
