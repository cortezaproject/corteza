import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'

// The list dropped its handle column, and the description subtitle carries a
// bounded max-width: uncapped, a long description is a nowrap run that widens
// the name cell until every other column is pushed off the right-hand edge.

const h = vi.hoisted(() => ({ capturedFields: null, items: null }))

vi.mock('@planetcrust/human-vue', async () => {
  const { ref, reactive } = await import('vue')
  h.items = ref([])
  return {
    changedAtField: header => ({ key: 'changedAt', header }),
    changedAtText: () => '',
    components: {
      CViewContainer: { name: 'CViewContainer', template: '<div><slot /></div>' },
      CResourceList: {
        name: 'CResourceList',
        props: {
          fields: { type: Array, default: () => [] },
          items: { type: Array, default: () => [] },
        },
        setup(props) {
          h.capturedFields = props.fields
        },
        template: `
          <div>
            <div v-for="(row, i) in items" :key="i">
              <slot name="body-name" :data="row" />
            </div>
          </div>`,
      },
      CRouterLinkButton: { template: '<div />' },
    },
    useAgentStore: () => ({ updateInList: vi.fn(), removeFromList: vi.fn() }),
    useConfirmDelete: () => ({ confirmDelete: vi.fn() }),
    usePermissions: () => ({ open: vi.fn() }),
    useRBACStore: () => ({ can: () => true }),
    useResourceList: () => ({
      items: h.items,
      loading: ref(false),
      filter: reactive({ query: '' }),
      sorting: reactive({ sortBy: 'name', sortDesc: false }),
      pagination: reactive({ limit: 50 }),
      handleSort: vi.fn(),
      handlePageChange: vi.fn(),
      filterList: vi.fn(),
    }),
  }
})

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))
vi.mock('vue-router', () => ({ useRouter: () => ({ push: vi.fn() }) }))

import Home from './Home.vue'

function mountHome(rows = []) {
  h.items.value = rows
  return mount(Home, {
    global: {
      mocks: { $t: k => k },
      provide: {
        $SystemAPI: { agentListCancellable: () => ({}) },
        $toast: { toastSuccess: vi.fn(), toastDanger: vi.fn(), toastErrorHandler: () => vi.fn() },
      },
      directives: { tooltip: {} },
      stubs: { Teleport: true, Tag: true, CPermissionsButton: true },
    },
  })
}

describe('agent list columns', () => {
  it('has no handle column', () => {
    mountHome()
    expect(h.capturedFields.map(f => f.key)).toEqual(['name', 'status', 'changedAt'])
  })

  it('bounds the description subtitle so it cannot widen the name cell', () => {
    const wrapper = mountHome([{ meta: { short: 'Probe', description: 'x'.repeat(400) } }])
    const subtitle = wrapper.findAll('span').find(s => s.text().startsWith('xxx'))
    expect(subtitle).toBeTruthy()
    expect(subtitle.classes()).toContain('truncate')
    expect(subtitle.classes().some(c => /^max-w-(?!full$)/.test(c))).toBe(true)
  })
})
