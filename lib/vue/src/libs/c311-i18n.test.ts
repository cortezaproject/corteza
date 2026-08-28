import { expect } from 'chai'
import { C311_LOCALES, readC311Locale, persistC311Locale } from './c311-i18n'

describe('C311 locale persistence', () => {
  const storage = new Map<string, string>()
  const localStorageStub = {
    getItem: (key: string) => storage.get(key) || null,
    setItem: (key: string, value: string) => storage.set(key, value),
  }

  beforeEach(() => {
    storage.clear()
    Object.defineProperty(globalThis, 'localStorage', { configurable: true, value: localStorageStub })
  })

  after(() => {
    delete (globalThis as { localStorage?: unknown }).localStorage
  })

  it('accepts supported locales and ignores unsupported values', () => {
    expect(C311_LOCALES).to.deep.equal(['en', 'es', 'vi'])
    persistC311Locale('es', 'actor-fixture-001')
    expect(readC311Locale('actor-fixture-001')).to.equal('es')
    storage.set('c311.locale.actor-fixture-001', 'fr')
    expect(readC311Locale('actor-fixture-001')).to.equal(null)
  })
})
