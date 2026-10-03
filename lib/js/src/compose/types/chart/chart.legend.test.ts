import { expect } from 'chai'
import Chart from './chart'

// echarts 6 moved the default legend position to the bottom; charts keep
// showing it at the top unless a position is configured
describe('chart legend position', () => {
  const data = { labels: ['a'], datasets: [] }

  it('stays at the top by default', () => {
    const c = new Chart({ config: { reports: [{ legend: { position: { isDefault: true } } }] } } as any)
    const { legend } = c.makeOptions(data)
    expect(legend.top).to.eq('top')
    expect(legend.bottom).to.eq(undefined)
  })

  it('honours a configured position', () => {
    const c = new Chart({ config: { reports: [{ legend: { position: { isDefault: false, bottom: 10, left: 5 } } }] } } as any)
    const { legend } = c.makeOptions(data)
    expect(legend.top).to.eq(undefined)
    expect(legend.bottom).to.eq(10)
    expect(legend.left).to.eq(5)
  })
})
