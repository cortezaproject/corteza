import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'

vi.mock('vue-i18n', () => ({
  // `t` echoes the key so assertions can target keys, but still exercises the
  // plural call signature (key, count) used for the change counter.
  useI18n: () => ({ t: (k, n) => (n === undefined ? k : `${k}:${n}`), locale: { value: 'en' } }),
}))

vi.mock('@planetcrust/human-vue', () => ({
  filters: { locFullDateTime: d => `FULL(${d})` },
}))

import EventTimelineItem from './EventTimelineItem.vue'

// Real payload shapes, copied from actual rows in the action log.
const updateEvent = {
  actionID: '1',
  timestamp: new Date().toISOString(),
  action: 'update',
  resource: 'corteza::system:user/482186405312724993',
  requestOrigin: 'api/rest',
  severity: 6,
  actorID: '482186405312724993',
  meta: { 'user.email': 'ada@example.com' },
  delta: [
    { key: 'meta.theme', old: ['light'], new: ['dark'] },
    { key: 'meta.avatarID', old: ['504643189915254785'], new: ['504643189915516929'] },
  ],
}

const createEvent = {
  actionID: '2',
  timestamp: new Date().toISOString(),
  action: 'create',
  resource: 'corteza::compose:record/504494280251604993/504494280253046785/504798810403307521',
  requestOrigin: 'api/rest',
  severity: 6,
  actorID: '482186405312724993',
  meta: { 'record.ID': '504798810403307521', 'module.name': 'Module A' },
  delta: [],
}

const mountItem = (data, props = {}) =>
  mount(EventTimelineItem, {
    props: { data, actorName: 'Ada Lovelace', ...props },
    global: {
      directives: { tooltip: {} },
      stubs: { UserCell: { props: ['name'], template: '<span>{{ name }}</span>' } },
    },
  })

describe('EventTimelineItem', () => {
  it('reads as a sentence: actor, verb, resource type', () => {
    const w = mountItem(updateEvent)
    const text = w.text()
    expect(text).toContain('Ada Lovelace')
    expect(text).toContain('updated')
    expect(text).toContain('User')
  })

  it('names the affected resource from meta rather than showing a bare ID', () => {
    const w = mountItem(updateEvent)
    expect(w.text()).toContain('ada@example.com')
  })

  it('falls back to the trailing ID when meta carries no name for this resource', () => {
    // module.name must NOT be borrowed for a record event's subject.
    const w = mountItem(createEvent)
    expect(w.text()).toContain('504798810403307521')
  })

  // The whole point of the redesign: the diff is answered without a click.
  it('shows the old→new diff inline, unexpanded', () => {
    const w = mountItem(updateEvent)
    const text = w.text()
    expect(text).toContain('meta.theme')
    expect(text).toContain('light')
    expect(text).toContain('dark')
  })

  it('caps the inline diff at 3 fields and offers the rest', () => {
    const many = {
      ...updateEvent,
      delta: Array.from({ length: 5 }, (_, i) => ({
        key: `field${i}`,
        old: [`old${i}`],
        new: [`new${i}`],
      })),
    }
    const w = mountItem(many)
    expect(w.text()).toContain('field2')
    expect(w.text()).not.toContain('field3')
    expect(w.text()).toContain('project.dashboard.activity.moreChanges:2')
  })

  it('reveals the capped fields on request', async () => {
    const many = {
      ...updateEvent,
      delta: Array.from({ length: 5 }, (_, i) => ({
        key: `field${i}`,
        old: [`old${i}`],
        new: [`new${i}`],
      })),
    }
    const w = mountItem(many)
    await w.find('button').trigger('click')
    expect(w.text()).toContain('field4')
    expect(w.text()).not.toContain('project.dashboard.activity.moreChanges')
  })

  it('offers no "more" affordance when every change already fits', () => {
    const w = mountItem(updateEvent) // 2 fields, under the cap
    expect(w.text()).not.toContain('project.dashboard.activity.moreChanges')
  })

  it('renders no diff for events without one', () => {
    const w = mountItem(createEvent)
    expect(w.findComponent({ name: 'EventDiff' }).exists()).toBe(false)
  })

  it('renders a relative timestamp, not a raw one', () => {
    const w = mountItem({ ...updateEvent, timestamp: new Date(Date.now() - 120000).toISOString() })
    expect(w.text()).toMatch(/minutes ago/)
  })

  it('emits a resource filter carrying the bare type, without the IDs', async () => {
    const w = mountItem(createEvent)
    await w.find('[role="button"]').trigger('click')
    expect(w.emitted('filter')[0]).toEqual(['resource', 'corteza::compose:record'])
  })

  it('falls back to System when there is no actor', () => {
    const w = mountItem(updateEvent, { actorName: '' })
    expect(w.text()).toContain('project.dashboard.activity.system')
  })

  // Icons must come from the project's own KIND_CONFIG, not a bespoke map — an
  // agent event should wear the same mark the agent wears everywhere else.
  it('uses the project kind icon for the resource', () => {
    const w = mountItem({ ...createEvent, resource: 'corteza::system:agent/1' })
    expect(w.html()).toContain('pi-sparkles') // KIND_CONFIG.agent.icon
    expect(w.html()).toContain('text-fuchsia-600')
  })

  it('maps a record to its parent module kind', () => {
    const w = mountItem(createEvent) // compose:record
    expect(w.html()).toContain('pi-database') // KIND_CONFIG.module.icon
  })

  it('falls back to the neutral kind for non-project resources', () => {
    const w = mountItem({ ...createEvent, resource: 'corteza::system:settings/1' })
    expect(w.html()).toContain('pi-circle') // kindConfig FALLBACK
  })

  it('repaints the icon for failed events but keeps the kind glyph', () => {
    const w = mountItem({ ...createEvent, resource: 'corteza::system:agent/1', severity: 3 })
    expect(w.html()).toContain('pi-sparkles') // still an agent
    expect(w.html()).toContain('text-red-600') // but alarmed
    expect(w.html()).not.toContain('text-fuchsia-600')
  })

  // The contract PrimeVue's Accordion would have given us; since this disclosure
  // is hand-rolled, assert it rather than trust it.
  it('wires the details toggle to its region via aria-controls/id', async () => {
    const w = mountItem(createEvent)
    const toggle = w.find('[aria-expanded]')

    expect(toggle.attributes('aria-expanded')).toBe('false')
    const controls = toggle.attributes('aria-controls')
    expect(controls).toBe('event-details-2')
    // Collapsed: nothing should claim that id.
    expect(w.find(`#${controls}`).exists()).toBe(false)

    await toggle.trigger('click')

    expect(toggle.attributes('aria-expanded')).toBe('true')
    const region = w.find(`#${controls}`)
    expect(region.exists()).toBe(true)
    expect(region.attributes('role')).toBe('region')
  })
})
