import { describe, it, expect, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { createApp } from 'vue'
import { compose } from '@planetcrust/human-js'
import { useModuleStore } from './useModuleStore'
import { useNamespaceStore } from './useNamespaceStore'

// A record answers for its namespace through its module, and the module only
// knows one if it was built with it. Constraint matching on a compose event
// reads exactly that, and throws when it is missing.

const NS = { namespaceID: '1', slug: 'sandbox', name: 'Sandbox', enabled: true }
const MODULE = {
  moduleID: '2',
  namespaceID: '1',
  handle: 'contact',
  name: 'Contact',
  fields: [{ name: 'email', kind: 'String' }],
}

function storeWith(namespaces: unknown[]) {
  const app = createApp({})
  app.provide('$ComposeAPI', {
    moduleList: vi.fn().mockResolvedValue({ set: [MODULE] }),
    namespaceList: vi.fn().mockResolvedValue({ set: namespaces }),
  })
  const pinia = createPinia()
  ;(pinia as any)._a = app
  app.use(pinia)
  setActivePinia(pinia)

  const namespaceStore = useNamespaceStore()
  namespaceStore.set.splice(
    0,
    namespaceStore.set.length,
    ...namespaces.map(n => new compose.Namespace(n as object)),
  )

  return useModuleStore()
}

describe('useModuleStore namespace', () => {
  it('builds a module with the namespace its own store holds', async () => {
    const store = storeWith([NS])
    await store.load({ namespaceID: '1' })

    const module = store.getByID('2')
    expect(module.namespace).toBeDefined()
    expect(module.namespace.slug).toBe('sandbox')
  })

  it('lets a record answer for its namespace', async () => {
    const store = storeWith([NS])
    await store.load({ namespaceID: '1' })

    const record = new compose.Record(store.getByID('2'), { recordID: '3' })
    expect(record.namespace.slug).toBe('sandbox')
  })

  it('leaves the namespace unset rather than guessing when it is not cached', async () => {
    const store = storeWith([])
    await store.load({ namespaceID: '1' })

    expect(store.getByID('2').namespace).toBeUndefined()
  })
})
