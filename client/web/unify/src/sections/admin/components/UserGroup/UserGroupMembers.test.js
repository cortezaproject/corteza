import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

const h = vi.hoisted(() => ({ confirm: null }))

vi.mock('@planetcrust/human-vue', () => ({
  components: {
    CInputUser: {
      name: 'CInputUser',
      props: { excludeUsers: { type: Array, default: () => [] } },
      emits: ['select'],
      template: '<div data-testid="add-picker" />',
    },
    CInputUserGroup: {
      name: 'CInputUserGroup',
      props: { modelValue: null, excludeUserGroups: { type: Array, default: () => [] } },
      emits: ['update:modelValue', 'select'],
      template: '<div data-testid="group-picker" />',
    },
  },
}))

vi.mock('@planetcrust/human-js', () => ({
  system: {
    User: class {
      constructor(u) {
        Object.assign(this, u)
      }
    },
  },
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: (k, p) => (p ? `${k} ${JSON.stringify(p)}` : k) }),
}))

vi.mock('primevue/useconfirm', () => ({
  useConfirm: () => ({ require: opts => (h.confirm = opts) }),
}))

import UserGroupMembers from './UserGroupMembers.vue'

const ENG = { userGroupID: 'G1', meta: { short: 'Engineering' }, canManageMembersOnUserGroup: true }
const SALES = { userGroupID: 'G2', meta: { short: 'Sales' } }

const USERS = {
  U1: { userID: 'U1', name: 'Jane Doe', userGroupID: 'G1' },
  U2: { userID: 'U2', name: 'John Roe', userGroupID: 'G1' },
}

let api
let toast

function mountMembers(userGroup = ENG) {
  return mount(UserGroupMembers, {
    props: { userGroup },
    global: {
      mocks: { $t: k => k },
      provide: { $SystemAPI: api, $toast: toast },
      directives: { tooltip: {} },
      stubs: {
        CFormItemList: {
          props: ['items'],
          template: `<div>
            <div v-for="item in items" :key="item.userID" data-testid="member-row">
              <slot :item="item" />
              <slot name="hover-actions" :item="item" />
            </div>
          </div>`,
        },
        CFormGroup: { template: '<div><slot /></div>' },
        Dialog: {
          props: ['visible'],
          template:
            '<div v-if="visible" data-testid="move-dialog"><slot /><slot name="footer" /></div>',
        },
        Button: {
          props: ['label', 'disabled'],
          emits: ['click'],
          template:
            '<button :disabled="disabled" @click="$emit(\'click\', $event)">{{ label }}</button>',
        },
      },
    },
  })
}

const rows = w => w.findAll('[data-testid="member-row"]').map(r => r.text())

beforeEach(() => {
  h.confirm = null
  toast = { toastSuccess: vi.fn(), toastErrorHandler: vi.fn(() => vi.fn()) }
  api = {
    userGroupMemberList: vi.fn(async () => ({ set: ['U1', 'U2'] })),
    userRead: vi.fn(async ({ userID }) => USERS[userID]),
    userGroupList: vi.fn(async () => ({ set: [ENG, SALES] })),
    userGroupRead: vi.fn(async () => SALES),
    userGroupMemberAdd: vi.fn(async () => ({})),
  }
})

describe('UserGroupMembers', () => {
  it('says each user belongs to one group', async () => {
    const w = mountMembers()
    await flushPromises()

    expect(w.text()).toContain('system.user-groups.editor.members.explainer')
  })

  it('moves a member through one dialog that names both groups', async () => {
    const w = mountMembers()
    await flushPromises()
    expect(rows(w)).toHaveLength(2)

    await w.findAll('[data-testid="member-move"]')[0].trigger('click')
    const picker = w.findComponent({ name: 'CInputUserGroup' })
    expect(picker.props('excludeUserGroups')).toEqual(['G1'])
    expect(w.find('[data-testid="member-move-submit"]').attributes('disabled')).toBeDefined()

    picker.vm.$emit('update:modelValue', 'G2')
    picker.vm.$emit('select', SALES)
    await flushPromises()

    expect(w.find('[data-testid="member-move-confirm"]').text()).toBe(
      'system.user-groups.editor.members.move.confirm {"user":"Jane Doe","from":"Engineering","to":"Sales"}',
    )

    await w.find('[data-testid="member-move-submit"]').trigger('click')
    await flushPromises()

    expect(api.userGroupMemberAdd).toHaveBeenCalledWith({ userGroupID: 'G2', userID: 'U1' })
    expect(rows(w)).toHaveLength(1)
    expect(rows(w)[0]).toContain('John Roe')
    expect(w.find('[data-testid="move-dialog"]').exists()).toBe(false)
  })

  it('keeps the member when the move fails', async () => {
    api.userGroupMemberAdd = vi.fn(async () => {
      throw new Error('nope')
    })
    const w = mountMembers()
    await flushPromises()

    await w.findAll('[data-testid="member-move"]')[0].trigger('click')
    const picker = w.findComponent({ name: 'CInputUserGroup' })
    picker.vm.$emit('update:modelValue', 'G2')
    picker.vm.$emit('select', SALES)
    await flushPromises()
    await w.find('[data-testid="member-move-submit"]').trigger('click')
    await flushPromises()

    expect(rows(w)).toHaveLength(2)
    expect(toast.toastErrorHandler).toHaveBeenCalled()
  })

  it('leaves users already in the group out of the add picker', async () => {
    const w = mountMembers()
    await flushPromises()

    expect(w.findComponent({ name: 'CInputUser' }).props('excludeUsers')).toEqual(['U1', 'U2'])
  })

  it('asks before adding a user who is in another group, and adds on yes', async () => {
    const w = mountMembers()
    await flushPromises()

    w.findComponent({ name: 'CInputUser' }).vm.$emit('select', {
      userID: 'U3',
      name: 'Max Poe',
      userGroupID: 'G2',
    })
    await flushPromises()

    expect(api.userGroupMemberAdd).not.toHaveBeenCalled()
    expect(h.confirm.message).toBe(
      'system.user-groups.editor.members.move.confirm {"user":"Max Poe","from":"Sales","to":"Engineering"}',
    )

    await h.confirm.accept()
    await flushPromises()

    expect(api.userGroupMemberAdd).toHaveBeenCalledWith({ userGroupID: 'G1', userID: 'U3' })
    expect(rows(w)).toHaveLength(3)
  })

  it('adds a user with no group without asking', async () => {
    const w = mountMembers()
    await flushPromises()

    w.findComponent({ name: 'CInputUser' }).vm.$emit('select', {
      userID: 'U4',
      name: 'Ann Loe',
      userGroupID: '0',
    })
    await flushPromises()

    expect(h.confirm).toBeNull()
    expect(api.userGroupMemberAdd).toHaveBeenCalledWith({ userGroupID: 'G1', userID: 'U4' })
  })

  it('offers neither add nor move without manage-members on the group', async () => {
    const w = mountMembers({ ...ENG, canManageMembersOnUserGroup: false })
    await flushPromises()

    expect(rows(w)).toHaveLength(2)
    expect(w.findComponent({ name: 'CInputUser' }).exists()).toBe(false)
    expect(w.find('[data-testid="member-move"]').exists()).toBe(false)
  })
})
