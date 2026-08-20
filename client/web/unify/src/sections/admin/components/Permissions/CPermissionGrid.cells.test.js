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
const cellIcons = wrapper =>
  wrapper.findAll('div.w-48.text-lg').map(c => {
    const i = c.find('i')
    return i.exists() ? i.attributes('class') : ''
  })

describe('CPermissionGrid cell states', () => {
  beforeEach(() => localStorage.clear())

  it('renders an edit column as green allow, red deny and a blank inherit', async () => {
    const wrapper = await mountGrid({
      mode: 'edit',
      ID: 'edit-R1',
      roleID: 'R1',
      name: ['Auditor'],
    })

    const [read, update, del] = cellIcons(wrapper)
    expect(read).toContain('pi-check')
    expect(read).toContain('text-green-500')
    expect(update).toContain('pi-times')
    expect(update).toContain('text-red-500')
    expect(del).toBe('')
  })

  it('never leaves an evaluation column blank: inherit evaluates to a deny', async () => {
    const wrapper = await mountGrid({
      mode: 'eval',
      ID: 'eval-R1',
      roleID: ['R1'],
      userID: null,
      name: ['Auditor'],
    })

    const icons = cellIcons(wrapper)
    expect(icons).toHaveLength(3)
    expect(icons.every(c => c !== '')).toBe(true)
    expect(icons[2]).toContain('pi-times')
    expect(icons[2]).toContain('text-red-500')
  })
})
