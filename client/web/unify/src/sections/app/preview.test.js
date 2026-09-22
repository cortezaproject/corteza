import { describe, expect, it } from 'vitest'
import { previewsSwitchedOffApp } from './preview'

const custom = { enabled: false, unify: { kind: 'custom' }, canManageSourceOnApplication: true }

describe('switched-off custom app preview', () => {
  it('opens for whoever may change its page', () => {
    expect(previewsSwitchedOffApp(custom)).toBe(true)
  })

  it('stays shut for everybody else', () => {
    expect(previewsSwitchedOffApp({ ...custom, canManageSourceOnApplication: false })).toBe(false)
  })

  it('is only for a switched-off custom app', () => {
    expect(previewsSwitchedOffApp({ ...custom, enabled: true })).toBe(false)
    expect(previewsSwitchedOffApp({ ...custom, unify: { kind: '' } })).toBe(false)
    expect(previewsSwitchedOffApp(undefined)).toBe(false)
  })
})
