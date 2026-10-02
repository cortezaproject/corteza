import { describe, expect, it } from 'vitest'
import { trimServer } from './trimServer'

describe('trimServer', () => {
  const padded = {
    host: '  smtp.example.com ',
    port: 587,
    user: ' mailer ',
    pass: ' keep spaces ',
    from: ' noreply@example.com ',
    tlsInsecure: false,
    tlsServerName: ' smtp.example.com ',
  }

  it('trims the host and addresses', () => {
    const server = trimServer(padded)
    expect(server.host).toBe('smtp.example.com')
    expect(server.user).toBe('mailer')
    expect(server.from).toBe('noreply@example.com')
    expect(server.tlsServerName).toBe('smtp.example.com')
  })

  it('keeps the password and non-string values as they are', () => {
    const server = trimServer(padded)
    expect(server.pass).toBe(' keep spaces ')
    expect(server.port).toBe(587)
    expect(server.tlsInsecure).toBe(false)
  })

  it('does not mutate the input', () => {
    trimServer(padded)
    expect(padded.host).toBe('  smtp.example.com ')
  })
})
