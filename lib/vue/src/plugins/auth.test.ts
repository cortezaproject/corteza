import { expect } from 'chai'
import sinon from 'sinon'
import axios from 'axios'
import { Auth } from './auth'

class MockStorage implements Storage {
  private map = new Map<string, string>()

  get length (): number { return this.map.size }
  clear (): void { this.map.clear() }
  getItem (key: string): string | null { return this.map.has(key) ? this.map.get(key) as string : null }
  key (index: number): string | null { return Array.from(this.map.keys())[index] ?? null }
  removeItem (key: string): void { this.map.delete(key) }
  setItem (key: string, value: string): void { this.map.set(key, value) }
}

interface Navigation { type: 'assign' | 'replace'; url: string }

class MockLocation {
  href: string
  navigations: Navigation[] = []

  constructor (href: string) {
    this.href = href
  }

  // Model real browser behavior: assign/replace record the pending navigation
  // but do NOT update href synchronously. href is only updated via the actual
  // navigation (which we don't simulate) or via history.replaceState below.
  assign (url: string): void { this.navigations.push({ type: 'assign', url }) }
  replace (url: string): void { this.navigations.push({ type: 'replace', url }) }
  toString (): string { return this.href }

  // The final navigation target the browser would resolve to once the current
  // task finishes — successive assign/replace calls in the same task collapse
  // into the most recent one.
  get pendingNavigation (): string | null {
    if (this.navigations.length === 0) return null
    return this.navigations[this.navigations.length - 1].url
  }
}

const tokenResponse = {
  aud: 'test-aud',
  sub: '12345',
  scope: 'profile api',
  access_token: 'access-token',
  refresh_token: 'refresh-token',
  expires_in: 3600,
  name: 'Test User',
  email: 'test@example.com',
}

describe(__filename, () => {
  let axiosCreateStub: sinon.SinonStub
  let originalWindow: any

  beforeEach(() => {
    // Stub axios so token exchange doesn't hit the network.
    axiosCreateStub = sinon.stub(axios, 'create').returns({
      post: () => Promise.resolve({ data: tokenResponse }),
      get: () => Promise.resolve({ data: tokenResponse }),
    } as any)

    originalWindow = (globalThis as any).window
  })

  afterEach(() => {
    axiosCreateStub.restore()
    ;(globalThis as any).window = originalWindow
  })

  function setupWindow (mockSessionStorage: MockStorage, mockLocation: MockLocation): void {
    ;(globalThis as any).window = {
      sessionStorage: mockSessionStorage,
      // history.replaceState would synchronously update the URL bar.
      // We model that so any fix using replaceState is observable in tests.
      history: {
        replaceState: (_state: unknown, _title: string, url: string) => {
          mockLocation.href = url
        },
      },
      setTimeout: () => 0 as unknown as number,
      clearTimeout: () => undefined,
    }
  }

  function makeAuth (opts: {
    sessionStorage: MockStorage
    location: MockLocation
    entrypointURL: string
  }): Auth {
    return new Auth({
      app: 'compose',
      verbose: false,
      cortezaAuthURL: 'https://example.com/auth',
      callbackURL: 'https://example.com/compose/auth/callback',
      location: opts.location as unknown as Location,
      sessionStorage: opts.sessionStorage,
      entrypointURL: opts.entrypointURL,
      refreshFactor: 0.75,
      registerEventListener: () => undefined,
    })
  }

  describe('duplicate-tab detection', () => {
    it('returns false on the first call and true on subsequent calls', () => {
      const ss = new MockStorage()
      const loc = new MockLocation('https://example.com/compose/')
      setupWindow(ss, loc)

      const auth = makeAuth({ sessionStorage: ss, location: loc, entrypointURL: loc.href })

      expect(auth.handleStateManagement()).to.equal(false)
      expect(auth.handleStateManagement()).to.equal(true)
    })

    it('handle() throws Unauthenticated when the dup flag is already present in sessionStorage', async () => {
      // Simulate sessionStorage cloned from a sibling tab: it already carries
      // the final-state flag and a refresh token. This is what a right-click
      // "open in new tab" produces.
      const ss = new MockStorage()
      ss.setItem('auth.state.final', String(Date.now()))
      ss.setItem('auth.refresh-token', 'cloned-refresh-token')

      const recordURL = 'https://example.com/compose/ns/crm/module/foo/record/123'
      const loc = new MockLocation(recordURL)
      setupWindow(ss, loc)

      const auth = makeAuth({ sessionStorage: ss, location: loc, entrypointURL: recordURL })

      let caught: Error | null = null
      try {
        await auth.handle()
      } catch (err) {
        caught = err as Error
      }

      expect(caught).to.be.instanceOf(Error)
      expect(caught?.message).to.equal('Unauthenticated')
    })
  })

  describe('callback handling after dup-tab redirect', () => {
    // This simulates the full real-world flow the user reported:
    //   1. Tab A is authenticated.
    //   2. User middle-clicks a record link → tab B opens at the record URL.
    //      sessionStorage is cloned (auth.state.final + refresh token present).
    //   3. Tab B's handle() throws Unauthenticated → startAuthenticationFlow
    //      saves `auth.state.<state>.location` = recordURL and redirects to OAuth.
    //   4. OAuth bounces back to /compose/auth/callback?code=…&state=…
    //   5. New handle() pass → handleCallbackRoute looks up the saved location
    //      and assigns it.
    //   6. **BUT** the webapp's .then() tail in client/web/compose/src/app.js
    //      reads window.location.href (still the callback URL) and calls
    //      location.replace(<callback URL without code>), clobbering the
    //      pending assign.
    //
    // The assertion below is: after the whole sequence, the final navigation
    // target is the saved record URL — not the callback URL.

    it('the final navigation after callback handling is the saved record URL', async () => {
      const recordURL = 'https://example.com/compose/ns/crm/module/foo/record/123'
      const state = 'state-uid-abc'
      const callbackEntrypoint = `https://example.com/compose/auth/callback?code=auth-code&state=${state}`

      // sessionStorage as it would look in tab B after startAuthenticationFlow
      // has run (in the previous load) and OAuth has redirected back.
      const ss = new MockStorage()
      ss.setItem('auth.refresh-token', 'cloned-refresh-token')
      ss.setItem(`auth.state.${state}.location`, recordURL)

      const loc = new MockLocation(callbackEntrypoint)
      setupWindow(ss, loc)

      const auth = makeAuth({ sessionStorage: ss, location: loc, entrypointURL: callbackEntrypoint })

      await auth.handle()

      // Simulate the cleanup tail that lives at the end of compose/src/app.js
      // (and admin/, one/, workflow/, etc.) — this is the racy code we want
      // to either remove or render harmless.
      const cur = new URL(loc.href)
      if (cur.searchParams.get('code')) {
        cur.searchParams.delete('code')
        loc.replace(cur.toString())
      }

      expect(loc.pendingNavigation).to.equal(recordURL)
    })

    it('does not navigate to a URL containing the OAuth "code" param', async () => {
      // Defensive: even if finalLocation somehow carries `code`, the URL we
      // navigate the user to must never expose it (it's a one-shot secret
      // and leaving it in the address bar is a code smell at best).
      const state = 'state-uid-xyz'
      // Saved location accidentally has a `code` query param of its own.
      const savedURL = 'https://example.com/compose/ns/crm/page/p1?code=something-else'
      const callbackEntrypoint = `https://example.com/compose/auth/callback?code=auth-code&state=${state}`

      const ss = new MockStorage()
      ss.setItem('auth.refresh-token', 'cloned-refresh-token')
      ss.setItem(`auth.state.${state}.location`, savedURL)

      const loc = new MockLocation(callbackEntrypoint)
      setupWindow(ss, loc)

      const auth = makeAuth({ sessionStorage: ss, location: loc, entrypointURL: callbackEntrypoint })

      await auth.handle()

      expect(loc.pendingNavigation).to.not.be.null
      const navigated = new URL(loc.pendingNavigation as string)
      expect(navigated.searchParams.get('code')).to.equal(null)
    })
  })
})
