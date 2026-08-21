import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { ref, reactive } from 'vue'

// A project row offers only what the backend would accept: rename and archive
// are both project updates, delete is its own permission, and creating one is
// a component-level operation with no row to carry a flag. Every one of these
// was offered unconditionally, so a read-only user got four dead buttons.

// vi.mock factories are hoisted above the file body, so anything they close
// over has to be hoisted with them.
const h = vi.hoisted(() => ({
  capturedActionItems: null,
  rbacRules: { 'system//project.create': true },
}))

vi.mock('@planetcrust/human-vue', () => ({
  changedAt: r => r?.deletedAt || r?.updatedAt || r?.createdAt,
  changedAtField: header => ({ key: 'changedAt', header }),
  components: {
    CResourceList: {
      name: 'CResourceList',
      props: {
        actionItems: { type: Function, default: null },
        items: { type: Array, default: () => [] },
      },
      setup(props) {
        h.capturedActionItems = props.actionItems
      },
      template: '<div><slot name="header" /></div>',
    },
    CViewContainer: { name: 'CViewContainer', template: '<div><slot /></div>' },
  },
  useRBACStore: () => ({
    can: (resource, op) => !!h.rbacRules[`${resource}/${op}`],
  }),
  useResourceList: () => ({
    items: ref([]),
    loading: ref(false),
    filter: reactive({ query: '', status: '' }),
    sorting: reactive({ sortBy: 'name', sortDesc: false }),
    pagination: reactive({ limit: 50 }),
    handleSort: vi.fn(),
    handlePageChange: vi.fn(),
    filterList: vi.fn(),
  }),
}))

vi.mock('@planetcrust/human-js', () => ({
  system: { Project: class Project {} },
}))

vi.mock('vue-router', () => ({ useRouter: () => ({ push: vi.fn() }) }))
vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))
vi.mock('primevue/useconfirm', () => ({ useConfirm: () => ({ require: vi.fn() }) }))
vi.mock('@/sections/project/config/publishState', () => ({ chainHasPublished: () => false }))

vi.mock('@/sections/project/components/project/NewProjectDialog.vue', () => ({
  default: { template: '<div />' },
}))
vi.mock('@/sections/project/components/project/RenameProjectDialog.vue', () => ({
  default: { template: '<div />' },
}))
vi.mock('@/sections/project/components/project/StatusChip.vue', () => ({
  default: { template: '<div />' },
}))

import ProjectList from './ProjectList.vue'

function mountList() {
  return mount(ProjectList, {
    global: {
      mocks: { $t: k => k },
      provide: {
        $toast: { toastSuccess: vi.fn(), toastErrorHandler: () => vi.fn() },
        $SystemAPI: {},
      },
      stubs: {
        Teleport: true,
        // Typed props — a valueless attribute reaches an array-props stub as ''.
        Button: { props: { label: { type: String, default: '' } }, template: '<button />' },
        Popover: true,
        RadioButton: true,
      },
    },
  })
}

// Labels the menu builder emits, so a test names the action rather than an index.
const RENAME = 'project.list.actions.rename'
const ARCHIVE = 'project.list.actions.archive'
const DELETE = 'general.label.delete'

const labelsFor = project =>
  (h.capturedActionItems(project) || []).map(i => i.label).filter(Boolean)

describe('ProjectList row actions', () => {
  beforeEach(() => {
    h.capturedActionItems = null
  })

  it('offers nothing on a project the user may neither update nor delete', () => {
    mountList()
    expect(h.capturedActionItems).toBeTypeOf('function')
    // An empty array is the contract: CResourceList hides the whole kebab
    // trigger on `actionItems(row).length`, so the menu itself disappears.
    expect(h.capturedActionItems({ canUpdateProject: false, canDeleteProject: false })).toEqual([])
  })

  it('offers rename and archive only with update permission', () => {
    mountList()
    const labels = labelsFor({ canUpdateProject: true, canDeleteProject: false })
    expect(labels).toContain(RENAME)
    expect(labels).toContain(ARCHIVE)
    expect(labels).not.toContain(DELETE)
  })

  it('offers delete only with delete permission', () => {
    mountList()
    const labels = labelsFor({ canUpdateProject: false, canDeleteProject: true })
    expect(labels).toEqual([DELETE])
  })

  it('offers every action when both permissions are held', () => {
    mountList()
    const labels = labelsFor({ canUpdateProject: true, canDeleteProject: true })
    expect(labels).toEqual([RENAME, ARCHIVE, DELETE])
  })

  it('labels the archive entry unarchive on an already-archived project', () => {
    mountList()
    const labels = labelsFor({ canUpdateProject: true, archivedAt: '2026-01-01T00:00:00Z' })
    expect(labels).toContain('project.list.actions.unarchive')
    expect(labels).not.toContain(ARCHIVE)
  })
})

describe('ProjectList create button', () => {
  it('is shown when the effective rules allow project.create', () => {
    h.rbacRules['system//project.create'] = true
    const wrapper = mountList()
    expect(wrapper.find('button').exists()).toBe(true)
  })

  it('is hidden when they do not', () => {
    h.rbacRules['system//project.create'] = false
    const wrapper = mountList()
    expect(wrapper.find('button').exists()).toBe(false)
    h.rbacRules['system//project.create'] = true
  })
})
