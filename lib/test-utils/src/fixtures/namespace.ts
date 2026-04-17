import { compose } from '@planetcrust/human-js'

let _seq = 10001
const nextID = () => String(_seq++)

export interface NamespaceOverrides {
  namespaceID?: string
  name?: string
  slug?: string
  enabled?: boolean
}

export function makeNamespace(overrides: NamespaceOverrides = {}): compose.Namespace {
  return new compose.Namespace({
    namespaceID: overrides.namespaceID ?? nextID(),
    name: overrides.name ?? 'Test Namespace',
    slug: overrides.slug ?? 'test-namespace',
    enabled: overrides.enabled ?? true,
  })
}
