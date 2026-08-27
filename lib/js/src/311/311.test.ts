import { expect } from 'chai'
import type { C311Scenario } from './enums'
import { C311ApiError } from './errors'
import { cloneFixtureSet, createDefaultFixtureSet } from './fixtures'
import { MockC311Provider } from './mock-provider'
import { C311FetchTransport, C311HttpProvider, type C311TransportRequest } from './provider'
import { C311_TIMEZONE, formatC311DateTime } from './time'
import type { PortalServiceRequestCreate, ServiceRequestCreate } from './types'

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
    expect(page.items[0].primary_requester.constituent_id).to.equal('constituent-fixture-001')
    expect(page.items[0].location?.address.postal_code).to.equal('14201')
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
    expect(await provider.reopenPortalRequest('request-fixture-001', 'fixture reason')).to.deep.equal({
      request_id: 'request-fixture-001',
      status: 'PENDING_APPROVAL',
    })
  })

  it('normalizes every failure scenario to the shared API error', async () => {
    const scenarios: Array<[string, string, boolean]> = [
      ['forbidden', 'FORBIDDEN', false],
      ['not-found', 'NOT_FOUND', false],
      ['validation', 'VALIDATION_ERROR', false],
      ['retryable', 'TEMPORARILY_UNAVAILABLE', true],
      ['terminal', 'OPERATION_FAILED', false],
      ['version-conflict', 'VERSION_CONFLICT', false],
    ]

    for (const [scenario, code, retryable] of scenarios) {
      const error = await expectError(() => new MockC311Provider({ scenario: scenario as C311Scenario }).getSession(), code)
      expect(error.retryable).to.equal(retryable)
      if (scenario === 'retryable') expect(error.retryAfter).to.equal('30')
    }
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
    const provider = new C311HttpProvider(transport)
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
  })

  it('formats the benchmark instant in the fixed timezone', () => {
    expect(C311_TIMEZONE).to.equal('America/New_York')
    expect(formatC311DateTime('2026-01-15T15:00:00.000Z', 'en-US')).to.contain('10:00')
    expect(formatC311DateTime('not-a-date')).to.equal('')
  })

  it('deep clones custom fixture sets', () => {
    const fixtures = createDefaultFixtureSet()
    const clone = cloneFixtureSet(fixtures)
    clone.requests[0].summary = 'changed in test'

    expect(fixtures.requests[0].summary).to.equal('Pothole on Example Street')
  })
})
