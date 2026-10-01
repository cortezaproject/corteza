import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { compose } from '@planetcrust/human-js'
import Tag from 'primevue/tag'

// Each row's ⋮ menu is built from the page's and the namespace's permission
// flags, a row click goes to the builder only for a page the viewer may
// update, a page is carried by its grip and saved once on release as its new
// parent plus the order of the level it landed in. These assert the gating,
// the routes and the saves; where the page lands is pageTreeDrop's, tested
// on its own.

const push = vi.fn()
const openPermissions = vi.fn()
const confirmDelete = vi.fn()
const pageDelete = vi.fn()
const pageUpdate = vi.fn()
const pageReorder = vi.fn()

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
  usePageStore: () => ({
    loadTree: async () => tree,
    delete: pageDelete,
    update: pageUpdate,
    reorder: pageReorder,
    load: vi.fn(),
  }),
  components: {
    CViewContainer: { name: 'CViewContainer', template: '<div><slot /></div>' },
    CInputSearch: { template: '<div />' },
    CPermissionsButton: { template: '<div />' },
    CRouterLinkButton: { template: '<div />' },
  },
}))

import List from './List.vue'

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
      components: {
        Menu,
        Tag,
        Dialog: {
          props: ['visible'],
          template: '<div v-if="visible"><slot /><slot name="footer" /></div>',
        },
        Select: { props: ['modelValue', 'options'], template: '<div />' },
        CFormGroup: { template: '<div><slot /></div>' },
      },
      stubs: {
        Teleport: true,
        Card: { template: '<div><slot name="header" /><slot name="content" /></div>' },
        Button: true,
      },
      directives: { tooltip: {}, ripple: {} },
      mocks: { $t: (k, n) => (n === undefined ? k : `${k}:${n}`) },
      provide: {
        $toast: { toastSuccess: vi.fn(), toastDanger: vi.fn(), toastErrorHandler: () => vi.fn() },
      },
    },
    attachTo: document.body,
  })
  await flushPromises()
  return wrapper
}

const nodeFor = pageID => {
  const walk = nodes => {
    for (const n of nodes) {
      if (n.key === pageID) return n
      const found = walk(n.children)
      if (found) return found
    }
  }
  return walk(wrapper.vm.$.setupState.treeNodes)
}
const rowOf = pageID => wrapper.find(`[data-test-id=page-tree-node][data-key="${pageID}"]`)
const titles = () => wrapper.findAll('[data-test-id=page-tree-title]').map(r => r.text())
const menuFor = pageID => wrapper.vm.$.setupState.actionItems(nodeFor(pageID).data)
const labels = items => items.filter(i => !i.separator).map(i => i.label)

const pointer = (type, x, y, target = document) =>
  target.dispatchEvent(
    new MouseEvent(type, { bubbles: true, clientX: x, clientY: y, button: 0, cancelable: true }),
  )

beforeEach(() => {
  push.mockReset()
  openPermissions.mockReset()
  confirmDelete.mockReset()
  pageUpdate.mockReset()
  pageReorder.mockReset()
})

afterEach(() => wrapper?.unmount())

describe('page tree actions', () => {
  it('offers everything to an admin, with both delete strategies on a parent', async () => {
    tree = [page('101', { children: [page('102', { selfID: '101' })] })]
    await mountList()

    // no "make sub-page of": nothing here to go under but its own sub-page
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

  it('opens record pages on a new record', async () => {
    tree = [page('201', { moduleID: '301', visible: false })]
    await mountList()

    const items = menuFor('201')
    expect(items.find(i => i.label === 'page.view').route).toEqual({
      name: 'page.record',
      params: { pageID: '201', recordID: '0' },
    })
    expect(labels(items)).toContain('page.list.addSubPage')
    expect(wrapper.vm.$.setupState.recordPageLabel(new compose.Page(tree[0]))).toBe(
      'page.list.recordPageOf:{"module":"Leads"}',
    )
    expect(wrapper.findComponent(Tag).text()).toBe('page.list.recordPageOf:{"module":"Leads"}')
  })

  it('marks a page missing from navigation with an eye-slash tag', async () => {
    tree = [page('101', { visible: false })]
    await mountList()

    const tag = wrapper.findComponent(Tag)
    expect(tag.text()).toContain('page.notVisible')
    expect(tag.find('span.pi-eye-slash').exists()).toBe(true)
  })

  it('renders every level with its key, parent and depth on the row', async () => {
    tree = [page('101', { children: [page('102', { selfID: '101' })] }), page('103')]
    await mountList()

    expect(titles()).toEqual(['101', '102', '103'])
    expect(
      [rowOf('101'), rowOf('102'), rowOf('103')].map(r => [
        r.attributes('data-parent'),
        r.attributes('data-depth'),
      ]),
    ).toEqual([
      ['0', '0'],
      ['101', '1'],
      ['0', '0'],
    ])
  })

  it('gives every page a grip except the single top-level page', async () => {
    tree = [page('101', { children: [page('102', { selfID: '101' })] })]
    await mountList()
    const gripOf = id => rowOf(id).find('.page-grip').exists()
    expect([gripOf('101'), gripOf('102')]).toEqual([false, true])

    wrapper.unmount()
    tree = [page('101'), page('103')]
    await mountList()
    expect([gripOf('101'), gripOf('103')]).toEqual([true, true])
  })

  it('carries a page once the pointer has travelled from its grip, and lets go on Escape', async () => {
    tree = [page('101'), page('102')]
    await mountList()
    const grip = rowOf('102').find('.page-grip').element
    const carried = () => wrapper.vm.$.setupState.carriedKey

    pointer('pointerdown', 100, 100, grip)
    pointer('pointermove', 102, 102)
    expect(carried()).toBe(null)

    pointer('pointermove', 100, 120)
    await flushPromises()
    expect(carried()).toBe('102')
    expect(rowOf('102').element.closest('li').classList.contains('page-node-carried')).toBe(true)
    expect(document.body.classList.contains('c-dragging')).toBe(true)

    document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }))
    await flushPromises()
    expect(carried()).toBe(null)
    expect(document.body.classList.contains('c-dragging')).toBe(false)
    expect(pageUpdate).not.toHaveBeenCalled()
  })

  it('carries a parent with its subtree listed on the card', async () => {
    tree = [
      page('101', {
        title: 'Sales',
        children: [
          page('102', {
            title: 'Accounts',
            selfID: '101',
            children: [page('103', { title: 'Lead', selfID: '102' })],
          }),
          page('104', { title: 'Contacts', selfID: '101' }),
          page('105', { title: 'Deals', selfID: '101' }),
          page('106', { title: 'Notes', selfID: '101' }),
          page('107', { title: 'Extra', selfID: '101' }),
        ],
      }),
      page('108', { title: 'Reports' }),
    ]
    await mountList()
    pointer('pointerdown', 100, 100, rowOf('101').find('.page-grip').element)
    pointer('pointermove', 100, 120)
    await flushPromises()

    const card = wrapper.find('[data-test-id=page-carry]')
    expect(card.text()).toContain('page.list.subPages:5')
    expect(card.findAll('.page-carry-kids li').map(li => li.text())).toEqual([
      'Accounts',
      'Lead',
      'Contacts',
      'Deals',
      'page.list.andMore:2',
    ])
    document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }))
  })

  it('lands the card on release and keeps it until the tree has caught up', async () => {
    // jsdom draws every row at 0,0: the first page released below the second
    // lands after it, a real change to save
    tree = [page('101'), page('102')]
    await mountList()
    const grip = rowOf('101').find('.page-grip').element
    const state = wrapper.vm.$.setupState

    pointer('pointerdown', 100, 100, grip)
    pointer('pointermove', 100, 120)
    await flushPromises()
    expect(state.drag.landing).toBe(false)

    let settle
    pageReorder.mockImplementationOnce(() => new Promise(r => (settle = r)))
    pointer('pointerup', 100, 120)
    expect(state.drag?.landing).toBe(true)
    expect(document.body.classList.contains('c-dragging')).toBe(false)

    settle()
    await flushPromises()
    expect(state.drag).toBe(null)
    expect(pageReorder).toHaveBeenCalledWith({
      namespaceID: '1',
      selfID: '0',
      pageIDs: ['102', '101'],
    })
  })

  it('saves a drop into another level as the new parent, then the order of that level', async () => {
    tree = [page('101', { children: [page('102', { selfID: '101' })] }), page('103')]
    await mountList()

    await wrapper.vm.$.setupState.applyDrop(nodeFor('103'), { parentKey: '101', afterKey: '102' })

    expect(pageUpdate).toHaveBeenCalledTimes(1)
    expect(pageUpdate.mock.calls[0][0]).toMatchObject({
      pageID: '103',
      selfID: '101',
      namespaceID: '1',
    })
    expect(pageReorder).toHaveBeenCalledWith({
      namespaceID: '1',
      selfID: '101',
      pageIDs: ['102', '103'],
    })
  })

  it('saves a drop within a level as a reorder alone, and nothing for a drop in place', async () => {
    tree = [page('101'), page('103')]
    await mountList()
    const { applyDrop } = wrapper.vm.$.setupState

    await applyDrop(nodeFor('103'), { parentKey: '0', afterKey: null })
    expect(pageUpdate).not.toHaveBeenCalled()
    expect(pageReorder).toHaveBeenCalledWith({
      namespaceID: '1',
      selfID: '0',
      pageIDs: ['103', '101'],
    })

    pageReorder.mockClear()
    await applyDrop(nodeFor('103'), { parentKey: '0', afterKey: '101' })
    expect(pageReorder).not.toHaveBeenCalled()
  })

  it('offers "make sub-page of" with the parents this page may go under, and moves it last among them', async () => {
    tree = [
      page('101', {
        title: 'Sales',
        children: [
          page('102', {
            title: 'Accounts',
            selfID: '101',
            children: [page('103', { title: 'Lead', selfID: '102' })],
          }),
        ],
      }),
      page('104', { title: 'Reports' }),
      page('201', { title: 'Record', moduleID: '301' }),
    ]
    await mountList()
    const { openMoveUnder, confirmMoveUnder, moveTargets } = wrapper.vm.$.setupState

    expect(moveTargets(nodeFor('104').data).map(o => o.value)).toEqual(['101', '102', '103', '201'])
    expect(moveTargets(nodeFor('103').data).map(o => o.value)).toEqual(['0', '101', '104', '201'])
    expect(labels(menuFor('104'))).toContain('page.list.makeSubPageOf')

    openMoveUnder(nodeFor('104').data)
    wrapper.vm.$.setupState.moveUnder.targetID = '102'
    await confirmMoveUnder()

    expect(pageUpdate.mock.calls[0][0]).toMatchObject({ pageID: '104', selfID: '102' })
    expect(pageReorder).toHaveBeenCalledWith({
      namespaceID: '1',
      selfID: '102',
      pageIDs: ['103', '104'],
    })
  })

  it('filters to matches, keeping the path down to them and what sits beneath, with no grips', async () => {
    tree = [
      page('101', {
        title: 'Sales',
        children: [
          page('102', {
            title: 'Accounts',
            selfID: '101',
            children: [page('103', { title: 'Lead' })],
          }),
          page('104', { title: 'Contacts', selfID: '101' }),
        ],
      }),
      page('105', { title: 'Reports' }),
    ]
    await mountList()
    expect(wrapper.findAll('.page-grip').length).toBe(5)

    wrapper.vm.$.setupState.filterValue = ' acc '
    await flushPromises()

    expect(titles()).toEqual(['Sales', 'Accounts', 'Lead'])
    expect(wrapper.findAll('.page-grip')).toHaveLength(0)

    wrapper.vm.$.setupState.filterValue = 'nothing like it'
    await flushPromises()
    expect(titles()).toEqual([])
    expect(wrapper.text()).toContain('general.label.noResults')
  })

  it('shows the server order again when a move fails to save', async () => {
    tree = [page('101', { title: 'First' }), page('102', { title: 'Second' })]
    await mountList()
    pageReorder.mockRejectedValueOnce(new Error('denied'))

    await wrapper.vm.$.setupState.applyDrop(nodeFor('102'), { parentKey: '0', afterKey: null })

    expect(titles()).toEqual(['First', 'Second'])
  })

  it('hides what the viewer may not do', async () => {
    tree = [page('101', { canUpdatePage: false, canDeletePage: false, canGrant: false })]
    await mountList({ namespaceID: '1', canCreatePage: false, canGrant: false })

    expect(labels(menuFor('101'))).toEqual(['page.view'])
    expect(wrapper.findAll('.page-grip')).toHaveLength(0)
  })

  it('opens the builder on click, and the page itself without update rights', async () => {
    tree = [page('101'), page('102', { canUpdatePage: false })]
    await mountList()

    await rowOf('101').trigger('click')
    await rowOf('102').trigger('click')

    expect(push.mock.calls).toEqual([
      [{ name: 'admin.pages.builder', params: { pageID: '101' } }],
      [{ name: 'page', params: { pageID: '102' } }],
    ])
  })
})
