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
