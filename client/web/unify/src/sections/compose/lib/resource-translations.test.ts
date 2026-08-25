import { describe, it, expect } from 'vitest'
import {
  applyPageLayoutTranslations,
  applyPageTranslations,
  moduleFieldKeyLabel,
  pageKeyLabel,
} from './resource-translations'

// The translator shows one row per translation key. A module field's keys are
// dotted paths the person translating never wrote, so each one is named; a key
// the mapper does not know returns '' and the form falls back to its own
// generic formatting.
const t = (k: string, p?: Record<string, string>) => (p ? `${k}:${Object.values(p).join('/')}` : k)

describe('moduleFieldKeyLabel', () => {
  it('names a select option after the option value it belongs to', () => {
    expect(moduleFieldKeyLabel('meta.options.new.text', t)).toBe('translator.keys.option:new')
  })

  it('keeps an option value that contains dots whole', () => {
    expect(moduleFieldKeyLabel('meta.options.a.b.text', t)).toBe('translator.keys.option:a.b')
  })

  it('names the label, description, hint and bool keys', () => {
    expect(moduleFieldKeyLabel('label', t)).toBe('translator.keys.label')
    expect(moduleFieldKeyLabel('meta.description.view', t)).toBe('translator.keys.description-view')
    expect(moduleFieldKeyLabel('meta.description.edit', t)).toBe('translator.keys.description-edit')
    expect(moduleFieldKeyLabel('meta.hint.view', t)).toBe('translator.keys.hint-view')
    expect(moduleFieldKeyLabel('meta.hint.edit', t)).toBe('translator.keys.hint-edit')
    expect(moduleFieldKeyLabel('meta.bool.true.label', t)).toBe('translator.keys.bool-true')
    expect(moduleFieldKeyLabel('meta.bool.false.label', t)).toBe('translator.keys.bool-false')
  })

  it('names a validator error whatever the validator id is', () => {
    expect(moduleFieldKeyLabel('expression.validator.12345.error', t)).toBe(
      'translator.keys.validator-error',
    )
  })

  it("returns '' for a key that is not a field's own", () => {
    expect(moduleFieldKeyLabel('name', t)).toBe('')
    expect(moduleFieldKeyLabel('title', t)).toBe('')
    expect(moduleFieldKeyLabel('pageBlock.3.title', t)).toBe('')
  })
})

// A page's title is `title`; a layout's is `meta.title`. Reading the page's key
// off a layout resource finds nothing, so a saved layout translation never
// reached the editor and the row kept showing the old title until a reload.
describe('applyPageLayoutTranslations', () => {
  const layout = () => ({
    namespaceID: '1',
    pageID: '2',
    pageLayoutID: '3',
    meta: { title: 'old title', description: 'old description' },
    config: { buttons: { back: { label: 'Back' } } },
  })

  const rows = (pairs: Array<[string, string]>) =>
    pairs.map(([key, message]) => ({
      resource: 'compose:page-layout/1/2/3',
      key,
      lang: 'fr',
      message,
    }))

  it('reads the layout title and description off their meta keys', () => {
    const l = layout()
    applyPageLayoutTranslations(
      l,
      rows([
        ['meta.title', 'Titre'],
        ['meta.description', 'Description'],
      ]),
      'fr',
    )
    expect(l.meta.title).toBe('Titre')
    expect(l.meta.description).toBe('Description')
  })

  it("ignores a page's own title key on a layout", () => {
    const l = layout()
    applyPageLayoutTranslations(l, rows([['title', 'Titre de la page']]), 'fr')
    expect(l.meta.title).toBe('old title')
  })

  it('still applies the toolbar button labels', () => {
    const l = layout()
    applyPageLayoutTranslations(l, rows([['config.buttons.back.label', 'Retour']]), 'fr')
    expect(l.config.buttons.back.label).toBe('Retour')
  })

  it('leaves a language it was not asked for alone', () => {
    const l = layout()
    applyPageLayoutTranslations(l, rows([['meta.title', 'Titre']]), 'sl')
    expect(l.meta.title).toBe('old title')
  })
})

// A page's dialog holds the page, its blocks and every layout at once, so one
// mapper names all three; a raw `pageBlock.3.button.2.label` says nothing about
// what it labels.
describe('pageKeyLabel', () => {
  const blocks = [
    {
      blockID: '3',
      title: 'Accounts',
      options: { buttons: [{ buttonID: '7', label: 'Run import' }] },
    },
    { blockID: '4', title: '' },
  ]

  it("names the page's own keys", () => {
    expect(pageKeyLabel('title', t, blocks)).toBe('translator.keys.page-title')
    expect(pageKeyLabel('description', t, blocks)).toBe('translator.keys.page-description')
  })

  it("names a layout's keys, which are meta keys", () => {
    expect(pageKeyLabel('meta.title', t, blocks)).toBe('translator.keys.layout-title')
    expect(pageKeyLabel('meta.description', t, blocks)).toBe('translator.keys.layout-description')
  })

  it('names each toolbar button and a custom action', () => {
    expect(pageKeyLabel('config.buttons.new.label', t, blocks)).toBe('translator.keys.toolbar-new')
    expect(pageKeyLabel('config.buttons.back.label', t, blocks)).toBe(
      'translator.keys.toolbar-back',
    )
    expect(pageKeyLabel('config.actions.12.meta.label', t, blocks)).toBe(
      'translator.keys.action-label',
    )
  })

  it('names a block key after the block, not its id', () => {
    expect(pageKeyLabel('pageBlock.3.title', t, blocks)).toBe(
      'translator.keys.block-title:Accounts',
    )
    expect(pageKeyLabel('pageBlock.3.content.body', t, blocks)).toBe(
      'translator.keys.block-content:Accounts',
    )
  })

  it('names a block button after the button it labels', () => {
    expect(pageKeyLabel('pageBlock.3.button.7.label', t, blocks)).toBe(
      'translator.keys.block-button:Accounts/Run import',
    )
  })

  it('falls back to the id where there is no name to use', () => {
    expect(pageKeyLabel('pageBlock.4.title', t, blocks)).toBe('translator.keys.block-title:#4')
    expect(pageKeyLabel('pageBlock.9.title', t, blocks)).toBe('translator.keys.block-title:#9')
    expect(pageKeyLabel('pageBlock.3.button.99.label', t, blocks)).toBe(
      'translator.keys.block-button:Accounts/#99',
    )
  })

  it("returns '' for a key that is not a page's or a layout's", () => {
    expect(pageKeyLabel('name', t, blocks)).toBe('')
    expect(pageKeyLabel('meta.options.new.text', t, blocks)).toBe('')
    expect(pageKeyLabel('pageBlock.3.somethingElse', t, blocks)).toBe('')
  })
})

describe('applyPageTranslations', () => {
  const rows = (pairs: Array<[string, string, string?]>) =>
    pairs.map(([key, message, resource]) => ({
      resource: resource || 'compose:page/1/2',
      key,
      lang: 'fr',
      message,
    }))

  const page = () => ({
    namespaceID: '1',
    pageID: '2',
    moduleID: '9',
    title: 'Account',
    description: '',
    config: { buttons: { back: { label: 'Back' } } },
    blocks: [{ blockID: '3', title: 'Accounts', description: '' }],
  })

  it('applies the page title and its blocks', () => {
    const p = page()
    applyPageTranslations(
      p,
      [],
      rows([
        ['title', 'Compte'],
        ['pageBlock.3.title', 'Comptes'],
      ]),
      'fr',
    )
    expect(p.title).toBe('Compte')
    expect(p.blocks[0].title).toBe('Comptes')
  })

  // Record-page toolbar labels live on the LAYOUT as config.buttons.*; the
  // page never carried a recordToolbar key, so reading one off it was dead.
  it('leaves the page config buttons to the layout', () => {
    const p = page()
    applyPageTranslations(p, [], rows([['recordToolbar.back.label', 'Retour']]), 'fr')
    expect(p.config.buttons.back.label).toBe('Back')
  })

  it("applies a layout's custom action labels", () => {
    const layout = {
      namespaceID: '1',
      pageID: '2',
      pageLayoutID: '3',
      meta: { title: '', description: '' },
      config: {
        buttons: {},
        actions: [{ actionID: '12', meta: { label: 'Go' } }, { meta: { label: 'Unsaved' } }],
      },
    }
    applyPageTranslations(
      page(),
      [layout],
      rows([
        ['config.actions.12.meta.label', 'Aller', 'compose:page-layout/1/2/3'],
        ['config.actions.1.meta.label', 'Non sauvé', 'compose:page-layout/1/2/3'],
      ]),
      'fr',
    )
    expect(layout.config.actions[0].meta.label).toBe('Aller')
    // The second has no id yet, so the server keys it by position.
    expect(layout.config.actions[1].meta.label).toBe('Non sauvé')
  })
})
