import { describe, expect, it } from 'vitest'
import { viewDataKey } from './recordLookup'
import {
  fromLines,
  initialValues,
  isReadOnly,
  linesOf,
  missingRequired,
  parseDraft,
  valuesToSave,
} from './recordForm'

const view = {
  module: { name: 'Deals' },
  fields: [
    { name: 'title', kind: 'String', required: true },
    { name: 'tags', kind: 'String', multi: true },
    { name: 'stage', kind: 'Select' },
    { name: 'notes', kind: 'String', options: { useRichTextEditor: true } },
    { name: 'files', kind: 'File', multi: true },
  ],
}

const draft = parseDraft({
  structuredContent: {
    draft: { title: 'Acme', tags: 'solo', stage: 'lead', notes: '<p>x</p>', files: ['9'] },
    namespaceID: '1',
    moduleID: '2',
    recordID: '3',
  },
  _meta: { [viewDataKey]: view },
})

describe('initialValues', () => {
  it('holds a multi-value field as an array and a single one as a value', () => {
    expect(initialValues(draft)).toMatchObject({ title: 'Acme', tags: ['solo'], stage: 'lead' })
  })
})

describe('valuesToSave', () => {
  it('sends filled fields, clears emptied ones, and skips read-only ones', () => {
    const form = { ...initialValues(draft), stage: '', tags: ['a', '', 'b'] }
    expect(valuesToSave(draft, form)).toEqual({ title: 'Acme', tags: ['a', 'b'], stage: [] })
  })

  it('leaves out a field that was empty and still is', () => {
    const d = parseDraft({ structuredContent: { draft: {} }, _meta: { [viewDataKey]: view } })
    expect(valuesToSave(d, { title: 'x', stage: '' })).toEqual({ title: 'x' })
  })
})

describe('read-only fields', () => {
  it('covers files and rich text, not plain text', () => {
    expect(view.fields.filter(isReadOnly).map(f => f.name)).toEqual(['notes', 'files'])
  })
})

describe('missingRequired', () => {
  it('names required fields left empty', () => {
    expect(missingRequired(draft, { title: '' })).toEqual(['title'])
    expect(missingRequired(draft, { title: 'x' })).toEqual([])
  })
})

describe('lines', () => {
  it('round-trips a multi-value String through text', () => {
    expect(fromLines(linesOf(['a', 'b']))).toEqual(['a', 'b'])
    expect(fromLines(' a \n\n b ')).toEqual(['a', 'b'])
  })
})
