import { compose } from '@planetcrust/human-js'
import { makeModule } from './module'
import type { ModuleOverrides } from './module'

let _seq = 30001
const nextID = () => String(_seq++)

export function makeRecord(
  values: Record<string, string | string[]> = {},
  moduleOrOverrides?: compose.Module | ModuleOverrides,
): compose.Record {
  let mod: compose.Module
  if (moduleOrOverrides instanceof compose.Module) {
    mod = moduleOrOverrides
  } else {
    mod = makeModule({
      fields: Object.keys(values).map(name => ({ name, kind: 'String' })),
      ...moduleOrOverrides,
    })
  }
  return new compose.Record(mod, { values })
}

export interface RawRecordOverrides {
  recordID?: string
  moduleID?: string
  namespaceID?: string
  [key: string]: unknown
}

export function makeRawRecord(
  values: Array<{ name: string; value?: string }> = [],
  overrides: RawRecordOverrides = {},
): Record<string, unknown> {
  return {
    recordID: overrides.recordID ?? nextID(),
    moduleID: overrides.moduleID ?? '20001',
    namespaceID: overrides.namespaceID ?? '10001',
    values,
    ...overrides,
  }
}
