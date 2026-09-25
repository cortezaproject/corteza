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
        filterDefaults: { query: '', deleted: '0', status: undefined },
        filterLabels: {
          deleted: 'Deleted roles',
          status: { label: 'Status', value: (v: string) => `is ${v}` },
        },
        sorting: {},
        pagination: { total: items.length, limit: 10, page: 1 },
      },
    })
  }

  it('shows no bar while every labelled filter is at its default', () => {
    const w = mountFiltered({ query: 'searching', deleted: '0' })
    expect(w.find('[data-test-id="active-filters"]').exists()).toBe(false)
  })

  it('shows a chip per changed filter, state and custom alike', () => {
    const w = mountFiltered({ query: '', deleted: '2', status: 'draft' })
    const chips = w.findAll('[data-test-id="active-filters"] [data-pc-name="chip"]')
    const parts = chips.map(c => c.findAll('span:not([data-pc-section])').map(s => s.text()))
    expect(parts).toEqual([
      ['Deleted roles', 'general.resourceList.filter.exclusive'],
      ['Status', 'is draft'],
    ])
  })

  it('puts a removed chip, or all of them on reset, back to the default', async () => {
    const w = mountFiltered({ query: 'q', deleted: '1', status: 'draft' })
    await w.find('[data-pc-name="chip"] [data-pc-section="removeicon"]').trigger('click')
    expect(w.emitted('update:filter')?.[0]).toEqual([{ deleted: '0' }])
    await w.find('[data-test-id="active-filters"] button').trigger('click')
    expect(w.emitted('update:filter')?.[1]).toEqual([{ deleted: '0', status: undefined }])
  })
})
