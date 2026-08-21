import { expect } from 'chai'

import { PageBlockMetric } from './metric'

/**
 * The record report drops deleted records before any filter is applied, so
 * `deletedAt IS NOT NULL` counts nothing. A metric that reports on deleted
 * records has to carry the report's own state constraint instead, which means
 * the block has to hand it to the reporter.
 */

type Params = Record<string, unknown>

const captured: Params[] = []

const reporter = async (p: Params) => {
  captured.push(p)
  return [{ count: 2 }]
}

const fetchWith = async (metric: Record<string, unknown>) => {
  captured.length = 0
  const block = new PageBlockMetric()
  const out = await block.fetch(
    { m: { ...block.makeMetric(), ...metric } } as never,
    reporter as never,
  )
  return { out, params: captured[0] }
}

describe('PageBlockMetric reporter params', () => {
  it('asks for deleted records when the metric reports on them', async () => {
    const { params } = await fetchWith({ moduleID: '100001', deleted: 2 })

    expect(params.deleted).to.equal(2)
    expect(params.moduleID).to.equal('100001')
  })

  it('leaves deleted records out by default', async () => {
    const { params } = await fetchWith({ moduleID: '100001' })

    expect(params.deleted).to.equal(0)
  })

  it('counts records without asking for a metric expression', async () => {
    const { params } = await fetchWith({ moduleID: '100001', metricField: 'count' })

    // 'count' is the report's own column; an expression would shadow it
    expect(params.metrics).to.equal('')
    expect(params.dimensions).to.equal('deletedAt')
  })

  it('sums the groups the report returns', async () => {
    const block = new PageBlockMetric()
    const many = async () => [{ count: 2 }, { count: 3 }]
    const out = (await block.fetch(
      { m: { ...block.makeMetric(), moduleID: '100001' } } as never,
      many as never,
    )) as Array<{ value: number }>

    // A deleted-records report groups by deletedAt, so it can return several
    expect(out[0].value).to.equal(5)
  })

  it('reports zero when nothing matches', async () => {
    const block = new PageBlockMetric()
    const none = async () => []
    const out = (await block.fetch(
      { m: { ...block.makeMetric(), moduleID: '100001' } } as never,
      none as never,
    )) as Array<{ value: number }>

    expect(out[0].value).to.equal(0)
  })
})
