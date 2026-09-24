// What the record form reads out of a compose_record_draft result, and the
// values it sends back to compose_record_create or compose_record_update.

import { viewDataKey } from './recordLookup'

export type FormField = {
  name: string
  label?: string
  kind: string
  multi?: boolean
  required?: boolean
  options?: Record<string, unknown>
}

export type Choice = { id: string; label: string }

export type FormView = {
  module: { name: string; handle?: string }
  fields: FormField[]
  choices?: Record<string, { items: Choice[]; more?: boolean }>
  link?: string
}

export type FormValue = string | string[]

export type Draft = {
  values: Record<string, FormValue>
  namespaceID: string
  moduleID: string
  recordID?: string
  refs: Record<string, string>
  view?: FormView
}

type ToolResult = {
  structuredContent?: Record<string, unknown>
  _meta?: Record<string, unknown>
}

// Kinds the form shows but cannot change: an upload or a map needs Human itself.
const readOnlyKinds = new Set(['File', 'Geometry'])

// Rich text is edited in Human too: its editor alone would double the view.
export function isReadOnly(f: FormField) {
  return readOnlyKinds.has(f.kind) || (f.kind === 'String' && !!f.options?.useRichTextEditor)
}

// A multi-value String is edited as text, one value per line.
export function linesOf(v: FormValue | undefined) {
  return [v ?? []].flat().join('\n')
}

export function fromLines(text: string) {
  return text
    .split('\n')
    .map(l => l.trim())
    .filter(Boolean)
}

export function parseDraft(result: ToolResult): Draft {
  const data = (result.structuredContent ?? {}) as Record<string, unknown>
  return {
    values: (data.draft as Record<string, FormValue>) ?? {},
    namespaceID: String(data.namespaceID ?? ''),
    moduleID: String(data.moduleID ?? ''),
    recordID: (data.recordID as string) || undefined,
    refs: (data.refs as Record<string, string>) ?? {},
    view: result._meta?.[viewDataKey] as FormView | undefined,
  }
}

// The form's starting state: every editable field, multi-value ones as arrays.
export function initialValues(draft: Draft): Record<string, FormValue> {
  const out: Record<string, FormValue> = {}
  for (const f of draft.view?.fields ?? []) {
    const v = draft.values[f.name]
    if (f.multi) {
      out[f.name] = Array.isArray(v) ? [...v] : v ? [v] : []
    } else {
      out[f.name] = Array.isArray(v) ? (v[0] ?? '') : (v ?? '')
    }
  }
  return out
}

function isEmpty(v: FormValue | undefined) {
  return v === undefined || v === '' || (Array.isArray(v) && v.filter(x => x !== '').length === 0)
}

// The values to save. An empty field is left out unless the draft had a value
// for it, and then it is sent as [] so the save clears it; a read-only field
// is never sent.
export function valuesToSave(
  draft: Draft,
  form: Record<string, FormValue>,
): Record<string, FormValue> {
  const out: Record<string, FormValue> = {}
  for (const f of draft.view?.fields ?? []) {
    if (isReadOnly(f)) continue

    const v = form[f.name]
    if (!isEmpty(v)) {
      out[f.name] = Array.isArray(v) ? v.filter(x => x !== '') : v
    } else if (!isEmpty(draft.values[f.name])) {
      out[f.name] = []
    }
  }
  return out
}

// Required fields left empty, by label.
export function missingRequired(draft: Draft, form: Record<string, FormValue>) {
  return (draft.view?.fields ?? [])
    .filter(f => f.required && !isReadOnly(f) && isEmpty(form[f.name]))
    .map(f => f.label || f.name)
}

// What the model is told once the user saved, so its next turn knows.
export function savedContext(draft: Draft, recordID: string, values: Record<string, FormValue>) {
  const verb = draft.recordID ? 'updated' : 'created'
  const module = draft.view?.module.name ?? draft.moduleID
  return `The user reviewed the draft and ${verb} record ${recordID} in ${module} (namespace ${draft.namespaceID}) with these values: ${JSON.stringify(values)}.`
}
