import { describe, it, expect, beforeAll } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import CInputRole from './CInputRole.vue'

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

const role = (roleID: string) => ({ roleID, name: `Role ${roleID}` })

type VM = { options: Array<{ roleID: string }>; selectedRoles: string[] }

// The list is gated so the test decides whether it lands before or after the
// picked role is resolved — the order is what the bug turned on.
function mountRole(modelValue: string[], listed = [role('R2'), role('R3')]) {
  let release: () => void = () => {}
  const gate = new Promise<void>(r => {
    release = r
  })

  const api = {
    roleListCancellable: () => ({
      response: async () => {
        await gate
        return { set: listed }
      },
      cancel: () => {},
    }),
    roleRead: async ({ roleID }: { roleID: string }) => role(roleID),
  }

  const w = mount(CInputRole, {
    props: { multiple: true, modelValue },
    global: { provide: { $SystemAPI: api } },
  })
  return { w, release, vm: () => w.vm as unknown as VM }
}

describe('CInputRole multiple', () => {
  it('keeps a picked role when the role list lands after it', async () => {
    const { release, vm } = mountRole(['R1'])

    // The picked role resolves first and goes into the options...
    await flushPromises()
    expect(vm().options.map(r => r.roleID)).toContain('R1')

    // ...then the list arrives and replaces them wholesale.
    release()
    await flushPromises()

    expect(
      vm()
        .options.map(r => r.roleID)
        .sort(),
    ).toEqual(['R1', 'R2', 'R3'])
    expect(vm().selectedRoles).toEqual(['R1'])
  })

  it('lists a picked role once even though several passes pin it', async () => {
    const { release, vm } = mountRole(['R2'])

    release()
    await flushPromises()

    const ids = vm().options.map(r => r.roleID)
    expect(ids.filter(id => id === 'R2')).toHaveLength(1)
  })

  it('emits the picked roles as an array of IDs', async () => {
    const { w, release } = mountRole([])
    release()
    await flushPromises()

    w.findComponent({ name: 'MultiSelect' }).vm.$emit('update:modelValue', ['R2', 'R3'])
    await flushPromises()

    expect(w.emitted('update:modelValue')?.at(-1)).toEqual([['R2', 'R3']])
  })
})
