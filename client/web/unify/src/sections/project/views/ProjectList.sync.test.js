import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { ref, reactive } from 'vue'

// The list fetches its own rows; the section sidebar renders the shared project
// cache. A mutation here that never reaches the store leaves the nav offering a
// project that has been shelved or deleted, and calling it by its old name.

const h = vi.hoisted(() => ({
  capturedActionItems: null,
  confirmed: null,
  store: null,
}))

vi.mock('@planetcrust/human-vue', () => ({
  changedAt: r => r?.updatedAt,
  changedAtField: header => ({ key: 'changedAt', header }),
  components: {
    CResourceList: {
      name: 'CResourceList',
      props: { actionItems: { type: Function, default: null } },
      setup(props) {
        h.capturedActionItems = props.actionItems
      },
      template: '<div />',
    },
    CViewContainer: { name: 'CViewContainer', template: '<div><slot /></div>' },
  },
  useRBACStore: () => ({ can: () => true }),
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

vi.mock('@planetcrust/human-js', () => ({ system: { Project: class Project {} } }))
vi.mock('vue-router', () => ({ useRouter: () => ({ push: vi.fn() }) }))
vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))
vi.mock('primevue/useconfirm', () => ({
  useConfirm: () => ({
    require: opts => {
      h.confirmed = opts
    },
  }),
}))
vi.mock('@/sections/project/config/publishState', () => ({ chainHasPublished: () => false }))
vi.mock('@/sections/project/stores/projects', () => ({ useProjectsStore: () => h.store }))

vi.mock('@/sections/project/components/project/NewProjectDialog.vue', () => ({
  default: { template: '<div />' },
}))
vi.mock('@/sections/project/components/project/RenameProjectDialog.vue', () => ({
  default: { name: 'RenameProjectDialog', emits: ['rename'], template: '<div />' },
}))
vi.mock('@/sections/project/components/project/StatusChip.vue', () => ({
  default: { template: '<div />' },
}))

import ProjectList from './ProjectList.vue'

const PROJECT = {
  projectID: '42',
  name: 'Ledger',
  handle: 'ledger',
  status: 'draft',
  config: {},
  meta: { short: 'Ledger' },
  labels: {},
  updatedAt: '2026-01-01T00:00:00Z',
  canUpdateProject: true,
  canDeleteProject: true,
}

const ARCHIVED = { ...PROJECT, projectID: '42', archivedAt: '2026-02-01T00:00:00Z' }

let api

function mountList() {
  return mount(ProjectList, {
    global: {
      mocks: { $t: k => k },
      provide: {
        $toast: { toastSuccess: vi.fn(), toastErrorHandler: () => vi.fn() },
        $SystemAPI: api,
      },
      stubs: {
        Teleport: true,
        Button: { props: { label: { type: String, default: '' } }, template: '<button />' },
        Popover: true,
        RadioButton: true,
      },
    },
  })
}

const commandFor = (project, label) =>
  h.capturedActionItems(project).find(i => i.label === label).command

beforeEach(() => {
  h.capturedActionItems = null
  h.confirmed = null
  h.store = { absorb: vi.fn(), removeProject: vi.fn().mockResolvedValue(undefined) }
  api = {
    projectUpdate: vi.fn(raw => Promise.resolve({ ...raw, meta: { short: 'Renamed' } })),
    projectArchive: vi.fn(() => Promise.resolve({ ...PROJECT, archivedAt: 'now' })),
    projectUnarchive: vi.fn(() => Promise.resolve({ ...PROJECT, archivedAt: null })),
    projectDelete: vi.fn(),
  }
})

describe('ProjectList mutations reach the shared project cache', () => {
  it('absorbs the project an archive gives back', async () => {
    mountList()
    commandFor(PROJECT, 'project.list.actions.archive')()
    await flushPromises()

    expect(api.projectArchive).toHaveBeenCalledWith({ projectID: '42' })
    expect(h.store.absorb).toHaveBeenCalledWith(expect.objectContaining({ archivedAt: 'now' }))
  })

  it('absorbs the project an unarchive gives back', async () => {
    mountList()
    commandFor(ARCHIVED, 'project.list.actions.unarchive')()
    await flushPromises()

    expect(api.projectUnarchive).toHaveBeenCalledWith({ projectID: '42' })
    expect(h.store.absorb).toHaveBeenCalledWith(expect.objectContaining({ archivedAt: null }))
  })

  it('deletes through the store, so the cache drops the row with it', async () => {
    mountList()
    commandFor(PROJECT, 'general.label.delete')()
    h.confirmed.accept()
    await flushPromises()

    expect(h.store.removeProject).toHaveBeenCalledWith('42')
    // The store owns the call; a second one here would delete twice.
    expect(api.projectDelete).not.toHaveBeenCalled()
  })

  it('absorbs the project a rename gives back', async () => {
    const wrapper = mountList()
    // The menu entry only arms the dialog; the rename lands when it emits.
    commandFor(PROJECT, 'project.list.actions.rename')()
    await flushPromises()
    wrapper.findComponent({ name: 'RenameProjectDialog' }).vm.$emit('rename', 'Renamed')
    await flushPromises()

    expect(api.projectUpdate).toHaveBeenCalledWith(
      expect.objectContaining({ projectID: '42', meta: { short: 'Renamed' } }),
    )
    expect(h.store.absorb).toHaveBeenCalledWith(
      expect.objectContaining({ projectID: '42', meta: { short: 'Renamed' } }),
    )
  })
})
