import { describe, it, expect } from 'vitest'
import {
  escapeQlString,
  evaluatePrefilter,
  getFieldFilter,
  getRecordListFilterSql,
  queryToFilter,
} from './record-filter'

// Helper to build a single-field filter group as consumed by queryToFilter.
// `groupCondition` controls how the group connects to the previous group.
const group = (filter: any[], groupCondition?: string) => ({
  filter,
  groupCondition,
})
const field = (
  name: string,
  value: string,
  { kind = 'String', operator = '=' }: { kind?: string; operator?: string } = {},
) => ({ name, kind, value, operator })

describe('lib/record-filter', () => {
  describe('escapeQlString', () => {
    it('escapes single quotes and backslashes', () => {
      expect(escapeQlString("O'Brien")).toBe("O\\'Brien")
      expect(escapeQlString('a\\b')).toBe('a\\\\b')
    })

    it('escapes LIKE wildcards only in like mode', () => {
      expect(escapeQlString('50%_off', false)).toBe('50%_off')
      expect(escapeQlString('50%_off', true)).toBe('50\\\\%\\\\_off')
    })
  })

  describe('getFieldFilter', () => {
    it('builds a simple equality for strings', () => {
      expect(getFieldFilter('Name', 'String', 'bar', '=')).toBe("(Name = 'bar')")
    })

    it('includes IS NULL fallback for != so nulls are not silently excluded', () => {
      expect(getFieldFilter('Name', 'String', 'bar', '!=')).toBe(
        "((Name != 'bar') OR (Name IS NULL))",
      )
    })

    it('flips operands for IN', () => {
      expect(getFieldFilter('Name', 'String', 'bar', 'IN')).toBe("('bar' IN Name)")
    })

    it('builds LIKE with wildcards', () => {
      expect(getFieldFilter('Name', 'String', 'foo', 'LIKE')).toBe("(Name LIKE '%foo%')")
    })

    it('treats empty value as a NULL check', () => {
      expect(getFieldFilter('Name', 'String', '', '=')).toBe('(Name IS NULL)')
      expect(getFieldFilter('Name', 'String', '', '!=')).toBe('(Name IS NOT NULL)')
    })

    // Date-only comparisons must be inclusive of the whole day: comparing a datetime
    // field against DATE('YYYY-MM-DD') resolves the value to midnight, which drops
    // records logged later that same day. Comparing on DATE(field) keeps it day-level.
    it('compares date-only upper bounds on the field date so the end day is inclusive', () => {
      expect(getFieldFilter('created_at', 'DateTime', '2024-06-29', '<=')).toBe(
        "(DATE(created_at) <= DATE('2024-06-29'))",
      )
    })

    it('compares date-only lower bounds on the field date', () => {
      expect(getFieldFilter('created_at', 'DateTime', '2024-06-26', '>=')).toBe(
        "(DATE(created_at) >= DATE('2024-06-26'))",
      )
    })

    it('builds an inclusive date-only BETWEEN range on the field date', () => {
      expect(
        getFieldFilter(
          'created_at',
          'DateTime',
          { start: '2024-06-26', end: '2024-06-29' },
          'BETWEEN',
        ),
      ).toBe("(DATE(created_at) BETWEEN DATE('2024-06-26') DATE('2024-06-29'))")
    })

    // These are the strings the filter UI actually produces: CInputDateTime emits
    // Date.toISOString() for a datetime field and HH:mm:ss for a time-only one
    // (lib/vue/src/components/input/CInputDateTime.vue). Neither matched the strict
    // formats parsed here, so getFieldFilter returned undefined and
    // getRecordListFilterSql dropped the condition — the filter looked applied and
    // never reached the API. Every other field kind built its condition fine, which
    // is what made it look like a date-only problem.
    // The emitted instant is re-rendered in the runner's timezone, so this pins the
    // shape of the comparison rather than a wall-clock string that only holds in UTC.
    it('accepts the ISO timestamp the date-time editor emits', () => {
      expect(getFieldFilter('created_at', 'DateTime', '2024-06-29T12:30:00.000Z', '>=')).toMatch(
        /^\(TIMESTAMP\(DATE_FORMAT\(created_at, '%Y-%m-%dT%H:%i:%s\.%f\+00:00'\)\) >= TIMESTAMP\(DATE_FORMAT\('2024-06-29T\d{2}:\d{2}:\d{2}[+-]\d{2}:\d{2}', '%Y-%m-%dT%H:%i:%s\.%f\+00:00'\)\)\)$/,
      )
    })

    it('accepts an ISO timestamp on both ends of a BETWEEN range', () => {
      expect(
        getFieldFilter(
          'created_at',
          'DateTime',
          { start: '2024-06-26T00:00:00.000Z', end: '2024-06-29T23:59:59.000Z' },
          'BETWEEN',
        ),
      ).toBeDefined()
    })

    // TIME(field), not the bare field: the value is stored as a timestamp, and
    // postgres rejects "timestamp with time zone >= time without time zone".
    it('accepts the HH:mm:ss the time-only editor emits, comparing on the field time', () => {
      expect(getFieldFilter('start_time', 'DateTime', '14:30:00', '>=')).toBe(
        "(TIME(start_time) >= TIME('14:30:00'))",
      )
    })

    it('compares a time-only BETWEEN range on the field time', () => {
      expect(
        getFieldFilter('start_time', 'DateTime', { start: '09:00:00', end: '17:00:00' }, 'BETWEEN'),
      ).toBe("(TIME(start_time) BETWEEN TIME('09:00:00') TIME('17:00:00'))")
    })
  })

  describe('evaluatePrefilter', () => {
    it('interpolates record/recordID/ownerID/userID template expressions', () => {
      const record = { values: { status: 'open' } }
      const user = { userID: '42' }
      expect(
        evaluatePrefilter('status = ${record.values.status} AND recordID = ${recordID}', {
          record,
          user,
          recordID: '7',
          ownerID: '3',
          userID: user.userID,
        }),
      ).toBe('status = open AND recordID = 7')
    })

    it('passes through plain strings with no expressions unchanged', () => {
      expect(
        evaluatePrefilter('LocalGroupID = 5', {
          record: undefined,
          user: undefined,
          recordID: '0',
          ownerID: '0',
          userID: '0',
        }),
      ).toBe('LocalGroupID = 5')
    })
  })

  describe('getRecordListFilterSql', () => {
    it('returns empty string for an empty filter', () => {
      expect(getRecordListFilterSql([])).toBe('')
    })

    it('wraps a single field condition', () => {
      expect(getRecordListFilterSql([field('A', '1')])).toBe("((A = '1'))")
    })

    it('groups OR conditions within a field group', () => {
      const sql = getRecordListFilterSql([
        { ...field('A', '1'), condition: '' },
        { ...field('B', '2'), condition: 'OR' },
      ])
      expect(sql).toBe("(((A = '1') OR (B = '2')))")
    })
  })

  describe('queryToFilter', () => {
    const PRE = "(LocalGroupID = '5')"

    it('returns just the prefilter when no user filter or search is given', () => {
      expect(queryToFilter('', PRE, [], [])).toBe(PRE)
    })

    it('AND-joins a single user filter group to the prefilter', () => {
      expect(queryToFilter('', PRE, [], [group([field('A', '1')])])).toBe(`${PRE} AND ((A = '1'))`)
    })

    // The regression this suite primarily guards: a top-level OR in the user
    // filter must be parenthesised so the prefilter constrains BOTH branches.
    // Before the fix the result was `PRE AND (a) OR (b)`, which SQL reads as
    // `(PRE AND a) OR b`, letting branch `b` escape the prefilter entirely.
    it('wraps a top-level OR so the prefilter applies to every branch', () => {
      const out = queryToFilter(
        '',
        PRE,
        [],
        [group([field('A', '1')]), group([field('B', '2')], 'OR')],
      )
      expect(out).toBe(`${PRE} AND (((A = '1')) OR ((B = '2')))`)

      // Semantic invariant: everything after the prefilter's AND is a single
      // parenthesised expression, so no OR branch can sit outside the prefilter.
      const tail = out.slice(`${PRE} AND `.length)
      expect(tail.startsWith('(')).toBe(true)
      expect(tail.endsWith(')')).toBe(true)
    })

    it('does not add a redundant outer wrap for pure AND groups', () => {
      const out = queryToFilter(
        '',
        PRE,
        [],
        [group([field('A', '1')]), group([field('B', '2')], 'AND')],
      )
      expect(out).toBe(`${PRE} AND (((A = '1')) AND ((B = '2')))`)
    })

    it('handles mixed AND/OR with correct precedence grouping', () => {
      const out = queryToFilter(
        '',
        PRE,
        [],
        [group([field('A', '1')]), group([field('B', '2')], 'AND'), group([field('C', '3')], 'OR')],
      )
      expect(out).toBe(`${PRE} AND ((((A = '1')) AND ((B = '2'))) OR ((C = '3')))`)
    })

    it('keeps the prefilter applied when an OR filter is combined with a search query', () => {
      const fields = [{ name: 'A', kind: 'String' }]
      const out = queryToFilter('foo', PRE, fields, [
        group([field('A', '1')]),
        group([field('B', '2')], 'OR'),
      ])
      expect(out).toBe(`${PRE} AND (((A = '1')) OR ((B = '2'))) AND ((A LIKE '%foo%'))`)
    })

    it('works without a prefilter (OR still grouped on its own)', () => {
      const out = queryToFilter(
        '',
        '',
        [],
        [group([field('A', '1')]), group([field('B', '2')], 'OR')],
      )
      expect(out).toBe("(((A = '1')) OR ((B = '2')))")
    })
  })
})
