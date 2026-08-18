import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'

vi.mock('./PageBlock.vue', () => ({
  default: { template: '<div><slot /></div>' },
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: k => k }),
}))

vi.mock('../registry', () => ({
  resolveBlock: () => ({ template: '<div class="inner" />' }),
}))

import TabsBlock from './TabsBlock.vue'

// A tab label is an author-typed template, the same as the block title beneath
// it — `${record.values.name}` names the tab after the record that is open.

const TabsStub = { name: 'TabsStub', template: '<div><slot /></div>' }
const TabListStub = { name: 'TabListStub', template: '<div><slot /></div>' }
const TabStub = { name: 'TabStub', template: '<div class="tab"><slot /></div>' }
const PanelsStub = { name: 'PanelsStub', template: '<div />' }

const USER = { userID: '42', name: 'Ada' }
const RECORD = { recordID: '7', ownedBy: '3', values: { name: 'Acme' } }

const INNER = { blockID: 'B2', kind: 'Content', title: 'Inner' }

function mountBlock(tabs, record) {
  return mount(TabsBlock, {
    props: {
      block: { blockID: 'B1', options: { tabs } },
      blocks: [INNER],
      record,
    },
    global: {
      stubs: {
        Tabs: TabsStub,
        TabList: TabListStub,
        Tab: TabStub,
        TabPanels: PanelsStub,
        TabPanel: PanelsStub,
        Menu: true,
      },
      directives: { ripple: {} },
      provide: { $Auth: { user: USER }, $pageBuilder: null },
      mocks: { $t: k => k },
    },
  })
}

const labels = w => w.findAll('.tab').map(t => t.text())

describe('TabsBlock tab labels', () => {
  it('interpolates the record into a tab label', () => {
    const w = mountBlock([{ blockID: 'B2', title: '${record.values.name} (${recordID})' }], RECORD)

    expect(labels(w)).toEqual(['Acme (7)'])
  })

  it('interpolates the signed-in user with no record', () => {
    const w = mountBlock([{ blockID: 'B2', title: 'Welcome, ${user.name}' }])

    expect(labels(w)).toEqual(['Welcome, Ada'])
  })

  it('shows a record template as authored off a record page', () => {
    const w = mountBlock([{ blockID: 'B2', title: '${record.values.name}' }])

    expect(labels(w)).toEqual(['${record.values.name}'])
  })

  it('leaves a plain label alone', () => {
    const w = mountBlock([{ blockID: 'B2', title: 'Details' }], RECORD)

    expect(labels(w)).toEqual(['Details'])
  })
})
