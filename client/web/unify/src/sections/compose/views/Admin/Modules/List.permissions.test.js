import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { ref } from 'vue'

// A compose module's RBAC resource is `module/<namespaceID>/<moduleID>`, and
// the server rejects any other shape outright. The permission dialog swallows
// that rejection and draws an empty rule list, so a malformed resource here
// costs a dead button and no error anywhere — these assert the exact strings.

const openPermissions = vi.fn()

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: vi.fn() }),
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: k => k }),
}))

const Menu = {
  name: 'Menu',
  props: {
    model: { type: Array, default: () => [] },
  },
  template: '<div />',
}

vi.mock('@planetcrust/human-vue', () => ({
  changedAtField: header => ({ key: 'changedAt', header }),
  changedAtText: () => '',
  useConfirmDelete: () => ({ confirmDelete: vi.fn() }),
  usePermissions: () => ({ open: openPermissions }),
  useResourceList: () => ({
    items: ref([]),
    loading: ref(false),
    filter: ref({ query: '' }),
    sorting: ref({}),
    pagination: ref({}),
    handleSort: vi.fn(),
    handlePageChange: vi.fn(),
    filterList: vi.fn(),
  }),
  useModuleStore: () => ({ delete: vi.fn() }),
  components: {
    CViewContainer: { name: 'CViewContainer', template: '<div><slot /></div>' },
    CResourceList: {
      name: 'CResourceList',
      props: {
        actionItems: { type: Function, default: undefined },
      },
      methods: {
        hideActionsMenu() {},
      },
      template: '<div><slot name="header" /></div>',
    },
    CRouterLinkButton: { template: '<div />' },
  },
  filters: { locFullDateTime: v => v },
}))

vi.mock('@/sections/compose/components/Modules/ModuleImporter.vue', () => ({
  default: { template: '<div />' },
}))

import { components } from '@planetcrust/human-vue'
import List from './List.vue'

const { CResourceList } = components

const namespace = { namespaceID: 'NS1', canGrant: true, canExportModules: true }
const module = { moduleID: 'M1', name: 'Leads', canGrant: true }

let wrapper

async function mountList() {
  wrapper = mount(List, {
    props: { namespace },
    global: {
      components: { Menu },
      stubs: { teleport: true, Teleport: true, Button: true, Tag: true },
      directives: { tooltip: {} },
      mocks: { $t: k => k },
      provide: {
        $toast: { toastSuccess: vi.fn(), toastErrorHandler: () => vi.fn() },
        $ComposeAPI: { moduleListCancellable: vi.fn() },
      },
    },
  })
  await flushPromises()
  return wrapper
}

function resourcesOf(items) {
  return items.map(i => {
    openPermissions.mockClear()
    i.command()
    return openPermissions.mock.calls[0][0].resource
  })
}

beforeEach(() => openPermissions.mockClear())

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
})

describe('Module list permissions', () => {
  it('scopes the header menu to every module in the namespace', async () => {
    await mountList()
    const items = wrapper.findComponent(Menu).props('model')

    expect(resourcesOf(items)).toEqual([
      'corteza::compose:module/NS1/*',
      'corteza::compose:module-field/NS1/*/*',
      'corteza::compose:record/NS1/*/*',
    ])
  })

  it('scopes the row menu to the one module', async () => {
    await mountList()
    const items = wrapper
      .findComponent(CResourceList)
      .props('actionItems')(module)
      .filter(i => i.icon === 'pi pi-lock')

    expect(resourcesOf(items)).toEqual([
      'corteza::compose:module/NS1/M1',
      'corteza::compose:module-field/NS1/M1/*',
      'corteza::compose:record/NS1/M1/*',
    ])
  })

  it('names the module in the row menu, and nothing in the header menu', async () => {
    await mountList()

    const header = wrapper.findComponent(Menu).props('model')
    header[0].command()
    expect(openPermissions.mock.calls[0][0].target).toBeUndefined()

    openPermissions.mockClear()
    const row = wrapper
      .findComponent(CResourceList)
      .props('actionItems')(module)
      .filter(i => i.icon === 'pi pi-lock')
    row[1].command()
    expect(openPermissions.mock.calls[0][0]).toMatchObject({
      title: 'Leads',
      target: 'Leads',
      allSpecific: true,
    })
  })
})
