import { expect } from 'chai'
import { readFileSync } from 'node:fs'
import { C311ApiError } from './errors'
import { cloneFixtureSet, createDefaultFixtureSet } from './fixtures'
import { MockC311Provider } from './mock-provider'
import { C311FetchTransport, C311HttpProvider, type C311Provider, type C311TransportRequest } from './provider'
import { C311_TIMEZONE, formatC311DateTime } from './time'
import type { PortalServiceRequestCreate, ReportDefinition, ServiceRequestCreate } from './types'
import { APPLICATION_ROLES } from './enums'

async function expectError (action: () => Promise<unknown>, code: string): Promise<C311ApiError> {
  try {
    await action()
    throw new Error(`expected ${code}`)
  } catch (error) {
    expect(error).to.be.instanceOf(C311ApiError)
    expect((error as C311ApiError).code).to.equal(code)
    return error as C311ApiError
  }
}

describe('City 311 frontend contract', () => {
  it('returns contract-v1 fixture data with wire-compatible fields', async () => {
    const fixtures = createDefaultFixtureSet()
    const page = await new MockC311Provider().listPortalRequests()

    expect(fixtures.fixture_id).to.equal('contract-v1')
    expect(fixtures.contract_version).to.equal('1.0.0')
    expect(page.items).to.have.length(1)
    expect(page.items[0].request_id).to.equal('request-fixture-001')
    expect(page.items[0].owning_department).to.equal('STREETS')
    expect(page.items[0]).to.not.have.property('primary_requester')
  })

  it('keeps the empty response shape stable', async () => {
    const page = await new MockC311Provider({ scenario: 'empty' }).listStaffRequests()

    expect(page.items).to.deep.equal([])
    expect(page.total_count).to.equal(0)
    expect(page.next_page_token).to.equal(null)
    expect(page.applied_filters).to.deep.equal({})
  })

  it('keeps anonymous lookup privacy and reopen response shapes', async () => {
    const provider = new MockC311Provider()
    expect(await provider.getPublicStatus({ request_number: 'SR-2026-99999', email: 'unknown@example.test' })).to.deep.equal({ request_detail: null })
    expect(await provider.getPublicStatus({ request_number: 'SR-2026-00001', email: 'wrong@example.test' })).to.deep.equal({ request_detail: null })
    expect((await provider.getPublicStatus({ request_number: 'SR-2026-00001', email: 'alex@example.test' })).request_detail?.request_number).to.equal('SR-2026-00001')
    expect((await provider.getPublicStatus({ request_number: 'SR-2026-00001', email: ' ALEX@EXAMPLE.TEST ' })).request_detail?.request_number).to.equal('SR-2026-00001')
    expect(await provider.reopenPortalRequest('request-fixture-001', 'fixture reason')).to.deep.equal({
      request_id: 'request-fixture-001',
      status: 'PENDING_APPROVAL',
    })
  })

  it('provides complete role and session-expiry fixtures for route checks', () => {
    const fixtures = createDefaultFixtureSet()
    const roles = Object.keys(fixtures.role_fixtures)
    expect(roles).to.have.members(APPLICATION_ROLES)

    for (const role of roles) {
      const fixture = fixtures.role_fixtures[role as keyof typeof fixtures.role_fixtures]
      expect(fixture.denied_route).to.be.a('string')
      expect(fixture.denied_capability).to.be.a('string')
      expect(fixture.denied_scope).to.be.a('string')
      expect(fixture.expired_session.expires_at).to.equal('2026-01-15T14:00:00.000Z')
      if (role === 'public_visitor') expect(fixture.session.authenticated).to.equal(false)
      else {
        expect(fixture.session.actor?.available_routes).to.not.include(fixture.denied_route)
        expect(fixture.session.actor?.capabilities).to.not.include(fixture.denied_capability)
        expect(fixture.session.actor?.scopes).to.not.include(fixture.denied_scope)
        expect(fixture.session.actor?.application_roles).to.include(role)
      }
    }
  })

  it('selects every role and expired session without changing the fixture set', async () => {
    for (const role of APPLICATION_ROLES) {
      const current = await new MockC311Provider({ role }).getSession()
      const expired = await new MockC311Provider({ role, sessionVariant: 'expired' }).getSession()
      expect(current).to.deep.equal(createDefaultFixtureSet().role_fixtures[role].session)
      expect(expired.expires_at).to.equal('2026-01-15T14:00:00.000Z')
    }
  })

  it('keeps all role fixture values inside the frozen contract vocabularies', () => {
    const contract = JSON.parse(readFileSync(new URL('../../../../server/compose/types/city311/contract.json', import.meta.url), 'utf8')) as {
      enums: Record<string, string[]>
    }
    const fixtures = createDefaultFixtureSet()
    for (const role of APPLICATION_ROLES) {
      const fixture = fixtures.role_fixtures[role]
      const actor = fixture.session.actor
      const assertKnown = (kind: string, value: string) => expect(contract.enums[kind]).to.include(value)
      assertKnown('route', fixture.denied_route)
      assertKnown('capability', fixture.denied_capability)
      assertKnown('oauth_scope', fixture.denied_scope)
      if (!actor) continue
      actor.application_roles.forEach(value => assertKnown('application_role', value))
      actor.department_codes.forEach(value => assertKnown('department_code', value))
      actor.district_codes.forEach(value => assertKnown('district_code', value))
      actor.capabilities.forEach(value => assertKnown('capability', value))
      actor.scopes.forEach(value => assertKnown('oauth_scope', value))
      actor.available_routes.forEach(value => assertKnown('route', value))
      expect(actor.available_routes).to.not.include(fixture.denied_route)
      expect(actor.capabilities).to.not.include(fixture.denied_capability)
      expect(actor.scopes).to.not.include(fixture.denied_scope)
    }
  })

  it('exposes only endpoint-declared failures and models terminal operations in-band', async () => {
    const portalInput = {} as PortalServiceRequestCreate
    const scenarios: Array<[string, () => Promise<unknown>, string, boolean]> = [
      ['forbidden', () => new MockC311Provider({ scenario: 'forbidden' }).listStaffRequests(), 'FORBIDDEN', false],
      ['not-found', () => new MockC311Provider({ scenario: 'not-found' }).getStaffRequest('missing'), 'NOT_FOUND', false],
      ['validation', () => new MockC311Provider({ scenario: 'validation' }).submitPortalRequest(portalInput), 'VALIDATION_ERROR', false],
      ['retryable', () => new MockC311Provider({ scenario: 'retryable' }).geocode({ address: 'fixture' }), 'MAP_TEMPORARILY_UNAVAILABLE', true],
      ['version-conflict', () => new MockC311Provider({ scenario: 'version-conflict' }).updateDraft('draft-fixture-001', {}, { expectedVersion: 1 }), 'VERSION_CONFLICT', false],
    ]

    for (const [scenario, action, code, retryable] of scenarios) {
      const error = await expectError(action, code)
      expect(error.retryable).to.equal(retryable)
      if (scenario === 'retryable') expect(error.retryAfter).to.equal('30')
      if (scenario === 'version-conflict') expect(error.currentVersion).to.equal(2)
    }

    expect((await new MockC311Provider({ scenario: 'version-conflict' }).getSession()).authenticated).to.equal(true)
    const terminal = await new MockC311Provider({ scenario: 'terminal' }).getOperation('operation-fixture-terminal')
    expect(terminal.status).to.equal('FAILED')
    expect(terminal.error?.error).to.equal('OPERATION_FAILED')
  })

  it('does not mutate caller fixtures while exercising draft flows', async () => {
    const fixtures = createDefaultFixtureSet()
    const before = JSON.stringify(fixtures)
    const provider = new MockC311Provider({ fixtures })

    await provider.createDraft({ summary: 'temporary draft' })
    await provider.updateDraft('draft-fixture-001', { summary: 'temporary update' })
    await provider.deleteDraft('draft-fixture-001')

    expect(JSON.stringify(fixtures)).to.equal(before)
  })

  it('uses the frozen endpoint paths and concurrency headers', async () => {
    const requests: C311TransportRequest[] = []
    const transport = {
      request: async <T> (request: C311TransportRequest): Promise<T> => {
        requests.push(request)
        return {} as T
      },
    }
    const provider: C311Provider = new C311HttpProvider(transport)
    const input: ServiceRequestCreate = {
      summary: 'Example request',
      description: 'A sufficiently long fixture description.',
      service_type: 'GENERAL_INQUIRY',
      requester: { display_name: 'Example User', email: 'user@example.test' },
    }

    await provider.createServiceRequest(input, { idempotencyKey: 'fixture-key', expectedVersion: 3 })
    const portalInput: PortalServiceRequestCreate = {
      summary: input.summary,
      description: input.description,
      service_type: input.service_type,
      requester: input.requester,
      attachment_tokens: ['attachment-token-fixture-001'],
    }
    await provider.submitPortalRequest(portalInput, { idempotencyKey: 'portal-fixture-key' })
    await provider.getPublicStatus({ request_number: 'SR-2026-00001', email: 'alex@example.test' })
    await provider.reopenPortalRequest('request-fixture-001', 'fixture reason')
    await provider.downloadAttachment('attachment-fixture-001')
    await provider.uploadPortalAttachment({ file: 'ZmFrZQ==', filename: 'fixture.txt', media_type: 'text/plain' })
    await provider.exportReport('report-fixture-001')
    await provider.updateWorkflow('workflow-fixture-001', { name: 'Updated fixture workflow' }, { expectedVersion: 4 })
    await provider.deleteDraft('draft-fixture-001', { expectedVersion: 2 })

    expect(requests[0]).to.deep.include({
      method: 'POST',
      path: '/api/v1/service-requests',
      body: input,
      headers: { 'Idempotency-Key': 'fixture-key' },
    })
    expect(requests[1]).to.deep.include({
      method: 'POST',
      path: '/api/v1/portal/service-requests',
      body: portalInput,
      headers: { 'Idempotency-Key': 'portal-fixture-key' },
    })
    expect(requests[2]).to.deep.include({ method: 'POST', path: '/api/v1/public/service-request-status', acceptedStatuses: [404] })
    expect(requests[3]).to.deep.include({ method: 'POST', path: '/api/v1/portal/service-requests/request-fixture-001/reopen' })
    expect(requests[4]).to.deep.include({ method: 'GET', path: '/api/v1/attachments/attachment-fixture-001' })
    expect(requests[5].body).to.be.instanceOf(FormData)
    expect(requests[6]).to.deep.include({
      method: 'POST',
      path: '/api/v1/staff/reports/report-fixture-001/export',
      body: { format: 'CSV' },
    })
    expect(requests[7]).to.deep.include({
      method: 'PATCH',
      path: '/api/v1/admin/workflows/workflow-fixture-001',
      headers: { 'If-Match': '"4"' },
    })
    expect(requests[8]).to.deep.include({
      method: 'DELETE',
      path: '/api/v1/portal/service-request-drafts/draft-fixture-001',
      headers: { 'If-Match': '"2"' },
    })
  })

  it('uses the frozen account maintenance endpoints and payloads', async () => {
    const requests: C311TransportRequest[] = []
    const session = createDefaultFixtureSet().session
    const provider = new C311HttpProvider({
      request: async <T> (request: C311TransportRequest): Promise<T> => {
        requests.push(request)
        if (request.path === '/api/v1/account/login-identifier') return session as T
        if (request.path === '/api/v1/account/link/confirm') return session as T
        return undefined as T
      },
    })

    expect(await provider.changeLoginIdentifier({ current_password: 'Current-password-1!', login_identifier: 'alex.new' })).to.deep.equal(session)
    expect(await provider.changePassword({ current_password: 'Current-password-1!', new_password: 'New-password-2!' })).to.equal(undefined)
    expect(await provider.confirmAccountLink()).to.deep.equal(session)
    expect(requests).to.deep.equal([
      { method: 'POST', path: '/api/v1/account/login-identifier', body: { current_password: 'Current-password-1!', login_identifier: 'alex.new' } },
      { method: 'POST', path: '/api/v1/account/password', body: { current_password: 'Current-password-1!', new_password: 'New-password-2!' } },
      { method: 'POST', path: '/api/v1/account/link/confirm', body: {} },
    ])
  })

  it('keeps the current session and submitted values unchanged when account maintenance fails', async () => {
    const conflict = new MockC311Provider({ role: 'constituent', scenario: 'version-conflict' })
    const before = await conflict.getSession()
    await expectError(() => conflict.changeLoginIdentifier({ current_password: 'Current-password-1!', login_identifier: 'alex.conflict' }), 'VERSION_CONFLICT')
    expect(await conflict.getSession()).to.deep.equal(before)

    const invalid = new MockC311Provider({ role: 'constituent', scenario: 'validation' })
    await expectError(() => invalid.changePassword({ current_password: 'wrong', new_password: 'short' }), 'VALIDATION_ERROR')
    expect(await invalid.getSession()).to.deep.equal(await new MockC311Provider({ role: 'constituent' }).getSession())

    const profile = new MockC311Provider({ role: 'constituent' })
    await expectError(() => profile.updateProfile({ display_name: 'Updated' }), 'EXPECTED_VERSION_REQUIRED')
  })

  it('uses one-time opaque reset tokens and identical forgot-password responses', async () => {
    const privacyProvider = new MockC311Provider()
    const known = await privacyProvider.requestPasswordReset({ email: 'alex@example.test' })
    const unknown = await privacyProvider.requestPasswordReset({ email: 'unknown@example.test' })
    expect(unknown).to.deep.equal(known)

    const provider = new MockC311Provider()
    await provider.requestPasswordReset({ email: 'alex@example.test' })
    await provider.requestPasswordReset({ email: 'alex@example.test' })
    await expectError(() => provider.confirmPasswordReset({ token: 'reset-token-fixture-001', password: 'New-password-2!' }), 'INVALID_RESET_TOKEN')
    expect(await provider.confirmPasswordReset({ token: 'reset-token-fixture-002', password: 'New-password-2!' })).to.include({ message: 'Your password has been reset.' })
    await expectError(() => provider.confirmPasswordReset({ token: 'reset-token-fixture-002', password: 'New-password-2!' }), 'INVALID_RESET_TOKEN')
  })

  it('exposes and completes an explicit link confirmation operation', async () => {
    const provider = new MockC311Provider({ scenario: 'link-confirmation-required' })
    await provider.startFederatedSignIn('oidc')
    const pending = await provider.completeFederatedSignIn('oidc', { code: 'fixture-code' })
    expect(pending.outcome).to.equal('link_confirmation_required')
    const confirmed = await provider.confirmAccountLink()
    expect(confirmed.authenticated).to.equal(true)
  })

  it('keeps the account session unchanged when a pending federated link is cancelled', async () => {
    const provider = new MockC311Provider({ scenario: 'account-link-cancelled' })
    const before = await provider.getSession()
    await expectError(() => provider.startFederatedSignIn('saml'), 'FORBIDDEN')
    expect(await provider.getSession()).to.deep.equal(before)
  })

  it('nests request filters under the contract filters parameter', async () => {
    const requests: C311TransportRequest[] = []
    const provider = new C311HttpProvider({
      request: async <T> (request: C311TransportRequest): Promise<T> => {
        requests.push(request)
        return { items: [], next_page_token: null, total_count: 0, applied_filters: {}, sort: [] } as T
      },
    })

    await provider.listPortalRequests({
      page_size: 20,
      status: 'SUBMITTED',
      service_type: 'POTHOLE',
      filters: { category: 'RESIDENT' },
      sort: '-updated_at',
    })

    expect(requests[0].query).to.deep.equal({
      page_size: 20,
      filters: { category: 'RESIDENT', status: 'SUBMITTED', service_type: 'POTHOLE' },
      sort: '-updated_at',
    })
    expect(requests[0].query).to.not.have.property('status')
  })

  it('continues report lookup across opaque result pages', async () => {
    const requests: C311TransportRequest[] = []
    const report: ReportDefinition = {
      report_id: 'report-page-2',
      name: 'Second page report',
      entity: 'service_requests',
      columns: ['request_number'],
      filters: {},
      sort: [],
      version: 1,
      updated_at: '2026-01-15T15:00:00.000Z',
    }
    const provider = new C311HttpProvider({
      request: async <T> (request: C311TransportRequest): Promise<T> => {
        requests.push(request)
        const secondPage = request.query?.page_token === 'page-2'
        return {
          items: secondPage ? [report] : [],
          next_page_token: secondPage ? null : 'page-2',
          total_count: 1,
          applied_filters: {},
          sort: [],
        } as T
      },
    })

    expect(await provider.getReport('report-page-2')).to.deep.equal(report)
    expect(requests).to.have.length(2)
    expect(requests[1].query).to.deep.equal({ page_token: 'page-2' })
  })

  it('maps fetch responses, query arrays and contract errors', async () => {
    const calls: Array<{ url: string; init: RequestInit }> = []
    const okTransport = new C311FetchTransport({
      baseURL: 'https://fixture.example',
      fetch: async (url, init) => {
        calls.push({ url: String(url), init })
        return {
          ok: true,
          status: 204,
          headers: { get: () => null, forEach: () => {} },
        } as unknown as Response
      },
    })

    await okTransport.request({
      method: 'DELETE',
      path: '/api/v1/session',
      query: { sort: ['updated_at', 'desc'], page_size: 20, filters: { status: 'SUBMITTED' } },
    })
    expect(calls[0].url).to.equal('https://fixture.example/api/v1/session?sort=updated_at&sort=desc&page_size=20&filters=%7B%22status%22%3A%22SUBMITTED%22%7D')
    expect(calls[0].init.credentials).to.equal('include')

    const errorTransport = new C311FetchTransport({
      fetch: async () => ({
        ok: false,
        status: 422,
        headers: { get: () => 'application/json', forEach: () => {} },
        json: async () => ({ error: 'VALIDATION_ERROR', message: 'Invalid input.', retryable: false, errors: [{ field: '/summary', code: 'REQUIRED' }] }),
      } as unknown as Response),
    })
    const error = await expectError(() => errorTransport.request({ method: 'POST', path: '/api/v1/portal/service-requests', body: {} }), 'VALIDATION_ERROR')
    expect(error.status).to.equal(422)
    expect(error.error).to.equal('VALIDATION_ERROR')
    expect(error.fieldErrors[0].field).to.equal('/summary')

    const privacyTransport = new C311FetchTransport({
      fetch: async () => ({
        ok: false,
        status: 404,
        headers: { get: () => 'application/json', forEach: () => {} },
        json: async () => ({ request_detail: null }),
      } as unknown as Response),
    })
    const privacyResponse = await privacyTransport.request<{ request_detail: null }>({
      method: 'POST',
      path: '/api/v1/public/service-request-status',
      body: {},
      acceptedStatuses: [404],
    })
    expect(privacyResponse).to.deep.equal({ request_detail: null })

    const networkError = await expectError(() => new C311FetchTransport({ fetch: async () => { throw new Error('offline') } }).request({ method: 'GET', path: '/healthz' }), 'TEMPORARILY_UNAVAILABLE')
    expect(networkError.retryable).to.equal(true)

    const malformedTransport = new C311FetchTransport({
      fetch: async () => ({
        ok: false,
        status: 503,
        headers: { get: () => 'application/json', forEach: () => {} },
        text: async () => '{invalid',
      } as unknown as Response),
    })
    const malformed = await expectError(() => malformedTransport.request({ method: 'GET', path: '/healthz' }), 'OPERATION_FAILED')
    expect(malformed.status).to.equal(503)
    expect(malformed.retryable).to.equal(false)
  })

  it('normalizes anonymous HTTP misses to the same non-disclosing response', async () => {
    const provider = new C311HttpProvider({
      request: async <T> (): Promise<T> => ({ error: 'NOT_FOUND', message: 'not found' } as T),
    })

    expect(await provider.getPublicStatus({ request_number: 'SR-2026-99999', email: 'unknown@example.test' })).to.deep.equal({ request_detail: null })
  })

  it('formats the benchmark instant in the fixed timezone', () => {
    expect(C311_TIMEZONE).to.equal('America/New_York')
    expect(formatC311DateTime('2026-01-15T15:00:00.000Z', 'en-US')).to.equal('01/15/2026 10:00 AM EST')
    expect(formatC311DateTime('2026-07-15T15:00:00.000Z', 'en-US')).to.equal('07/15/2026 11:00 AM EDT')
    expect(formatC311DateTime('2026-07-15T15:00:00.000Z', 'es')).to.equal('07/15/2026 11:00 AM EDT')
    expect(formatC311DateTime('not-a-date')).to.equal('')
  })

  it('deep clones custom fixture sets', () => {
    const fixtures = createDefaultFixtureSet()
    const clone = cloneFixtureSet(fixtures)
    clone.requests[0].summary = 'changed in test'

    expect(fixtures.requests[0].summary).to.equal('Pothole on Example Street')
  })

  it('maps public identity operations to the frozen contract paths', async () => {
    const requests: C311TransportRequest[] = []
    const provider: any = new C311HttpProvider({
      request: async <T> (request: C311TransportRequest): Promise<T> => {
        requests.push(request)
        if (request.path === '/api/v1/public/branding') return { organisation_name: 'City 311' } as T
        if (request.path.includes('/content/')) return { content_key: 'HOME', body: '<p>Welcome</p>' } as T
        if (request.path.includes('/help/')) return { help_key: 'public.request.submit', language: 'EN', body: '<p>Help</p>', version: 1, updated_at: '2026-01-15T15:00:00.000Z' } as T
        if (request.path === '/api/v1/account/profile') return createDefaultFixtureSet().requests[0].primary_requester as T
        if (request.method === 'POST' && request.path === '/api/v1/session') return createDefaultFixtureSet().session as T
        return { accepted: true, message: 'accepted' } as T
      },
    })

    await provider.registerAccount({ display_name: 'Example', email: 'example@example.test', login_identifier: 'example', password: 'ValidPassword1!', preferred_language: 'EN' })
    await provider.requestPasswordReset({ email: 'example@example.test' })
    await provider.confirmPasswordReset({ token: 'ephemeral-token', password: 'ValidPassword1!' })
    await provider.startFederatedSignIn('oidc')
    await provider.completeFederatedSignIn('saml', { code: 'ephemeral-code', state: 'ephemeral-state' })
    await provider.getBranding()
    await provider.getPublicContent('HOME')
    await provider.getPublicHelp('public.request.submit', 'EN')
    await provider.getProfile()
    await provider.updateProfile({ display_name: 'Updated' }, { expectedVersion: 1 })
    await provider.updateLanguage('ES')
    await provider.changeLoginIdentifier({ current_password: 'Current-password-1!', login_identifier: 'updated.login' })
    await provider.changePassword({ current_password: 'Current-password-1!', new_password: 'New-password-2!' })

    expect(requests.map(request => `${request.method} ${request.path}`)).to.include.members([
      'POST /api/v1/accounts',
      'POST /api/v1/auth/password-reset/request',
      'POST /api/v1/auth/password-reset/confirm',
      'GET /api/v1/auth/oidc/start',
      'GET /api/v1/auth/saml/callback',
      'GET /api/v1/public/branding',
      'GET /api/v1/public/content/HOME',
      'GET /api/v1/public/help/public.request.submit',
      'GET /api/v1/account/profile',
      'PATCH /api/v1/account/profile',
      'PATCH /api/v1/preferences/language',
      'POST /api/v1/account/login-identifier',
      'POST /api/v1/account/password',
    ])
    expect(requests.find(request => request.path === '/api/v1/account/login-identifier')?.body).to.deep.equal({ current_password: 'Current-password-1!', login_identifier: 'updated.login' })
    expect(requests.find(request => request.path === '/api/v1/account/password')?.body).to.deep.equal({ current_password: 'Current-password-1!', new_password: 'New-password-2!' })
    expect(requests.find(request => request.path === '/api/v1/auth/saml/callback')?.query).to.deep.equal({ code: 'ephemeral-code', state: 'ephemeral-state' })
  })

  it('validates password policy without persisting credentials', () => {
    const fixture = createDefaultFixtureSet()
    expect(fixture).to.not.have.property('password')
    expect(fixture).to.not.have.property('reset_token')
    const provider: any = new MockC311Provider({ scenario: 'registration-validation' })
    return expectError(() => provider.registerAccount({ display_name: '', email: 'bad', login_identifier: 'x', password: 'short', preferred_language: 'EN' }), 'VALIDATION_ERROR')
  })

  it('models identity and public-content fixture scenarios with contract errors', async () => {
    const cases: Array<[string, string, number]> = [
      ['invalid-credentials', 'UNAUTHENTICATED', 401],
      ['expired-reset-token', 'EXPIRED_RESET_TOKEN', 422],
      ['invalid-reset-token', 'INVALID_RESET_TOKEN', 422],
      ['oidc-failure', 'TEMPORARILY_UNAVAILABLE', 503],
      ['saml-failure', 'FORBIDDEN', 403],
      ['identity-claims-failure', 'UNAUTHENTICATED', 401],
      ['branding-failure', 'TEMPORARILY_UNAVAILABLE', 503],
      ['content-loading-failure', 'TEMPORARILY_UNAVAILABLE', 503],
      ['help-loading-failure', 'TEMPORARILY_UNAVAILABLE', 503],
      ['account-loading', 'TEMPORARILY_UNAVAILABLE', 503],
    ]
    for (const [scenario, code, status] of cases) {
      const provider: any = new MockC311Provider({ scenario: scenario as any })
      const action = scenario === 'invalid-credentials'
        ? () => provider.signIn({ login_identifier: 'fixture', password: 'not-a-secret' })
        : scenario === 'expired-reset-token' || scenario === 'invalid-reset-token'
          ? () => provider.confirmPasswordReset({ token: 'ephemeral-token', password: 'ValidPassword1!' })
          : scenario === 'oidc-failure' || scenario === 'saml-failure'
            ? () => provider.startFederatedSignIn(scenario === 'oidc-failure' ? 'oidc' : 'saml')
            : scenario === 'identity-claims-failure'
              ? () => provider.completeFederatedSignIn('oidc', { code: 'fixture-code' })
            : scenario === 'branding-failure'
              ? () => provider.getBranding()
              : scenario === 'help-loading-failure'
                ? () => provider.getPublicHelp('public.request.submit', 'EN')
                : scenario === 'account-loading'
                  ? () => provider.getProfile()
                  : () => provider.getPublicContent('HOME')
      const error = await expectError(action, code)
      expect(error.status).to.equal(status)
    }
  })

  it('rejects public SAML callbacks in the mock provider', async () => {
    await expectError(() => new MockC311Provider().completeFederatedSignIn('saml'), 'FORBIDDEN')
  })

  it('keeps accounts and the current session unchanged when identity claims fail', async () => {
    const provider: any = new MockC311Provider({ scenario: 'identity-claims-failure', role: 'constituent' })
    const accountsBefore = JSON.parse(JSON.stringify(provider.fixtures.role_fixtures))
    const sessionBefore = await provider.getSession()

    await expectError(() => provider.completeFederatedSignIn('oidc', { code: 'fixture-code' }), 'UNAUTHENTICATED')

    expect(provider.fixtures.role_fixtures).to.deep.equal(accountsBefore)
    expect(await provider.getSession()).to.deep.equal(sessionBefore)
  })

  it('maps a cancelled federated callback to the generic unauthenticated result', async () => {
    const provider = new MockC311Provider()
    const error = await expectError(() => provider.completeFederatedSignIn('oidc', { error: 'access_denied' }), 'UNAUTHENTICATED')
    expect(error.status).to.equal(401)
    expect(error.message).to.equal('Federated sign-in was cancelled.')
  })

  it('supports an empty authenticated request catalogue fixture', async () => {
    const page = await new MockC311Provider({ scenario: 'empty-my-requests', role: 'constituent' }).listPortalRequests()
    expect(page.items).to.deep.equal([])
    expect(page.total_count).to.equal(0)
  })

  it('maps FE-03 submit, draft, and staff-assist operations to contract paths and headers', async () => {
    const requests: C311TransportRequest[] = []
    const provider = new C311HttpProvider({
      request: async <T> (request: C311TransportRequest): Promise<T> => {
        requests.push(request)
        if (request.path === '/api/v1/staff/service-requests') return { request: { request_id: 'staff-request' } } as T
        if (request.path.includes('/submit')) return { request_id: 'draft-request', request_number: 'SR-2026-00004', status: 'SUBMITTED', version: 2, created_at: '2026-01-15T15:00:00.000Z', links: { self: '/api/v1/service-requests/draft-request' } } as T
        if (request.method === 'DELETE') return undefined as T
        return { request_id: 'draft-request', status: 'DRAFT', version: 2 } as T
      },
    })
    const portalInput = {
      summary: 'Pothole near library',
      description: 'The road surface is damaged near the library entrance.',
      service_type: 'POTHOLE' as const,
      requester: { display_name: 'Alex Example', email: 'alex@example.test' },
      location: { address: '100 Example Street', latitude: 42.9001, longitude: -88.8801 },
      attachment_tokens: ['upload-00031'],
      custom_fields: { ward: 'NORTH' },
    }
    await provider.submitPortalRequest(portalInput, { idempotencyKey: 'fe03-submit-1' })
    await provider.createDraft(portalInput)
    await provider.getDraft('draft-request')
    await provider.updateDraft('draft-request', { summary: 'Updated summary' }, { expectedVersion: 2 })
    await provider.deleteDraft('draft-request', { expectedVersion: 2 })
    await provider.submitDraft('draft-request', { expectedVersion: 2 })
    await provider.createStaffServiceRequest({
      constituent: { constituent_id: 'constituent-fixture-001' },
      request: portalInput,
    })

    expect(requests.find(request => request.path === '/api/v1/portal/service-requests')?.headers).to.deep.equal({ 'Idempotency-Key': 'fe03-submit-1' })
    expect(requests.find(request => request.path === '/api/v1/portal/service-request-drafts/draft-request' && request.method === 'PATCH')?.headers).to.deep.equal({ 'If-Match': '"2"' })
    expect(requests.find(request => request.path === '/api/v1/portal/service-request-drafts/draft-request' && request.method === 'DELETE')?.headers).to.deep.equal({ 'If-Match': '"2"' })
    expect(requests.find(request => request.path.endsWith('/draft-request/submit'))?.headers).to.deep.equal({ 'If-Match': '"2"' })
    expect(requests.find(request => request.path === '/api/v1/staff/service-requests')?.body).to.deep.equal({ constituent: { constituent_id: 'constituent-fixture-001' }, request: portalInput })
  })

  it('persists mock drafts, increments versions, and deduplicates logical portal submissions', async () => {
    const provider = new MockC311Provider({ role: 'constituent' })
    const input = {
      summary: 'Pothole near library',
      description: 'The road surface is damaged near the library entrance.',
      service_type: 'POTHOLE' as const,
      requester: { display_name: 'Alex Example', email: 'alex@example.test' },
      location: { address: '100 Example Street', latitude: 42.9001, longitude: -88.8801 },
    }
    const first = await provider.submitPortalRequest(input, { idempotencyKey: 'fe03-submit-1' })
    const replay = await provider.submitPortalRequest(input, { idempotencyKey: 'fe03-submit-1' })
    expect(replay).to.deep.equal(first)
    expect(provider.getWriteCount('portal_service_request_submit')).to.equal(1)
    await expectError(() => provider.submitPortalRequest({ ...input, summary: 'Different summary' }, { idempotencyKey: 'fe03-submit-1' }), 'IDEMPOTENCY_CONFLICT')

    const created = await provider.createDraft(input)
    expect(created.status).to.equal('DRAFT')
    const loaded = await provider.getDraft(created.request_id)
    expect(loaded.summary).to.equal(input.summary)
    const updated = await provider.updateDraft(created.request_id, { summary: 'Changed draft' }, { expectedVersion: created.version })
    expect(updated.summary).to.equal('Changed draft')
    expect(updated.version).to.equal(created.version + 1)
    await expectError(() => provider.updateDraft(created.request_id, { summary: 'Stale update' }, { expectedVersion: created.version }), 'VERSION_CONFLICT')
    await provider.deleteDraft(created.request_id, { expectedVersion: updated.version })
    await expectError(() => provider.getDraft(created.request_id), 'NOT_FOUND')
  })

  it('models attachment staging, one-time token use, and expected-version failures', async () => {
    const provider = new MockC311Provider({ role: 'public_visitor' })
    const attachment = await provider.uploadPortalAttachment({ file: 'opaque-fixture-bytes', filename: 'fixture.txt', media_type: 'text/plain' })
    expect(attachment.attachment_token).to.match(/^attachment-token-fixture-/)
    expect(provider.getWriteCount('portal_attachment_upload')).to.equal(1)
    const input: PortalServiceRequestCreate = {
      summary: 'Fixture attachment request',
      description: 'This request validates one-time attachment staging.',
      service_type: 'GENERAL_INQUIRY',
      requester: { display_name: 'Fixture Resident', email: 'resident@example.test' },
      attachment_tokens: [attachment.attachment_token],
    }
    const first = await provider.submitPortalRequest(input, { idempotencyKey: 'attachment-submit-1' })
    expect(first.status).to.equal('SUBMITTED')
    expect(await provider.submitPortalRequest(input, { idempotencyKey: 'attachment-submit-1' })).to.deep.equal(first)
    await expectError(() => provider.submitPortalRequest(input, { idempotencyKey: 'attachment-submit-2' }), 'IDEMPOTENCY_CONFLICT')
    const missingVersion = await expectError(() => new MockC311Provider({ scenario: 'expected-version-required', role: 'constituent' }).updateDraft('draft-fixture-001', { summary: 'x' }), 'EXPECTED_VERSION_REQUIRED')
    expect(missingVersion.status).to.equal(428)
  })

  it('exposes explicit idempotency and retryable attachment scenarios', async () => {
    const conflict = await expectError(() => new MockC311Provider({ scenario: 'idempotency-conflict' }).submitPortalRequest({
      summary: 'Fixture request',
      description: 'This request is long enough for the fixture.',
      service_type: 'GENERAL_INQUIRY',
      requester: { display_name: 'Fixture Resident', email: 'resident@example.test' },
    }, { idempotencyKey: 'fixture-key' }), 'IDEMPOTENCY_CONFLICT')
    expect(conflict.status).to.equal(409)
    const retryable = await expectError(() => new MockC311Provider({ scenario: 'retryable' }).uploadPortalAttachment({ file: 'fixture', filename: 'fixture.txt', media_type: 'text/plain' }), 'TEMPORARILY_UNAVAILABLE')
    expect(retryable.status).to.equal(503)
    expect(retryable.retryable).to.equal(true)
  })

  it('exposes retryable and terminal portal submission failures', async () => {
    const input = {
      service_type: 'GENERAL_INQUIRY' as const,
      summary: 'A valid summary',
      description: 'A valid description for a portal request.',
      requester: { display_name: 'Fixture Resident', email: 'resident@example.test' },
    }
    const retryable = await expectError(() => new MockC311Provider({ scenario: 'retryable' }).submitPortalRequest(input), 'TEMPORARILY_UNAVAILABLE')
    expect(retryable.status).to.equal(503)
    expect(retryable.retryable).to.equal(true)
    const terminal = await expectError(() => new MockC311Provider({ scenario: 'terminal' }).submitPortalRequest(input), 'OPERATION_FAILED')
    expect(terminal.status).to.equal(500)
    expect(terminal.retryable).to.equal(false)
  })
})
