import { describe, it, expect, vi, beforeAll, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { defineComponent } from 'vue'

// The evaluation columns say who a comparison is against; the accesses they
// show belong to whichever resource the dialog was opened for. So what survives
// a close is the identity, and reopening re-traces it against the new resource.

beforeAll(() => {
  // PrimeVue overlays bind a matchMedia listener on mount; jsdom has none.
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

vi.mock('../input/CInputRole.vue', () => ({ default: { template: '<div />' } }))
vi.mock('../input/CInputUser.vue', () => ({ default: { template: '<div />' } }))

import CPermissionsDialog from './CPermissionsDialog.vue'
import { providePermissions, type PermissionsContext } from '../../composables/usePermissions'

// The dialog is mounted once per app and driven by the composable's state; a
// host that provides it is the honest way in.
let ctx: PermissionsContext
const Host = defineComponent({
  components: { CPermissionsDialog },
  setup() {
    ctx = providePermissions()
    return {}
  },
  template: '<CPermissionsDialog />',
})

const EVAL_KEY = 'permissionsDialog.evalColumns'

const ROLES: Record<string, { roleID: string; name: string }> = {
  R1: { roleID: 'R1', name: 'Administrator' },
  R2: { roleID: 'R2', name: 'Tester' },
}

const permissionsTrace = vi.fn(async () => [{ operation: 'read', access: 'allow' }])
const roleRead = vi.fn(async ({ roleID }: { roleID: string }) => {
  if (!ROLES[roleID]) throw new Error('not found')
  return ROLES[roleID]
})
const userRead = vi.fn(async ({ userID }: { userID: string }) => ({ userID, name: 'Ada' }))

const $SystemAPI = {
  permissionsList: vi.fn(async () => [{ op: 'read' }, { op: 'update' }]),
  permissionsRead: vi.fn(async () => []),
  permissionsTrace,
  permissionsUpdate: vi.fn(),
  roleList: vi.fn(async () => ({ set: [{ roleID: 'R9', name: 'Editor' }] })),
  roleRead,
  userRead,
}

let wrapper: ReturnType<typeof mount> | null = null

async function openDialog(resource: string) {
  if (!wrapper) {
    wrapper = mount(Host, {
      global: {
        directives: { tooltip: {} },
        provide: {
          $SystemAPI,
          $ComposeAPI: $SystemAPI,
          $toast: { toastSuccess: vi.fn(), toastErrorHandler: () => vi.fn() },
        },
      },
    })
  }
  ctx.open({ resource })
  await flushPromises()
  return wrapper
}

/** The dialog's own state, which is what these assertions are about. */
function vm() {
  return wrapper!.findComponent(CPermissionsDialog).vm as unknown as {
    evaluate: Array<{ roleIDs: string[]; userID: string | null; roleNames: string[] }>
    onAddEvalColumn: () => Promise<void>
    removeEvalColumn: (_i: number) => void
    addEval: { roleIDs: unknown[]; userID: string | null }
  }
}

async function addColumn(identity: { roleIDs?: string[]; userID?: string | null }) {
  vm().addEval.roleIDs = identity.roleIDs ?? []
  vm().addEval.userID = identity.userID ?? null
  await vm().onAddEvalColumn()
  await flushPromises()
}

function stored() {
  return JSON.parse(localStorage.getItem(EVAL_KEY) || 'null')
}

beforeEach(() => {
  localStorage.clear()
  permissionsTrace.mockClear()
  roleRead.mockClear()
})

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
})

describe('CPermissionsDialog evaluation columns', () => {
  it('stores only the identity of an added column, never its accesses', async () => {
    await openDialog('corteza::compose:module/NS1/M1')
    await addColumn({ roleIDs: ['R1'] })

    expect(vm().evaluate).toHaveLength(1)
    expect(stored()).toEqual([{ roleIDs: ['R1'], userID: null }])
  })

  it('restores the columns on the next open, traced against the new resource', async () => {
    await openDialog('corteza::compose:module/NS1/M1')
    await addColumn({ roleIDs: ['R1'] })
    await addColumn({ userID: 'U1' })

    permissionsTrace.mockClear()
    await openDialog('corteza::system:user/U9')

    expect(vm().evaluate.map(c => c.roleNames[0] ?? c.userID)).toEqual(['Administrator', 'U1'])
    for (const call of permissionsTrace.mock.calls) {
      expect((call[0] as { resource: string }).resource).toBe('corteza::system:user/U9')
    }
  })

  it('forgets a column that was removed', async () => {
    await openDialog('corteza::compose:module/NS1/M1')
    await addColumn({ roleIDs: ['R1'] })
    await addColumn({ roleIDs: ['R2'] })

    vm().removeEvalColumn(0)
    expect(stored()).toEqual([{ roleIDs: ['R2'], userID: null }])
  })

  it('drops a stored column whose role no longer exists, and stops asking for it', async () => {
    localStorage.setItem(
      EVAL_KEY,
      JSON.stringify([
        { roleIDs: ['R1'], userID: null },
        { roleIDs: ['GONE'], userID: null },
      ]),
    )

    await openDialog('corteza::compose:module/NS1/M1')

    expect(vm().evaluate.map(c => c.roleNames[0])).toEqual(['Administrator'])
    expect(stored()).toEqual([{ roleIDs: ['R1'], userID: null }])
  })

  it('restores at most four columns', async () => {
    localStorage.setItem(
      EVAL_KEY,
      JSON.stringify(Array.from({ length: 6 }, () => ({ roleIDs: ['R1'], userID: null }))),
    )

    await openDialog('corteza::compose:module/NS1/M1')

    expect(vm().evaluate).toHaveLength(4)
  })

  it('ignores a corrupt stored value instead of throwing', async () => {
    localStorage.setItem(EVAL_KEY, '{not json')

    await openDialog('corteza::compose:module/NS1/M1')

    expect(vm().evaluate).toEqual([])
  })

  it('ignores a stored value that parses but is not a list of columns', async () => {
    localStorage.setItem(EVAL_KEY, '{"roleIDs":["R1"]}')

    await openDialog('corteza::compose:module/NS1/M1')

    expect(vm().evaluate).toEqual([])
  })
})
