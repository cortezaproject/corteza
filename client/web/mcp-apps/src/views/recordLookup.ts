// What the record lookup view reads out of a compose_record_lookup result.

export type ViewField = {
  name: string
  label?: string
  kind: string
  multi?: boolean
  options?: unknown
}

export type RecordView = {
  module: { name: string; handle?: string }
  fields: ViewField[]
  links?: Record<string, string>
}

type RawValue = { name: string; value: string }

type RawRecord = {
  recordID: string
  values?: RawValue[]
}

export type Row = {
  _dataKey: string
  recordID: string
  values: Record<string, string | string[]>
}

export type Page = {
  rows: Row[]
  refs: Record<string, string>
  cursor?: string
  view?: RecordView
}

type ToolResult = {
  structuredContent?: Record<string, unknown>
  _meta?: Record<string, unknown>
}

// The result _meta key the server puts the view's columns and links under.
export const viewDataKey = 'human.dev/view'

// Kinds whose value is the ID of something else, shown by the name the result's
// refs dictionary gives it.
const refKinds = new Set(['Record', 'User'])

export function isRefKind(kind: string) {
  return refKinds.has(kind)
}

// A multi-value field repeats its name once per value; everything else holds one.
export function toRow(rec: RawRecord, fields: ViewField[]): Row {
  const multi = new Set(fields.filter(f => f.multi).map(f => f.name))
  const values: Record<string, string | string[]> = {}

  for (const { name, value } of rec.values ?? []) {
    if (multi.has(name)) {
      values[name] = [...((values[name] as string[]) ?? []), value]
    } else {
      values[name] = value
    }
  }

  return { _dataKey: rec.recordID, recordID: rec.recordID, values }
}

// A list carries 'records'; a lookup by recordID returns the record itself.
export function parseResult(result: ToolResult): Page {
  const data = (result.structuredContent ?? {}) as Record<string, unknown>
  const view = result._meta?.[viewDataKey] as RecordView | undefined
  const fields = view?.fields ?? []

  let records: RawRecord[] = []
  if (Array.isArray(data.records)) {
    records = data.records as RawRecord[]
  } else if (typeof data.recordID === 'string') {
    records = [data as unknown as RawRecord]
  }

  return {
    rows: records.map(r => toRow(r, fields)),
    refs: (data.refs as Record<string, string>) ?? {},
    cursor: (data.nextPageCursor as string) || undefined,
    view,
  }
}

// Columns to show: the module's fields when the server described them, or else
// whatever names the records carry, in the order first seen.
export function columnsOf(page: Page): ViewField[] {
  if (page.view?.fields?.length) {
    return page.view.fields
  }

  const seen = new Map<string, ViewField>()
  for (const row of page.rows) {
    for (const name of Object.keys(row.values)) {
      if (!seen.has(name)) seen.set(name, { name, kind: 'String' })
    }
  }
  return [...seen.values()]
}

export function refLabel(value: string | string[] | undefined, refs: Record<string, string>) {
  if (value === undefined) return ''
  const ids = Array.isArray(value) ? value : [value]
  return ids.map(id => refs[id] ?? id).join(', ')
}

// Appends a further page to what is shown, keeping the columns the first page
// described.
export function appendPage(shown: Page, next: Page): Page {
  return {
    rows: [...shown.rows, ...next.rows],
    refs: { ...shown.refs, ...next.refs },
    cursor: next.cursor,
    view: shown.view && {
      ...shown.view,
      links: { ...shown.view.links, ...next.view?.links },
    },
  }
}
