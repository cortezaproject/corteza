import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import CResourceList from './CResourceList.vue'
import { changedAtField } from '../../composables/useChangedAt'

window.matchMedia =
  window.matchMedia ||
  ((() => ({ matches: false, addEventListener() {}, removeEventListener() {} })) as never)

const at = '2026-03-03T10:00:00Z'
const items = [
  { id: '1', name: 'live', updatedAt: at },
  { id: '2', name: 'gone', deletedAt: at },
  { id: '3', name: 'paused', suspendedAt: at },
  { id: '4', name: 'shelved', archivedAt: at },
]

function mountList() {
  return mount(CResourceList, {
    props: {
      primaryKey: 'id',
      fields: [{ key: 'name', header: 'Name' }, changedAtField('Last change')],
      items,
      filter: {},
      sorting: {},
      pagination: { total: items.length, limit: 10, page: 1 },
    },
  })
}

describe('CResourceList row state', () => {
  it('tags the last-change cell of a deleted, suspended or archived row', () => {
    const rows = mountList().findAll('tbody tr')
    const tag = (i: number) => rows[i].find('[data-pc-name="tag"]')
    expect(tag(0).exists()).toBe(false)
    expect(tag(1).text()).toBe('general.resourceList.state.deleted')
    expect(tag(2).text()).toBe('general.resourceList.state.suspended')
    expect(tag(3).text()).toBe('general.resourceList.state.archived')
    expect(rows[1].findAll('td')[1].find('[data-pc-name="tag"]').exists()).toBe(true)
  })

  it('colours deleted red and suspended dark, so the two never blur', () => {
    const rows = mountList().findAll('tbody tr')
    const sev = (i: number) => rows[i].find('[data-pc-name="tag"]').attributes('data-p')
    expect([sev(1), sev(2), sev(3)]).toEqual(['danger', 'contrast', 'secondary'])
  })

  it('dims only the deleted row', () => {
    const rows = mountList().findAll('tbody tr')
    expect(rows.map(r => r.classes('[&>td]:text-muted-color'))).toEqual([false, true, false, false])
  })
})

describe('CResourceList active filters', () => {
  function mountFiltered(filter: Record<string, unknown>) {
    return mount(CResourceList, {
      props: {
        primaryKey: 'id',
        fields: [{ key: 'name', header: 'Name' }],
        items,
        filter,
        states: ['suspended', 'deleted'],
        filterDefaults: { query: '', suspended: '0', deleted: '0', kind: undefined },
        filterLabels: { kind: { label: 'Kind', value: (v: string) => `is ${v}` } },
        sorting: {},
        pagination: { total: items.length, limit: 10, page: 1 },
      },
    })
  }

  const chipsOf = (w: ReturnType<typeof mountFiltered>) =>
    w.findAll('[data-test-id="active-filters"] [data-pc-name="chip"]')

  it('shows no bar while the status is Active and nothing else is set', () => {
    const w = mountFiltered({ query: 'searching', suspended: '0', deleted: '0' })
    expect(w.find('[data-test-id="active-filters"]').exists()).toBe(false)
  })

  it('shows the status as one chip with its row-state tag', () => {
    const w = mountFiltered({ query: '', suspended: '1', deleted: '2' })
    const chips = chipsOf(w)
    expect(chips).toHaveLength(1)
    expect(chips[0].find('span').text()).toBe('general.resourceList.status.label')
    const tag = chips[0].find('[data-pc-name="tag"]')
    expect(tag.text()).toBe('general.resourceList.state.deleted')
    expect(tag.attributes('data-p')).toBe('danger')
  })

  it('shows other labelled filters beside it', () => {
    const w = mountFiltered({ suspended: '2', deleted: '0', kind: 'bot' })
    const custom = chipsOf(w)[1]
      .findAll('span:not([data-pc-section])')
      .map(s => s.text())
    expect(custom).toEqual(['Kind', 'is bot'])
  })

  it('puts a removed status back to Active, and reset clears everything', async () => {
    const w = mountFiltered({ query: 'q', suspended: '1', deleted: '2', kind: 'bot' })
    await w.find('[data-pc-name="chip"] [data-pc-section="removeicon"]').trigger('click')
    expect(w.emitted('update:filter')?.[0]).toEqual([{ suspended: '0', deleted: '0' }])
    await w.find('[data-test-id="active-filters"] button').trigger('click')
    expect(w.emitted('update:filter')?.[1]).toEqual([
      { suspended: '0', deleted: '0', kind: undefined },
    ])
  })
})
