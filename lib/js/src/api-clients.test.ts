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

  it('carries validation details through, so a caller can place them by field', async () => {
    respondWith({
      error: {
        message: '2 issue(s) found',
        details: [
          { kind: 'empty', message: 'This field is required', meta: { field: 'title' } },
          { kind: 'invalidValue', message: 'Invalid field value', meta: { field: 'email' } },
        ],
      },
    })

    try {
      await new Compose({ baseURL: '/' }).namespaceList({})
      expect.fail('expected a rejection')
    } catch (e) {
      // The summary message names no field, so without the details a caller can
      // only show '2 issue(s) found' and drop what the server actually said.
      const details = (e as Error & { details?: Array<Record<string, unknown>> }).details
      expect(details).to.be.an('array').with.lengthOf(2)
      expect(details?.[0]).to.deep.include({ kind: 'empty', message: 'This field is required' })
      expect(details?.[0]?.meta).to.deep.equal({ field: 'title' })
      expect(details?.[1]?.meta).to.deep.equal({ field: 'email' })
    }
  })

  it('leaves details unset when the error carries none', async () => {
    respondWith({ error: { message: 'module does not exist' } })

    try {
      await new Compose({ baseURL: '/' }).namespaceList({})
      expect.fail('expected a rejection')
    } catch (e) {
      expect((e as Error & { details?: unknown }).details).to.equal(undefined)
      expect((e as Error).message).to.equal('module does not exist')
    }
  })

  it('still resolves a successful response unchanged', async () => {
    respondWith({ response: { set: [1, 2] } })

    const out = await new Compose({ baseURL: '/' }).namespaceList({})
    expect(out).to.deep.equal({ set: [1, 2] })
  })
})
