import { expect } from 'chai'
import { passwordPolicy, validatePassword, resetTokenFromLocation } from './c311-identity'

describe('C311 public identity helpers', () => {
  it('exposes the contract password policy and validates character classes', () => {
    expect(passwordPolicy.minLength).to.equal(12)
    expect(validatePassword('short')).to.include('too-short')
    expect(validatePassword('abcdefghijkl')).to.include('character-classes')
    expect(validatePassword('ValidPassword1!')).to.deep.equal([])
  })

  it('reads reset tokens once and removes them from the address bar', () => {
    const originalHistory = (globalThis as any).history
    const history = { replaceState: ((_state: unknown, _title: string, url?: string | URL | null) => { replaced = String(url || '') }) as History['replaceState'] } as unknown as History
    let replaced = ''
    Object.defineProperty(globalThis, 'history', { configurable: true, value: history })
    const token = resetTokenFromLocation({ search: '?token=ephemeral-token&email=example%40example.test', pathname: '/c311/reset-password' } as Location)
    expect(token).to.equal('ephemeral-token')
    expect(replaced).to.equal('/c311/reset-password?email=example%40example.test')
    Object.defineProperty(globalThis, 'history', { configurable: true, value: originalHistory })
  })
})
