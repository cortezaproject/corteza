import { describe, expect, it } from 'vitest'
import { RESOURCES, WORKFLOW_OUTCOMES, foldOutcomes, liveCount, orderedStatuses } from './resources'

describe('dashboard resources', () => {
  it('orders listed statuses first and appends unknown ones sorted', () => {
    const projects = RESOURCES.find(r => r.key === 'projects')
    const status = { deleted: 3, zzz: 1, active: 2, aaa: 4, draft: 5 }
    expect(orderedStatuses(projects, status)).toEqual(['active', 'draft', 'deleted', 'aaa', 'zzz'])
  })

  it('counts only live statuses toward the leading number', () => {
    expect(liveCount({ active: 2, enabled: 3, published: 1, deleted: 9, disabled: 4 })).toBe(6)
    expect(liveCount({})).toBe(0)
  })

  it('folds raw session statuses into display outcomes, per bucket', () => {
    const runs = {
      byStatus: { completed: 5, failed: 2, started: 1, prompted: 1, suspended: 1, canceled: 0 },
      series: {
        completed: [1, 4],
        failed: [0, 2],
        started: [1, 0],
        prompted: [0, 1],
        suspended: [1, 0],
      },
    }
    const { byStatus, series } = foldOutcomes(runs, WORKFLOW_OUTCOMES, 2)
    expect(byStatus).toEqual({ failed: 2, running: 3, canceled: 0, completed: 5 })
    expect(series.running).toEqual([2, 1])
    expect(series.canceled).toEqual([0, 0])
    expect(series.completed).toEqual([1, 4])
  })

  it('folds a missing section into zeros of the bucket length', () => {
    const { byStatus, series } = foldOutcomes(null, WORKFLOW_OUTCOMES, 3)
    expect(byStatus.failed).toBe(0)
    expect(series.failed).toEqual([0, 0, 0])
  })
})
