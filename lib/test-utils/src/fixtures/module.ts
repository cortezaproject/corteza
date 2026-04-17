import { compose } from '@planetcrust/human-js'

let _seq = 20001
const nextID = () => String(_seq++)

export interface FieldOverride {
  name: string
  kind?: string
  label?: string
  isMulti?: boolean
  isRequired?: boolean
  options?: Record<string, unknown>
}

export interface ModuleOverrides {
  moduleID?: string
  namespaceID?: string
  name?: string
  handle?: string
  fields?: FieldOverride[]
}

export function makeModule(overrides: ModuleOverrides = {}): compose.Module {
  const fields: FieldOverride[] = overrides.fields ?? [{ name: 'title', kind: 'String' }]

  return new compose.Module({
    moduleID: overrides.moduleID ?? nextID(),
    namespaceID: overrides.namespaceID ?? '10001',
    name: overrides.name ?? 'Test Module',
    handle: overrides.handle ?? 'test-module',
    fields: fields.map((f, i) => ({
      moduleFieldID: String(i + 1),
      kind: 'String',
      label: f.name,
      isMulti: false,
      isRequired: false,
      options: {},
      ...f,
    })),
  })
}
