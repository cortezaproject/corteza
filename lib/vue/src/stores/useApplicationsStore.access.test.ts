import { describe, it, expect, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { createApp } from 'vue'
import { useApplicationsStore } from './useApplicationsStore'

// Section entry is decided from this store, so its access answer has to fail
// closed. isPathEnabled, next door, deliberately fails open — it answers a
// different question (is this registered app switched on) where "no idea"
// genuinely means "carry on".

function withAPI(api: Record<string, unknown>) {
  const app = createApp({})
  app.provide('$SystemAPI', api)
  const pinia = createPinia()
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ;(pinia as any)._a = app
  app.use(pinia)
  setActivePinia(pinia)
  return useApplicationsStore()
}

const listing = (set: unknown[]) => withAPI({ applicationList: vi.fn().mockResolvedValue({ set }) })

const ADMIN = {
  applicationID: 'A1',
  name: 'Admin Area',
  enabled: true,
  unify: { listed: true, url: 'admin/' },
}

describe('useApplicationsStore access', () => {
  it('allows an application the payload says the user may access', async () => {
    const store = listing([{ ...ADMIN, canAccessApplication: true }])
    await store.ready()
    expect(store.canAccessApp('admin/')).toBe(true)
  })

  it('refuses one the payload says the user may not access', async () => {
    const store = listing([{ ...ADMIN, canAccessApplication: false }])
    await store.ready()
    expect(store.canAccessApp('admin/')).toBe(false)
  })

  it('refuses an application the list does not carry', async () => {
    const store = listing([{ ...ADMIN, canAccessApplication: true }])
    await store.ready()
    expect(store.canAccessApp('compose/')).toBe(false)
  })

  it('refuses before the list has been fetched', () => {
    const store = listing([{ ...ADMIN, canAccessApplication: true }])
    expect(store.canAccessApp('admin/')).toBe(false)
  })

  it('refuses when the list could not be fetched at all', async () => {
    const store = withAPI({
      applicationList: vi.fn().mockRejectedValue(new Error('not allowed to search')),
    })
    await store.ready()
    expect(store.error).toBeTruthy()
    expect(store.canAccessApp('admin/')).toBe(false)
  })

  it('matches the url however its slashes are written', async () => {
    const store = listing([
      { ...ADMIN, unify: { listed: true, url: '/Admin/' }, canAccessApplication: true },
    ])
    await store.ready()
    expect(store.canAccessApp('admin')).toBe(true)
    expect(store.canAccessApp('/admin/')).toBe(true)
  })

  it('serves one fetch to every caller that asks for the list', async () => {
    const applicationList = vi.fn().mockResolvedValue({ set: [ADMIN] })
    const store = withAPI({ applicationList })
    await Promise.all([store.ready(), store.ready(), store.ready()])
    expect(applicationList).toHaveBeenCalledTimes(1)
  })

  it('still fails open on enabled — a different question', async () => {
    const store = listing([{ ...ADMIN, canAccessApplication: false }])
    await store.ready()
    // Unregistered path: nothing to say it is switched off.
    expect(store.isPathEnabled('/nowhere')).toBe(true)
    // Registered and switched on, whatever the access answer is.
    expect(store.isPathEnabled('/admin/users')).toBe(true)
  })
})
