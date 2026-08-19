import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'

import PageBlock from './PageBlock.vue'

// A block's title and description are author-typed templates, the way its
// prefilter is — the card chrome, the magnify dialog header and the iframe's
// accessible name all read the same string.

// Surfaces what the real card root gets: `id` falls through as an attribute,
// the classes ride the pass-through `pt.root`.
const CardStub = {
  name: 'CardStub',
  props: ['pt'],
  template: `<div :class="pt?.root?.class">
    <div class="card-title"><slot name="title" /></div>
    <div class="card-subtitle"><slot name="subtitle" /></div>
    <slot name="content" />
  </div>`,
}

// Renders the header slot whether or not the dialog is open, so the header's
// binding is assertable without driving the magnify button.
const DialogStub = {
  name: 'DialogStub',
  template: '<div class="dialog"><slot name="header" /></div>',
}

const USER = { userID: '42', name: 'Ada' }
const RECORD = { recordID: '7', ownedBy: '3', values: { name: 'Acme' } }

function mountBlock(block, record) {
  return mount(PageBlock, {
    props: { block, record },
    global: {
      stubs: { Card: CardStub, Dialog: DialogStub, Button: true },
      directives: { tooltip: {} },
      provide: { $Auth: { user: USER } },
      mocks: { $t: k => k },
    },
  })
}

const title = w => w.get('.card-title').text()
const subtitle = w => w.get('.card-subtitle').text()

describe('PageBlock title and description', () => {
  it('interpolates the record into the title and the description', () => {
    const w = mountBlock(
      { title: '${record.values.name} — ${recordID}', description: 'owned by ${ownerID}' },
      RECORD,
    )

    expect(title(w)).toBe('Acme — 7')
    expect(subtitle(w)).toBe('owned by 3')
  })

  it('interpolates the signed-in user with no record', () => {
    const w = mountBlock({ title: 'Welcome, ${user.name}', description: 'user ${userID}' })

    expect(title(w)).toBe('Welcome, Ada')
    expect(subtitle(w)).toBe('user 42')
  })

  it('shows a record template as authored off a record page', () => {
    const w = mountBlock({ title: '${record.values.name}', description: '' })

    expect(title(w)).toBe('${record.values.name}')
  })

  it('shows an unparsable template as authored', () => {
    const w = mountBlock({ title: '${record.values.name', description: '' }, RECORD)

    expect(title(w)).toBe('${record.values.name')
  })

  it('interpolates the magnify dialog header from the same title', () => {
    const w = mountBlock(
      { title: '${record.values.name}', options: { magnifyOption: 'modal' } },
      RECORD,
    )

    expect(w.get('.dialog').text()).toBe('Acme')
  })

  it('leaves a plain title alone', () => {
    const w = mountBlock({ title: 'Open leads', description: 'All of them' })

    expect(title(w)).toBe('Open leads')
    expect(subtitle(w)).toBe('All of them')
  })
})

// The block root is the only hook a stylesheet written in admin > theming has:
// the author's ID and classes from the configurator, plus `page-block` and the
// kind so a rule can reach a whole class of blocks.
describe('PageBlock custom CSS hooks', () => {
  const root = w => w.get('.page-block')

  it('puts the authored ID and classes on the block root', () => {
    const w = mountBlock({
      kind: 'Chart',
      blockID: '12',
      meta: { customID: 'sales-chart', customCSSClass: 'wide highlight' },
    })

    expect(root(w).attributes('id')).toBe('sales-chart')
    expect(root(w).classes()).toEqual(
      expect.arrayContaining(['page-block', 'Chart', 'wide', 'highlight']),
    )
  })

  it('falls back to the block ID when none is authored', () => {
    const w = mountBlock({ kind: 'Content', blockID: '12' })

    expect(root(w).attributes('id')).toBe('12')
  })

  it('leaves an unsaved block without an ID', () => {
    const w = mountBlock({ kind: 'Content', blockID: '0', meta: { tempID: 'abc' } })

    expect(root(w).attributes('id')).toBeUndefined()
  })
})
