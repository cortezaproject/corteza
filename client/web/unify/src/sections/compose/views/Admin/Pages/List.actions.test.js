import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { compose } from '@planetcrust/human-js'

// Each tree node's ⋮ menu is built from the page's and the namespace's
// permission flags, and a row click goes to the builder only for a page the
// viewer may update. These assert the gating and the routes.

const push = vi.fn()
const openPermissions = vi.fn()
const confirmDelete = vi.fn()
const pageDelete = vi.fn()

vi.mock('vue-router', () => ({
  useRouter: () => ({ push }),
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: (k, p) => (p ? `${k}:${JSON.stringify(p)}` : k) }),
}))

let tree = []

vi.mock('@planetcrust/human-vue', () => ({
  useConfirmDelete: () => ({ confirmDelete }),
  usePermissions: () => ({ open: openPermissions }),
  useModuleStore: () => ({ getByID: id => (id === '301' ? { name: 'Leads' } : undefined) }),
  usePageStore: () => ({ loadTree: async () => tree, delete: pageDelete }),
  components: {
    CInputSearch: { template: '<div />' },
    CPermissionsButton: { template: '<div />' },
    CRouterLinkButton: { template: '<div />' },
  },
}))

import List from './List.vue'

const Tree = {
  name: 'Tree',
  props: { value: { type: Array, default: () => [] } },
  emits: ['node-select', 'node-drop'],
  template: `<div><template v-for="n in value" :key="n.key"><slot :node="n" /></template></div>`,
}

const Menu = {
  name: 'Menu',
  props: { model: { type: Array, default: () => [] } },
  methods: { show() {} },
  template: '<div />',
}

const admin = { namespaceID: '1', canCreatePage: true, canGrant: true }
const page = (id, extra = {}) => ({
  pageID: id,
  title: id,
  selfID: '0',
  moduleID: '0',
  visible: true,
  canUpdatePage: true,
  canDeletePage: true,
  canGrant: true,
  children: [],
  ...extra,
})

let wrapper

async function mountList(namespace = admin) {
  wrapper = mount(List, {
    props: { namespace },
    global: {
      components: { Tree, Menu },
      stubs: {
        Teleport: true,
        Card: { template: '<div><slot name="header" /><slot name="content" /></div>' },
        Button: true,
        Tag: true,
      },
      directives: { tooltip: {}, ripple: {} },
      mocks: { $t: k => k },
      provide: {
        $toast: { toastSuccess: vi.fn(), toastDanger: vi.fn(), toastErrorHandler: () => vi.fn() },
      },
    },
  })
  await flushPromises()
  return wrapper
}

function menuFor(pageID) {
  const node = wrapper
    .findComponent(Tree)
    .props('value')
    .find(n => n.key === pageID)
  return wrapper.vm.$.setupState.actionItems(node.data)
}

const labels = items => items.filter(i => !i.separator).map(i => i.label)

beforeEach(() => {
  push.mockReset()
  openPermissions.mockReset()
  confirmDelete.mockReset()
})

afterEach(() => wrapper?.unmount())

describe('page tree actions', () => {
  it('offers everything to an admin, with both delete strategies on a parent', async () => {
    tree = [page('101', { children: [page('102', { selfID: '101' })] })]
    await mountList()

    expect(labels(menuFor('101'))).toEqual([
      'general.label.pageBuilder',
      'page.list.settings',
      'page.view',
      'page.list.addSubPage',
      'page.list.pagePermissions',
      'page.list.layoutPermissions',
      'page.list.deleteKeepSubPages',
      'page.list.deleteWithSubPages',
    ])
  })

  it('opens page and page-layout permissions on their exact resources', async () => {
    tree = [page('101')]
    await mountList()

    const items = menuFor('101')
    items.find(i => i.label === 'page.list.pagePermissions').command()
    items.find(i => i.label === 'page.list.layoutPermissions').command()

    expect(openPermissions.mock.calls.map(([o]) => [o.resource, !!o.allSpecific])).toEqual([
      ['corteza::compose:page/1/101', false],
      ['corteza::compose:page-layout/1/101/*', true],
    ])
  })

  it('routes add sub-page to the create form with the parent', async () => {
    tree = [page('101')]
    await mountList()

    const add = menuFor('101').find(i => i.label === 'page.list.addSubPage')
    expect(add.route).toEqual({ name: 'admin.pages.create', query: { parent: '101' } })
  })

  it('gives a leaf a single delete that aborts on children', async () => {
    tree = [page('101')]
    await mountList()

    const del = menuFor('101').filter(i => i.class === 'text-red-500')
    expect(del.map(i => i.label)).toEqual(['general.label.delete'])
    del[0].command()
    confirmDelete.mock.calls[0][0].onConfirm()
    await flushPromises()
    expect(pageDelete).toHaveBeenCalledWith({ namespaceID: '1', pageID: '101', strategy: 'abort' })
  })

  it('opens record pages on a new record and never nests under them', async () => {
    tree = [page('201', { moduleID: '301', visible: false })]
    await mountList()

    const items = menuFor('201')
    expect(items.find(i => i.label === 'page.view').route).toEqual({
      name: 'page.record',
      params: { pageID: '201', recordID: '0' },
    })
    expect(labels(items)).not.toContain('page.list.addSubPage')
    expect(wrapper.vm.$.setupState.recordPageLabel(new compose.Page(tree[0]))).toBe(
      'page.list.recordPageOf:{"module":"Leads"}',
    )
  })

  it('hides what the viewer may not do', async () => {
    tree = [page('101', { canUpdatePage: false, canDeletePage: false, canGrant: false })]
    await mountList({ namespaceID: '1', canCreatePage: false, canGrant: false })

    expect(labels(menuFor('101'))).toEqual(['page.view'])
  })

  it('opens the builder on click, and the page itself without update rights', async () => {
    tree = [page('101'), page('102', { canUpdatePage: false })]
    await mountList()
    const nodes = wrapper.findComponent(Tree).props('value')

    wrapper.findComponent(Tree).vm.$emit('node-select', nodes[0])
    wrapper.findComponent(Tree).vm.$emit('node-select', nodes[1])

    expect(push.mock.calls).toEqual([
      [{ name: 'admin.pages.builder', params: { pageID: '101' } }],
      [{ name: 'page', params: { pageID: '102' } }],
    ])
  })
})
