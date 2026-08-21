import { describe, it, expect } from 'vitest'
import { adminRecordBlocks, adminRecordListBlocks, isWideField } from './record-blocks'

const SYSTEM = [
  'recordID',
  'ownedBy',
  'createdBy',
  'createdAt',
  'updatedBy',
  'updatedAt',
  'revision',
  'deletedBy',
  'deletedAt',
]

const field = (name: string, extra = {}) => ({ name, kind: 'String', ...extra })

const makeModule = (fields, config = {}) => ({
  moduleID: '510151835201437697',
  fields,
  config,
  systemFields: () => SYSTEM.map(name => ({ name, isSystem: true })),
})

const module_ = makeModule([
  field('code'),
  field('title'),
  field('notes', { options: { multiLine: true } }),
  field('plan', { kind: 'File' }),
])

const build = (mod = module_, opts = {}) =>
  adminRecordBlocks(mod, { idPrefix: '_t', systemTitle: 'System fields', ...opts })

describe('isWideField', () => {
  it('calls out the fields a column cannot hold', () => {
    expect(isWideField(field('notes', { options: { multiLine: true } }))).toBe(true)
    expect(isWideField(field('body', { options: { useRichTextEditor: true } }))).toBe(true)
    expect(isWideField(field('plan', { kind: 'File' }))).toBe(true)
    expect(isWideField(field('where', { kind: 'Geometry' }))).toBe(true)
    expect(isWideField(field('tags', { isMulti: true }))).toBe(true)
  })

  it('leaves an ordinary field to the columns', () => {
    expect(isWideField(field('code'))).toBe(false)
    expect(isWideField(field('signed', { kind: 'DateTime' }))).toBe(false)
    expect(isWideField(field('owner', { kind: 'User' }))).toBe(false)
  })
})

describe('adminRecordBlocks', () => {
  it('puts the column-friendly fields in one wrapped card', () => {
    const [compact] = build()

    expect(compact.options.fields).toEqual(['code', 'title'])
    expect(compact.options.recordFieldLayoutOption).toBe('wrap')
    expect(compact.title).toBe('')
  })

  it('gives every wide field a card of its own', () => {
    const [, notes, plan] = build()

    expect(notes.options.fields).toEqual(['notes'])
    expect(plan.options.fields).toEqual(['plan'])
    // A wide field is read down the page, not across columns
    expect(notes.options.recordFieldLayoutOption).not.toBe('wrap')
  })

  it('ends with the system fields', () => {
    const blocks = build()
    const system = blocks[blocks.length - 1]

    expect(system.title).toBe('System fields')
    expect(system.options.fields).toEqual(SYSTEM)
  })

  it('leaves out the compact card when every field is wide', () => {
    const blocks = build(makeModule([field('notes', { options: { multiLine: true } })]))

    // An empty field list would read as "every field", drawing them twice
    expect(blocks.every(b => b.options.fields.length > 0)).toBe(true)
    expect(blocks).toHaveLength(2)
  })

  it('stacks every block below the one before it', () => {
    const blocks = build()

    let y = 0
    for (const b of blocks) {
      const [x, top, w, h] = b.xywh
      expect([x, top]).toEqual([0, y])
      expect(w).toBe(48)
      expect(h).toBeGreaterThan(0)
      y += h
    }
  })

  it('leaves out what the server assigns, on a record that has none of it yet', () => {
    const blocks = build(module_, { isNew: true })
    const system = blocks[blocks.length - 1]

    expect(system.options.fields).not.toContain('recordID')
    expect(system.options.fields).not.toContain('revision')
    expect(system.options.fields).toContain('ownedBy')
  })

  it('offers revisions only where the module keeps them', () => {
    const kinds = mod => build(mod, { withRevisions: true }).map(b => b.kind)

    expect(kinds(module_)).not.toContain('RecordRevisions')
    expect(kinds(makeModule(module_.fields, { recordRevisions: { enabled: true } }))).toContain(
      'RecordRevisions',
    )
  })

  it('does not offer revisions to a screen that did not ask', () => {
    const mod = makeModule(module_.fields, { recordRevisions: { enabled: true } })

    expect(build(mod).map(b => b.kind)).not.toContain('RecordRevisions')
  })

  it('survives a module that is not loaded yet', () => {
    const blocks = adminRecordBlocks(null, { idPrefix: '_c', systemTitle: 'x' })

    expect(blocks).toHaveLength(1)
    expect(blocks[0].xywh[3]).toBeGreaterThan(0)
  })
})

describe('adminRecordListBlocks', () => {
  const built = () =>
    adminRecordListBlocks({
      moduleID: '510151835201437697',
      idPrefix: '_t',
      metricLabels: {
        total: 'Total records',
        createdRecently: 'New',
        updatedRecently: 'Touched',
        ownedByMe: 'Mine',
      },
    })

  it('lays the tiles out across one row', () => {
    const { tiles } = built()

    expect(tiles).toHaveLength(4)
    expect(tiles.map(t => t.xywh[0])).toEqual([0, 12, 24, 36])
    expect(tiles.every(t => t.xywh[1] === 0 && t.xywh[2] === 12)).toBe(true)
  })

  it('labels each tile and points it at the module', () => {
    const { tiles } = built()
    const metrics = tiles.map(t => t.options.metrics[0])

    expect(metrics.map(m => m.label)).toEqual(['Total records', 'New', 'Touched', 'Mine'])
    expect(metrics.every(m => m.moduleID === '510151835201437697')).toBe(true)
  })

  it('counts the whole module on the first tile and narrows on the rest', () => {
    const [total, ...rest] = built().tiles.map(t => t.options.metrics[0])

    // An empty filter is both "count everything" and what clears the table
    expect(total.filter).toBe('')
    expect(rest.every(m => !!m.filter)).toBe(true)
  })

  it('drills every tile into the table beside it', () => {
    const { tiles, list } = built()

    for (const tile of tiles) {
      const { drillDown } = tile.options.metrics[0]
      expect(drillDown.enabled).toBe(true)
      expect(drillDown.blockID).toBe(list.blockID)
    }
  })

  it('keeps the table configurable and showing deleted records', () => {
    const { list } = built()

    expect(list.options.moduleID).toBe('510151835201437697')
    expect(list.options.showDeletedRecordsOption).toBe(true)
    // The block would otherwise default this button away
    expect(list.options.hideConfigureFieldsButton).toBe(false)
  })
})
