import { describe, it, expect, vi, beforeEach } from 'vitest'
import { usePageVisibility, fetchBlockID, getBreakpoint } from './usePageVisibility'

function makeSystemAPI(results: Record<string, boolean> = {}) {
  return {
    expressionEvaluate: vi.fn().mockResolvedValue(results),
  }
}

function makeAuth(roles: string[] = []) {
  return { user: { roles } }
}

function makeLayout(pageLayoutID: string, opts: { expression?: string; roles?: string[] } = {}) {
  return {
    pageLayoutID,
    config: {
      visibility: {
        expression: opts.expression,
        roles: opts.roles || [],
      },
    },
  }
}

function makeBlock(
  blockID: string,
  opts: {
    expression?: string
    roles?: string[]
    kind?: string
    tempID?: string
    tabs?: any[]
  } = {},
) {
  return {
    blockID,
    kind: opts.kind || 'Content',
    meta: {
      tempID: opts.tempID,
      visibility: { expression: opts.expression, roles: opts.roles || [] },
    },
    options: opts.tabs !== undefined ? { tabs: opts.tabs } : undefined,
  }
}

describe('fetchBlockID', () => {
  it('returns blockID when set and non-zero', () => {
    expect(fetchBlockID({ blockID: '100', meta: {} })).toBe('100')
  })

  it('returns tempID when blockID is "0"', () => {
    expect(fetchBlockID({ blockID: '0', meta: { tempID: 'tmp-1' } })).toBe('tmp-1')
  })

  it('returns tempID when blockID absent', () => {
    expect(fetchBlockID({ meta: { tempID: 'tmp-2' } })).toBe('tmp-2')
  })

  it('returns empty string when both missing', () => {
    expect(fetchBlockID({ meta: {} })).toBe('')
  })
})

describe('getBreakpoint', () => {
  it('returns lg for wide screens', () => {
    Object.defineProperty(window, 'innerWidth', { writable: true, configurable: true, value: 1400 })
    expect(getBreakpoint()).toBe('lg')
  })

  it('returns md for 996-1199px', () => {
    Object.defineProperty(window, 'innerWidth', { writable: true, configurable: true, value: 1000 })
    expect(getBreakpoint()).toBe('md')
  })

  it('returns xxs for narrow screens', () => {
    Object.defineProperty(window, 'innerWidth', { writable: true, configurable: true, value: 300 })
    expect(getBreakpoint()).toBe('xxs')
  })
})

describe('usePageVisibility', () => {
  describe('buildExpressionVariables()', () => {
    it('includes user and screen', () => {
      const { buildExpressionVariables } = usePageVisibility(null, makeAuth(['role-1']))
      const vars = buildExpressionVariables()
      expect((vars.user as any).roles).toContain('role-1')
      expect(vars.screen).toBeTruthy()
    })

    it('adds mode flags when isRecordPage+mode provided', () => {
      const { buildExpressionVariables } = usePageVisibility(null, null)
      const vars = buildExpressionVariables({ isRecordPage: true, mode: 'edit' })
      expect(vars.isEdit).toBe(true)
      expect(vars.isView).toBe(false)
      expect(vars.isCreate).toBe(false)
    })

    it('omits mode flags when isRecordPage=false', () => {
      const { buildExpressionVariables } = usePageVisibility(null, null)
      const vars = buildExpressionVariables({ isRecordPage: false, mode: 'edit' })
      expect(vars.isEdit).toBeUndefined()
    })
  })

  describe('determineLayout()', () => {
    it('returns null for empty layouts', async () => {
      const { determineLayout } = usePageVisibility(null, null)
      expect(await determineLayout([], {})).toBeNull()
    })

    it('returns first layout with no visibility rules', async () => {
      const { determineLayout } = usePageVisibility(null, null)
      const layout = { pageLayoutID: 'L1' }
      expect(await determineLayout([layout], {})).toBe(layout)
    })

    it('evaluates expression via API and matches when true', async () => {
      const api = makeSystemAPI({ L1: true })
      const { determineLayout } = usePageVisibility(api, null)
      const layout = makeLayout('L1', { expression: 'user.age > 18' })
      const result = await determineLayout([layout], {})
      expect(result).toBe(layout)
      expect(api.expressionEvaluate).toHaveBeenCalled()
    })

    it('skips layout when expression evaluates to false', async () => {
      const api = makeSystemAPI({ L1: false })
      const { determineLayout } = usePageVisibility(api, null)
      const layout = makeLayout('L1', { expression: 'user.age > 18' })
      expect(await determineLayout([layout], {})).toBeNull()
    })

    it('skips layout when user lacks required role', async () => {
      const { determineLayout } = usePageVisibility(null, makeAuth(['other-role']))
      const layout = makeLayout('L1', { roles: ['admin-role'] })
      expect(await determineLayout([layout], {})).toBeNull()
    })

    it('matches layout when user has required role', async () => {
      const { determineLayout } = usePageVisibility(null, makeAuth(['admin-role']))
      const layout = makeLayout('L1', { roles: ['admin-role'] })
      expect(await determineLayout([layout], {})).toBe(layout)
    })

    it('passes the record on to layout expressions', async () => {
      const api = makeSystemAPI({ L1: true })
      const { buildExpressionVariables, determineLayout } = usePageVisibility(api, null)
      const record = { serialize: () => ({ values: { status: 'closed' } }) }
      const vars = buildExpressionVariables({ record, isRecordPage: true, mode: 'view' })

      const layout = makeLayout('L1', { expression: 'record.values.status == "closed"' })
      expect(await determineLayout([layout], vars)).toBe(layout)
      expect(api.expressionEvaluate).toHaveBeenCalledWith({
        variables: expect.objectContaining({ record: { values: { status: 'closed' } } }),
        expressions: { L1: 'record.values.status == "closed"' },
      })
    })

    it('retries without requestedLayoutID when specific layout does not match', async () => {
      const api = makeSystemAPI({ L1: false, L2: true })
      const { determineLayout } = usePageVisibility(api, null)
      const L1 = makeLayout('L1', { expression: 'false' })
      const L2 = makeLayout('L2', { expression: 'true' })
      const result = await determineLayout([L1, L2], {}, 'L1')
      expect(result?.pageLayoutID).toBe('L2')
    })

    it('opens the requested layout when its own condition passes', async () => {
      const api = makeSystemAPI({ L1: true, L2: true })
      const { determineLayout } = usePageVisibility(api, null)
      const L1 = makeLayout('L1', { expression: 'true' })
      const L2 = makeLayout('L2', { expression: 'true' })
      // L1 comes first in default order, so picking L2 proves the request won
      const result = await determineLayout([L1, L2], {}, 'L2')
      expect(result?.pageLayoutID).toBe('L2')
    })

    it('falls back to default order when the requested layout is barred by roles', async () => {
      const { determineLayout } = usePageVisibility(null, makeAuth(['viewer']))
      const L1 = makeLayout('L1')
      const L2 = makeLayout('L2', { roles: ['admin-role'] })
      const result = await determineLayout([L1, L2], {}, 'L2')
      expect(result?.pageLayoutID).toBe('L1')
    })

    it('treats expression as failed when API throws', async () => {
      const api = { expressionEvaluate: vi.fn().mockRejectedValue(new Error('eval error')) }
      const { determineLayout } = usePageVisibility(api, null)
      const layout = makeLayout('L1', { expression: 'user.x' })
      expect(await determineLayout([layout], {})).toBeNull()
    })
  })

  describe('evaluateBlocks()', () => {
    it('returns empty set when no visibility rules', async () => {
      const { evaluateBlocks } = usePageVisibility(null, null)
      const blocks = [makeBlock('B1'), makeBlock('B2')]
      const invisible = await evaluateBlocks(blocks, {})
      expect(invisible.size).toBe(0)
    })

    it('marks block invisible when expression is false', async () => {
      const api = makeSystemAPI({ B1: false })
      const { evaluateBlocks } = usePageVisibility(api, null)
      const blocks = [makeBlock('B1', { expression: 'user.x' })]
      const invisible = await evaluateBlocks(blocks, {})
      expect(invisible.has('B1')).toBe(true)
    })

    it('keeps block visible when expression is true', async () => {
      const api = makeSystemAPI({ B1: true })
      const { evaluateBlocks } = usePageVisibility(api, null)
      const blocks = [makeBlock('B1', { expression: 'user.x' })]
      const invisible = await evaluateBlocks(blocks, {})
      expect(invisible.has('B1')).toBe(false)
    })

    it('marks block invisible when user lacks required role', async () => {
      const { evaluateBlocks } = usePageVisibility(null, makeAuth(['user-role']))
      const blocks = [makeBlock('B1', { roles: ['admin-role'] })]
      const invisible = await evaluateBlocks(blocks, {})
      expect(invisible.has('B1')).toBe(true)
    })

    it('Tabs with all invisible children becomes invisible', async () => {
      const { evaluateBlocks } = usePageVisibility(null, makeAuth(['user']))
      const child = makeBlock('C1', { roles: ['admin'] }) // will be invisible
      const tabs: any = {
        blockID: 'T1',
        kind: 'Tabs',
        meta: { visibility: {} },
        options: { tabs: [{ blockID: 'C1' }] },
      }
      const invisible = await evaluateBlocks([child, tabs], {})
      expect(invisible.has('C1')).toBe(true)
      expect(invisible.has('T1')).toBe(true)
    })

    it('Tabs with at least one visible child stays visible', async () => {
      const { evaluateBlocks } = usePageVisibility(null, makeAuth(['admin']))
      const child1 = makeBlock('C1', { roles: ['admin'] }) // visible
      const child2 = makeBlock('C2', { roles: ['other'] }) // invisible
      const tabs: any = {
        blockID: 'T1',
        kind: 'Tabs',
        meta: { visibility: {} },
        options: { tabs: [{ blockID: 'C1' }, { blockID: 'C2' }] },
      }
      const invisible = await evaluateBlocks([child1, child2, tabs], {})
      expect(invisible.has('T1')).toBe(false)
    })

    it('uses tempID when blockID is 0', async () => {
      const api = makeSystemAPI({ 'tmp-1': false })
      const { evaluateBlocks } = usePageVisibility(api, null)
      const block: any = {
        blockID: '0',
        kind: 'Content',
        meta: { tempID: 'tmp-1', visibility: { expression: 'false' } },
      }
      const invisible = await evaluateBlocks([block], {})
      expect(invisible.has('tmp-1')).toBe(true)
    })
  })
})
