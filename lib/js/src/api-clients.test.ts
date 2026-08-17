import { expect } from 'chai'
import axios from 'axios'
import sinon from 'ts-sinon'
import { Compose, System } from './api-clients'

// The API reports failure in the body of a 200. Every client is generated from
// one template, so what stdResolve rejects with is the same question for all of
// them — and a rejected plain object stringifies to '[object Object]' wherever
// a caller shows it, which is what reached the chart builder's preview.
describe('api-clients stdResolve', () => {
  let restore: () => void

  const respondWith = (data: unknown) => {
    const stub = sinon.stub(axios, 'create').returns({
      request: () => Promise.resolve({ data }),
    } as never)
    restore = () => stub.restore()
  }

  afterEach(() => restore?.())

  it('rejects an { message } error as an Error carrying that message', async () => {
    respondWith({ error: { message: 'found an illegal token' } })

    try {
      await new Compose({ baseURL: '/' }).namespaceList({})
      expect.fail('expected a rejection')
    } catch (e) {
      expect(e).to.be.instanceOf(Error)
      expect((e as Error).message).to.equal('found an illegal token')
      expect(String(e)).to.not.contain('[object Object]')
    }
  })

  it('rejects a bare string error as an Error too', async () => {
    respondWith({ error: 'not allowed to read this' })

    try {
      await new Compose({ baseURL: '/' }).namespaceList({})
      expect.fail('expected a rejection')
    } catch (e) {
      expect(e).to.be.instanceOf(Error)
      expect((e as Error).message).to.equal('not allowed to read this')
    }
  })

  it('names the failure even when the error object carries no message', async () => {
    respondWith({ error: { code: 42 } })

    try {
      await new Compose({ baseURL: '/' }).namespaceList({})
      expect.fail('expected a rejection')
    } catch (e) {
      expect(e).to.be.instanceOf(Error)
      expect(String(e)).to.not.contain('[object Object]')
    }
  })

  it('applies to every generated client, not just compose', async () => {
    respondWith({ error: { message: 'nope' } })

    try {
      await new System({ baseURL: '/' }).userList({})
      expect.fail('expected a rejection')
    } catch (e) {
      expect(e).to.be.instanceOf(Error)
      expect((e as Error).message).to.equal('nope')
    }
  })

  it('still resolves a successful response unchanged', async () => {
    respondWith({ response: { set: [1, 2] } })

    const out = await new Compose({ baseURL: '/' }).namespaceList({})
    expect(out).to.deep.equal({ set: [1, 2] })
  })
})
