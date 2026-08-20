import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'

// The workflow section's own list offered enable/disable, delete and undelete
// on every row, while the admin section's list of the same resource gated all
// three. These assert the two now agree with each other and with the backend.

const h = vi.hoisted(() => ({ capturedActionItems: null }))

vi.mock('@planetcrust/human-vue', async () => {
  const { ref, reactive } = await import('vue')
  return {
    components: {
      CResourceList: {
        name: 'CResourceList',
        props: { actionItems: { type: Function, default: null } },
        setup(props) {
          h.capturedActionItems = props.actionItems
        },
        template: '<div />',
      },
      CRouterLinkButton: { template: '<div />' },
    },
    filters: { locFullDateTime: () => '' },
    useConfirmDelete: () => ({ confirmDelete: vi.fn() }),
    useRBACStore: () => ({ can: () => true }),
    useWorkflowStore: () => ({ set: [], load: vi.fn() }),
    useResourceList: () => ({
      items: ref([]),
      loading: ref(false),
      filter: reactive({ query: '', subWorkflow: '1', disabled: '1', deleted: '0' }),
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
vi.mock('primevue/usetoast', () => ({ useToast: () => ({ add: vi.fn() }) }))
vi.mock('file-saver', () => ({ saveAs: vi.fn() }))
vi.mock('../stores/labels', () => ({
  useLabelsStore: () => ({ namespaceLabels: [], moduleLabels: [], load: vi.fn() }),
}))
vi.mock('../components/Import.vue', () => ({ default: { template: '<div />' } }))
vi.mock('../components/Export.vue', () => ({ default: { template: '<div />' } }))
vi.mock('../components/NamespaceModuleSelector.vue', () => ({ default: { template: '<div />' } }))

import Home from './Home.vue'

function mountHome() {
  return mount(Home, {
    global: {
      mocks: { $t: k => k },
      provide: {
        $AutomationAPI: { workflowListCancellable: () => ({}) },
        $ComposeAPI: {},
        $Auth: { user: { userID: 'U1' } },
      },
      directives: { tooltip: {} },
      stubs: {
        Teleport: true,
        Button: true,
        Popover: true,
        RadioButton: true,
        InputText: true,
        CPermissionsButton: true,
        Tag: true,
        Divider: true,
      },
    },
  })
}

const EXPORT = 'general.export'
const DELETE = 'general.label.delete'
const RESTORE = 'general.label.restore'
const DISABLE = 'general.disable'
const ENABLE = 'general.enable'

function labelsFor(workflow) {
  mountHome()
  return (h.capturedActionItems(workflow) || []).map(i => i.label).filter(Boolean)
}

describe('workflow list row actions', () => {
  it('hides enable/disable without update permission', () => {
    const labels = labelsFor({ workflowID: 'W1', enabled: true, canUpdateWorkflow: false })
    expect(labels).not.toContain(DISABLE)
    expect(labels).not.toContain(ENABLE)
  })

  it('offers disable on an enabled workflow with update permission', () => {
    const labels = labelsFor({ workflowID: 'W1', enabled: true, canUpdateWorkflow: true })
    expect(labels).toContain(DISABLE)
  })

  it('offers enable on a disabled workflow with update permission', () => {
    const labels = labelsFor({ workflowID: 'W1', enabled: false, canUpdateWorkflow: true })
    expect(labels).toContain(ENABLE)
  })

  it('hides delete without delete permission', () => {
    const labels = labelsFor({ workflowID: 'W1', canDeleteWorkflow: false })
    expect(labels).not.toContain(DELETE)
  })

  it('offers delete with delete permission', () => {
    const labels = labelsFor({ workflowID: 'W1', canDeleteWorkflow: true })
    expect(labels).toContain(DELETE)
  })

  it('hides restore on a deleted workflow without undelete permission', () => {
    const labels = labelsFor({
      workflowID: 'W1',
      deletedAt: '2026-01-01T00:00:00Z',
      canUndeleteWorkflow: false,
    })
    expect(labels).not.toContain(RESTORE)
  })

  it('offers restore on a deleted workflow with undelete permission', () => {
    const labels = labelsFor({
      workflowID: 'W1',
      deletedAt: '2026-01-01T00:00:00Z',
      canUndeleteWorkflow: true,
    })
    expect(labels).toContain(RESTORE)
  })

  it('always keeps edit and export, which need no write permission', () => {
    const labels = labelsFor({ workflowID: 'W1' })
    expect(labels).toContain('general.label.edit')
    expect(labels).toContain(EXPORT)
  })
})
