import { describe, it, expect, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'

// Saving a workflow used to re-send every trigger, so a user allowed to update
// the workflow but not to manage its triggers could not save step changes.
// Only added or modified triggers may reach the API.

vi.mock('@planetcrust/human-vue', () => ({
  useRBACStore: () => ({ can: () => true }),
  useUnsavedGuard: () => {},
  useWorkflowStore: () => ({ updateInList: vi.fn(), removeFromList: vi.fn() }),
}))
vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))
vi.mock('vue-router', () => ({
  useRoute: () => ({ params: { workflowID: 'W1' } }),
  useRouter: () => ({ push: vi.fn() }),
}))
vi.mock('primevue/usetoast', () => ({ useToast: () => ({ add: vi.fn() }) }))
vi.mock('../components/WorkflowEditor.vue', () => ({
  default: { name: 'WorkflowEditor', emits: ['save'], template: '<div />' },
}))

import Editor from './Editor.vue'

const storedTrigger = () => ({
  triggerID: 'T1',
  workflowID: 'W1',
  stepID: '3',
  enabled: true,
  resourceType: 'compose',
  eventType: 'onManual',
  constraints: [],
  input: {},
  meta: {
    description: '',
    visual: { id: '2', parent: '1', value: 'Manual', xywh: [100, 200, 200, 80], edges: [] },
  },
  createdAt: '2026-10-05T00:00:00Z',
})

async function mountEditor() {
  const api = {
    triggerList: vi.fn(async () => ({ set: [storedTrigger()] })),
    workflowRead: vi.fn(async () => ({ workflowID: 'W1', steps: [], paths: [] })),
    workflowUpdate: vi.fn(async wf => wf),
    triggerUpdate: vi.fn(async t => t),
    triggerCreate: vi.fn(async t => t),
    triggerDelete: vi.fn(async () => {}),
  }

  const wrapper = mount(Editor, {
    global: {
      provide: { $AutomationAPI: api, $Auth: { user: { userID: 'U1' } } },
    },
  })
  await flushPromises()
  return { wrapper, api }
}

// What the graph encoder hands back on save: the editor adds a name to meta
function encoded(trigger) {
  return { ...trigger, meta: { name: trigger.meta.visual.value, ...trigger.meta } }
}

describe('workflow editor trigger saving', () => {
  it('does not update unchanged triggers', async () => {
    const { wrapper, api } = await mountEditor()

    wrapper.findComponent({ name: 'WorkflowEditor' }).vm.$emit('save', {
      workflowID: 'W1',
      triggers: [encoded(storedTrigger())],
    })
    await flushPromises()

    expect(api.triggerUpdate).not.toHaveBeenCalled()
    expect(api.triggerDelete).not.toHaveBeenCalled()
    expect(api.workflowUpdate).toHaveBeenCalledOnce()
  })

  it('updates a moved trigger', async () => {
    const { wrapper, api } = await mountEditor()

    const moved = storedTrigger()
    moved.meta.visual.xywh = [140, 200, 200, 80]
    wrapper.findComponent({ name: 'WorkflowEditor' }).vm.$emit('save', {
      workflowID: 'W1',
      triggers: [encoded(moved)],
    })
    await flushPromises()

    expect(api.triggerUpdate).toHaveBeenCalledOnce()
    expect(api.workflowUpdate).toHaveBeenCalledOnce()
  })
})
