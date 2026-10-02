import { expect } from 'chai'
import { PageBlockProgress, aggregateReport } from './progress'

describe('compose/types/page-block/progress', () => {
  describe('aggregateReport', () => {
    const rows = [{ rp: 4 }, { rp: 10 }, { count: 1 }]

    it('sums by default', () => {
      expect(aggregateReport(rows, '', 0)).to.equal(15)
    })

    it('handles min, max and avg', () => {
      expect(aggregateReport(rows, 'min', 0)).to.equal(1)
      expect(aggregateReport(rows, 'max', 0)).to.equal(10)
      expect(aggregateReport(rows, 'avg', 0)).to.equal(5)
    })

    it('falls back to the default on an empty report instead of NaN', () => {
      for (const op of ['', 'min', 'max', 'avg']) {
        expect(aggregateReport([], op, 0)).to.equal(0)
        expect(aggregateReport([], op, 100)).to.equal(100)
      }
    })

    it('ignores rows without a usable number', () => {
      expect(aggregateReport([{ rp: null }, { rp: 'x' }, { rp: '3' }], 'avg', 0)).to.equal(3)
    })
  })

  describe('fetch', () => {
    const api: any = {
      recordReport: () => Promise.resolve([]),
    }

    it('reports the configured defaults when the module has no records', async () => {
      const block = new PageBlockProgress({
        options: {
          value: { moduleID: '1', field: 'count', operation: 'avg' },
          maxValue: { moduleID: '1', field: 'amount', operation: 'max', default: 50 },
        },
      } as any)

      const { value, min, max } = await block.fetch({ value: {}, minValue: {}, maxValue: {} } as any, api, '1') as any

      expect(value).to.equal(0)
      expect(min).to.equal(0)
      expect(max).to.equal(50)
    })
  })
})
