/* eslint-disable no-unused-expressions */
import { expect } from 'chai'
import {
  escapeQlString,
  getFieldFilter,
  getRecordListFilterSql,
  queryToFilter,
} from 'corteza-webapp-compose/src/lib/record-filter'

// Helper to build a single-field filter group as consumed by queryToFilter.
// `groupCondition` controls how the group connects to the previous group.
const group = (filter, groupCondition) => ({ filter, groupCondition })
const field = (name, value, { kind = 'String', operator = '=' } = {}) => ({ name, kind, value, operator })

describe('lib/record-filter', () => {
  describe('escapeQlString', () => {
    it('escapes single quotes and backslashes', () => {
      expect(escapeQlString("O'Brien")).to.equal("O\\'Brien")
      expect(escapeQlString('a\\b')).to.equal('a\\\\b')
    })

    it('escapes LIKE wildcards only in like mode', () => {
      expect(escapeQlString('50%_off', false)).to.equal('50%_off')
      expect(escapeQlString('50%_off', true)).to.equal('50\\\\%\\\\_off')
    })
  })

  describe('getFieldFilter', () => {
    it('builds a simple equality for strings', () => {
      expect(getFieldFilter('Name', 'String', 'bar', '=')).to.equal("(Name = 'bar')")
    })

    it('includes IS NULL fallback for != so nulls are not silently excluded', () => {
      expect(getFieldFilter('Name', 'String', 'bar', '!=')).to.equal("((Name != 'bar') OR (Name IS NULL))")
    })

    it('flips operands for IN', () => {
      expect(getFieldFilter('Name', 'String', 'bar', 'IN')).to.equal("('bar' IN Name)")
    })

    it('builds LIKE with wildcards', () => {
      expect(getFieldFilter('Name', 'String', 'foo', 'LIKE')).to.equal("(Name LIKE '%foo%')")
    })

    it('treats empty value as a NULL check', () => {
      expect(getFieldFilter('Name', 'String', '', '=')).to.equal('(Name IS NULL)')
      expect(getFieldFilter('Name', 'String', '', '!=')).to.equal('(Name IS NOT NULL)')
    })

    // Date-only comparisons must be inclusive of the whole day: comparing a datetime
    // field against DATE('YYYY-MM-DD') resolves the value to midnight, which drops
    // records logged later that same day. Comparing on DATE(field) keeps it day-level.
    it('compares date-only upper bounds on the field date so the end day is inclusive', () => {
      expect(getFieldFilter('created_at', 'DateTime', '2024-06-29', '<='))
        .to.equal("(DATE(created_at) <= DATE('2024-06-29'))")
    })

    it('compares date-only lower bounds on the field date', () => {
      expect(getFieldFilter('created_at', 'DateTime', '2024-06-26', '>='))
        .to.equal("(DATE(created_at) >= DATE('2024-06-26'))")
    })

    it('builds an inclusive date-only BETWEEN range on the field date', () => {
      expect(getFieldFilter('created_at', 'DateTime', { start: '2024-06-26', end: '2024-06-29' }, 'BETWEEN'))
        .to.equal("(DATE(created_at) BETWEEN DATE('2024-06-26') DATE('2024-06-29'))")
    })
  })

  describe('getRecordListFilterSql', () => {
    it('returns empty string for an empty filter', () => {
      expect(getRecordListFilterSql([])).to.equal('')
    })

    it('wraps a single field condition', () => {
      const sql = getRecordListFilterSql([field('A', '1')])
      expect(sql).to.equal("((A = '1'))")
    })

    it('groups OR conditions within a field group', () => {
      const sql = getRecordListFilterSql([
        { ...field('A', '1'), condition: '' },
        { ...field('B', '2'), condition: 'OR' },
      ])
      expect(sql).to.equal("(((A = '1') OR (B = '2')))")
    })
  })

  describe('queryToFilter', () => {
    const PRE = "(LocalGroupID = '5')"

    it('returns just the prefilter when no user filter or search is given', () => {
      expect(queryToFilter('', PRE, [], [])).to.equal(PRE)
    })

    it('AND-joins a single user filter group to the prefilter', () => {
      const out = queryToFilter('', PRE, [], [group([field('A', '1')])])
      expect(out).to.equal(`${PRE} AND ((A = '1'))`)
    })

    // The regression this suite primarily guards: a top-level OR in the user
    // filter must be parenthesised so the prefilter constrains BOTH branches.
    // Before the fix the result was `PRE AND (a) OR (b)`, which SQL reads as
    // `(PRE AND a) OR b`, letting branch `b` escape the prefilter entirely.
    it('wraps a top-level OR so the prefilter applies to every branch', () => {
      const out = queryToFilter('', PRE, [], [
        group([field('A', '1')]),
        group([field('B', '2')], 'OR'),
      ])
      expect(out).to.equal(`${PRE} AND (((A = '1')) OR ((B = '2')))`)

      // Semantic invariant: everything after the prefilter's AND is a single
      // parenthesised expression, so no OR branch can sit outside the prefilter.
      const tail = out.slice(`${PRE} AND `.length)
      expect(tail.startsWith('(')).to.be.true
      expect(tail.endsWith(')')).to.be.true
    })

    it('does not add a redundant outer wrap for pure AND groups', () => {
      const out = queryToFilter('', PRE, [], [
        group([field('A', '1')]),
        group([field('B', '2')], 'AND'),
      ])
      expect(out).to.equal(`${PRE} AND (((A = '1')) AND ((B = '2')))`)
    })

    it('handles mixed AND/OR with correct precedence grouping', () => {
      const out = queryToFilter('', PRE, [], [
        group([field('A', '1')]),
        group([field('B', '2')], 'AND'),
        group([field('C', '3')], 'OR'),
      ])
      expect(out).to.equal(`${PRE} AND ((((A = '1')) AND ((B = '2'))) OR ((C = '3')))`)
    })

    it('keeps the prefilter applied when an OR filter is combined with a search query', () => {
      const fields = [{ name: 'A', kind: 'String' }]
      const out = queryToFilter('foo', PRE, fields, [
        group([field('A', '1')]),
        group([field('B', '2')], 'OR'),
      ])
      expect(out).to.equal(`${PRE} AND (((A = '1')) OR ((B = '2'))) AND ((A LIKE '%foo%'))`)
    })

    it('works without a prefilter (OR still grouped on its own)', () => {
      const out = queryToFilter('', '', [], [
        group([field('A', '1')]),
        group([field('B', '2')], 'OR'),
      ])
      expect(out).to.equal("(((A = '1')) OR ((B = '2')))")
    })
  })
})
