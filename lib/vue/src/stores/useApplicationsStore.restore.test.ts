import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { createApp } from 'vue'
import { useApplicationsStore } from './useApplicationsStore'

// A restored application has to come back into the app selector. Deleting drops
// it from the local list, so restoring has to put it back — otherwise the menu
// keeps saying the app is gone until the next full page load.

function withAPI(api: Record<string, unknown>) {
  const app = createApp({})
  app.provide('$SystemAPI', api)
  const pinia = createPinia()
  pinia._a = app
  app.use(pinia)
  setActivePinia(pinia)
  return useApplicationsStore()
}

const APP = { applicationID: 'A1', name: 'Reports', unify: { listed: true } }

describe('useApplicationsStore.restore', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('undeletes, then puts the application back in the list', async () => {
    const applicationUndelete = vi.fn().mockResolvedValue({})
    const applicationRead = vi.fn().mockResolvedValue(APP)
    const store = withAPI({
      applicationUndelete,
      applicationRead,
      applicationDelete: vi.fn().mockResolvedValue({}),
      applicationList: vi.fn().mockResolvedValue({ set: [] }),
    })

    store.apps = [APP]
    await store.delete('A1')
    expect(store.apps).toEqual([])

    await store.restore('A1')

    expect(applicationUndelete).toHaveBeenCalledWith({ applicationID: 'A1' })
    expect(store.apps).toEqual([APP])
  })

  it('replaces rather than duplicates an application still in the list', async () => {
    const store = withAPI({
      applicationUndelete: vi.fn().mockResolvedValue({}),
      applicationRead: vi.fn().mockResolvedValue({ ...APP, name: 'Reports (restored)' }),
    })

    store.apps = [APP]
    await store.restore('A1')

    expect(store.apps).toHaveLength(1)
    expect(store.apps[0].name).toBe('Reports (restored)')
  })

  it('leaves the list alone when the undelete is refused', async () => {
    const store = withAPI({
      applicationUndelete: vi.fn().mockRejectedValue(new Error('not allowed')),
      applicationRead: vi.fn(),
    })

    store.apps = []
    await expect(store.restore('A1')).rejects.toThrow('not allowed')
    expect(store.apps).toEqual([])
  })
})
