import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k, te: () => false }) }))

vi.mock('@planetcrust/human-vue', () => ({
  components: {
    CInputRole: { name: 'CInputRole', props: ['modelValue'], template: '<div />' },
  },
}))

import CPermissionGrid from './CPermissionGrid.vue'

const RES = 'corteza::system:project/*'
const OPS = ['read', 'update', 'delete']

// permissionsRead returns only the rules a role actually stores: `delete` has
// none, which is the inherit case. permissionsTrace always answers for every
// operation, and answers `inherit`/no-match where nothing matched.
const api = {
  permissionsList: async () => OPS.map(op => ({ type: 'corteza::system:project', any: RES, op })),
  permissionsRead: async () => [
    { resource: RES, operation: 'read', access: 'allow' },
    { resource: RES, operation: 'update', access: 'deny' },
  ],
  permissionsTrace: async () => [
    { resource: RES, operation: 'read', access: 'allow' },
    { resource: RES, operation: 'update', access: 'deny' },
    { resource: RES, operation: 'delete', access: 'inherit', resolution: 'no-match' },
  ],
}

const stub = (name, props = []) => ({ name, props, template: '<div><slot /></div>' })

const EDIT_COLUMN = { mode: 'edit', ID: 'edit-R1', roleID: 'R1', name: ['Auditor'] }
const EVAL_COLUMN = { mode: 'eval', ID: 'eval-R1', roleID: ['R1'], userID: null, name: ['Auditor'] }

async function mountGrid(column) {
  localStorage.setItem('permissionList.roles', JSON.stringify([column]))

  const wrapper = mount(CPermissionGrid, {
    props: { api, component: 'system' },
    global: {
      provide: {
        $SystemAPI: { userList: async () => ({ set: [] }) },
        $toast: { toastSuccess() {}, toastErrorHandler: () => () => {} },
      },
      mocks: { $t: k => k },
      stubs: {
        Card: { name: 'Card', template: '<div><slot name="content" /></div>' },
        Dialog: { name: 'Dialog', template: '<div />' },
        ProgressSpinner: stub('ProgressSpinner'),
        Message: stub('Message'),
        Button: stub('Button', ['label', 'icon']),
        SelectButton: stub('SelectButton', ['modelValue', 'options']),
        AutoComplete: stub('AutoComplete', ['modelValue', 'suggestions']),
        Chip: stub('Chip', ['label']),
        CFormGroup: stub('CFormGroup', ['label']),
        CEditorActions: stub('CEditorActions'),
      },
    },
  })

  await flushPromises()
  await flushPromises()
  return wrapper
}

// One cell per operation row; the icon is the whole of the cell's content.
const cells = wrapper => wrapper.findAll('div.w-48.text-lg')

const cellIcons = wrapper =>
  cells(wrapper).map(c => {
    const i = c.find('i')
    return i.exists() ? i.attributes('class') : ''
  })

// What a cell says it is: its icon, or 'blank' where the rule is inherited.
const state = cell => {
  const i = cell.find('i')
  if (!i.exists()) return 'blank'
  return i.attributes('class').includes('pi-check') ? 'allow' : 'deny'
}

const isChanged = cell => cell.attributes('class').includes('bg-amber-500/15')

const clickThrough = async (cell, times) => {
  const seen = []
  for (let n = 0; n < times; n++) {
    await cell.trigger('click')
    seen.push(state(cell))
  }
  return seen
}

describe('CPermissionGrid cell states', () => {
  beforeEach(() => localStorage.clear())

  it('renders an edit column as green allow, red deny and a blank inherit', async () => {
    const wrapper = await mountGrid(EDIT_COLUMN)

    const [read, update, del] = cellIcons(wrapper)
    expect(read).toContain('pi-check')
    expect(read).toContain('text-green-500')
    expect(update).toContain('pi-times')
    expect(update).toContain('text-red-500')
    expect(del).toBe('')
  })

  it('never leaves an evaluation column blank: inherit evaluates to a deny', async () => {
    const wrapper = await mountGrid(EVAL_COLUMN)

    const icons = cellIcons(wrapper)
    expect(icons).toHaveLength(3)
    expect(icons.every(c => c !== '')).toBe(true)
    expect(icons[2]).toContain('pi-times')
    expect(icons[2]).toContain('text-red-500')
  })

  it('cycles a cell inherit -> allow -> deny -> inherit', async () => {
    const wrapper = await mountGrid(EDIT_COLUMN)
    const [, , del] = cells(wrapper)

    expect(state(del)).toBe('blank')
    expect(await clickThrough(del, 4)).toEqual(['allow', 'deny', 'blank', 'allow'])
  })

  it('cycles a cell that starts on an explicit rule from where it stands', async () => {
    const wrapper = await mountGrid(EDIT_COLUMN)
    const [read] = cells(wrapper)

    expect(state(read)).toBe('allow')
    expect(await clickThrough(read, 3)).toEqual(['deny', 'blank', 'allow'])
  })

  it('tints a cell amber while it differs from what was loaded', async () => {
    const wrapper = await mountGrid(EDIT_COLUMN)
    const [, , del] = cells(wrapper)

    expect(isChanged(del)).toBe(false)

    await del.trigger('click') // allow
    expect(isChanged(del)).toBe(true)
    await del.trigger('click') // deny
    expect(isChanged(del)).toBe(true)

    // Round-tripped back to the loaded value: no longer a pending change.
    await del.trigger('click') // inherit
    expect(isChanged(del)).toBe(false)
  })

  it('leaves an evaluation column inert', async () => {
    const wrapper = await mountGrid(EVAL_COLUMN)
    const [, , del] = cells(wrapper)

    expect(await clickThrough(del, 2)).toEqual(['deny', 'deny'])
    expect(isChanged(del)).toBe(false)
  })
})
