import { describe, it, expect } from 'vitest'
import { adminRecordBlocks } from './record-blocks'

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

const module_ = {
  fields: [
    { name: 'title', kind: 'String' },
    { name: 'body', kind: 'String' },
  ],
  systemFields: () => SYSTEM.map(name => ({ name, isSystem: true })),
}

const build = () =>
  adminRecordBlocks(module_, { idPrefix: '_test', systemTitle: 'System fields' })

describe('adminRecordBlocks', () => {
  it('draws the module fields and the system fields as two blocks', () => {
    const [own, system] = build()

    expect(own.kind).toBe('Record')
    expect(system.kind).toBe('Record')
    // An empty list is what a Record block reads as "every module field"
    expect(own.options.fields).toEqual([])
    expect(system.options.fields).toEqual(SYSTEM)
  })

  it('titles only the system block', () => {
    const [own, system] = build()

    expect(own.title).toBe('')
    expect(system.title).toBe('System fields')
  })

  it('gives the two blocks distinct ids', () => {
    const [own, system] = build()

    expect(own.blockID).not.toBe(system.blockID)
    expect(own.meta.tempID).toBe(own.blockID)
    expect(system.meta.tempID).toBe(system.blockID)
  })

  it('stacks the system block directly below the module fields', () => {
    const [own, system] = build()
    const [ownX, ownY, ownW, ownH] = own.xywh
    const [sysX, sysY, sysW] = system.xywh

    expect([ownX, ownY]).toEqual([0, 0])
    expect(ownW).toBe(48)
    expect(sysW).toBe(48)
    expect(sysX).toBe(0)
    expect(sysY).toBe(ownH)
  })

  it('grows the module block with the field count', () => {
    const [small] = adminRecordBlocks(module_, { idPrefix: '_a', systemTitle: 'x' })
    const [large] = adminRecordBlocks(
      { ...module_, fields: Array.from({ length: 20 }, (_, i) => ({ name: `f${i}` })) },
      { idPrefix: '_b', systemTitle: 'x' },
    )

    expect(large.xywh[3]).toBeGreaterThan(small.xywh[3])
  })

  it('leaves out what the server assigns, on a record that has none of it yet', () => {
    const [, system] = adminRecordBlocks(module_, {
      idPrefix: '_new',
      systemTitle: 'System fields',
      isNew: true,
    })

    expect(system.options.fields).not.toContain('recordID')
    expect(system.options.fields).not.toContain('revision')
    expect(system.options.fields).toContain('ownedBy')
    expect(system.options.fields).toContain('createdAt')
  })

  it('survives a module that is not loaded yet', () => {
    const blocks = adminRecordBlocks(null, { idPrefix: '_c', systemTitle: 'x' })

    expect(blocks).toHaveLength(2)
    expect(blocks[1].options.fields).toEqual([])
    expect(blocks[0].xywh[3]).toBeGreaterThan(0)
  })
})
