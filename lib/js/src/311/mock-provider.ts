import { C311ApiError } from './errors'
import { C311_SCENARIOS, type C311Scenario } from './enums'
import { cloneFixtureSet, createDefaultFixtureSet } from './fixtures'
import type {
  AnonymousStatusLookupRequest,
  AnonymousStatusLookupResponse,
  BinaryAttachment,
  C311FixtureSet,
  DraftWrite,
  GeocodeRequest,
  GeocodeResponse,
  LocalSignIn,
  Operation,
  PageResponse,
  PortalAttachment,
  PortalServiceRequestCreate,
  ReportDefinition,
  RequestListQuery,
  RequestQueueItem,
  ReopenRequestResponse,
  ServiceRequest,
  ServiceRequestCreate,
  ServiceRequestResponse,
  Session,
  StaffServiceRequestDetail,
  WorkflowDefinition,
} from './types'
import type { C311Provider, C311RequestOptions, PortalAttachmentUpload, ReportExportOptions, RequestTransition } from './provider'

export interface MockC311ProviderOptions {
  scenario?: C311Scenario
  fixtures?: C311FixtureSet
}

const statusByScenario: Partial<Record<C311Scenario, number>> = {
  forbidden: 403,
  'not-found': 404,
  validation: 422,
  retryable: 503,
  terminal: 500,
  'version-conflict': 409,
}

function copy<T> (value: T): T {
  return JSON.parse(JSON.stringify(value)) as T
}

export class MockC311Provider implements C311Provider {
  private readonly fixtures: C311FixtureSet
  private readonly scenario: C311Scenario

  constructor (options: MockC311ProviderOptions = {}) {
    this.fixtures = cloneFixtureSet(options.fixtures || createDefaultFixtureSet())
    this.scenario = options.scenario || 'success'

    if (!C311_SCENARIOS.includes(this.scenario)) {
      throw new Error(`Unsupported City 311 fixture scenario: ${this.scenario}`)
    }
  }

  private failIfNeeded (): void {
    if (this.scenario === 'success' || this.scenario === 'empty') return

    const payload = this.fixtures.errors[this.scenario]
    if (payload) {
      const headers = this.scenario === 'retryable' ? { 'Retry-After': '30' } : undefined
      throw new C311ApiError(payload, statusByScenario[this.scenario], headers)
    }
  }

  private page<T> (items: T[], query: RequestListQuery = {}): PageResponse<T> {
    return {
      items: copy(items),
      next_page_token: null,
      total_count: items.length,
      applied_filters: { ...query },
      sort: query.sort ? [query.sort] : [],
    }
  }

  private request (requestID: string): ServiceRequest {
    const request = this.fixtures.requests.find(item => item.request_id === requestID)
    if (!request) throw new C311ApiError(this.fixtures.errors['not-found'], 404)
    return request
  }

  async getSession (): Promise<Session> {
    this.failIfNeeded()
    return copy(this.fixtures.session)
  }

  async signIn (_input: LocalSignIn): Promise<Session> {
    this.failIfNeeded()
    return copy(this.fixtures.session)
  }

  async signOut (): Promise<void> {
    this.failIfNeeded()
  }

  async getOperation (operationID: string): Promise<Operation> {
    this.failIfNeeded()
    return {
      operation_id: operationID,
      kind: 'fixture',
      status: 'SUCCEEDED',
      progress: 100,
      result: {},
      error: null,
      created_at: '2026-01-15T15:00:00.000Z',
      updated_at: '2026-01-15T15:00:00.000Z',
      completed_at: '2026-01-15T15:00:00.000Z',
    }
  }

  async uploadPortalAttachment (_input: PortalAttachmentUpload): Promise<PortalAttachment> {
    this.failIfNeeded()
    return {
      attachment_token: 'attachment-token-fixture-001',
      filename: 'fixture.txt',
      media_type: 'text/plain',
      size: 17,
      expires_at: '2026-01-15T16:00:00.000Z',
    }
  }

  async downloadAttachment (attachmentID: string): Promise<BinaryAttachment> {
    this.failIfNeeded()
    const attachment = this.fixtures.attachments[attachmentID]
    if (!attachment) throw new C311ApiError(this.fixtures.errors['not-found'], 404)
    return copy(attachment)
  }

  async createServiceRequest (_input: ServiceRequestCreate, _options: C311RequestOptions = {}): Promise<ServiceRequestResponse> {
    this.failIfNeeded()
    return {
      request_id: 'request-fixture-created',
      request_number: 'SR-2026-00002',
      status: 'SUBMITTED',
      version: 1,
      created_at: '2026-01-15T15:00:00.000Z',
      links: { self: '/api/v1/service-requests/request-fixture-created' },
    }
  }

  async submitPortalRequest (_input: PortalServiceRequestCreate, _options: C311RequestOptions = {}): Promise<ServiceRequestResponse> {
    this.failIfNeeded()
    return {
      request_id: 'request-fixture-submitted',
      request_number: 'SR-2026-00002',
      status: 'SUBMITTED',
      version: 1,
      created_at: '2026-01-15T15:00:00.000Z',
      links: { self: '/api/v1/portal/service-requests/request-fixture-submitted' },
    }
  }

  async createDraft (input: DraftWrite, _options: C311RequestOptions = {}): Promise<ServiceRequest> {
    this.failIfNeeded()
    const base = this.fixtures.requests[0]
    return copy({
      ...base,
      request_id: input.request_id || 'draft-fixture-created',
      status: 'DRAFT',
      summary: input.summary || base.summary,
      description: input.description || base.description,
      service_type: input.service_type || base.service_type,
    })
  }

  async getDraft (requestID: string): Promise<ServiceRequest> {
    this.failIfNeeded()
    const draft = this.fixtures.drafts[requestID]
    if (!draft) throw new C311ApiError(this.fixtures.errors['not-found'], 404)
    const base = this.fixtures.requests[0]
    return copy({
      ...base,
      request_id: requestID,
      status: 'DRAFT',
      summary: 'summary' in draft && typeof draft.summary === 'string' ? draft.summary : base.summary,
      description: 'description' in draft && typeof draft.description === 'string' ? draft.description : base.description,
      service_type: 'service_type' in draft && draft.service_type ? draft.service_type : base.service_type,
    })
  }

  async updateDraft (requestID: string, input: DraftWrite, _options: C311RequestOptions = {}): Promise<ServiceRequest> {
    this.failIfNeeded()
    const current = await this.getDraft(requestID)
    return copy({
      ...current,
      request_id: requestID,
      status: 'DRAFT',
      summary: input.summary || current.summary,
      description: input.description || current.description,
      service_type: input.service_type || current.service_type,
    })
  }

  async deleteDraft (requestID: string): Promise<void> {
    this.failIfNeeded()
    if (!this.fixtures.drafts[requestID]) throw new C311ApiError(this.fixtures.errors['not-found'], 404)
  }

  async submitDraft (requestID: string, _options: C311RequestOptions = {}): Promise<ServiceRequestResponse> {
    this.failIfNeeded()
    await this.getDraft(requestID)
    return {
      request_id: requestID,
      request_number: 'SR-2026-00003',
      status: 'SUBMITTED',
      version: 1,
      created_at: '2026-01-15T15:00:00.000Z',
      links: { self: `/api/v1/portal/service-requests/${requestID}` },
    }
  }

  async listPortalRequests (query: RequestListQuery = {}): Promise<PageResponse<ServiceRequest>> {
    this.failIfNeeded()
    const items = this.scenario === 'empty' ? [] : this.fixtures.requests
    return this.page(items, query)
  }

  async linkAnonymousRequest (input: AnonymousStatusLookupRequest): Promise<ServiceRequest> {
    this.failIfNeeded()
    const item = this.fixtures.requests.find(request => request.request_number === input.request_number && request.primary_requester.emails.includes(input.email))
    if (!item) throw new C311ApiError(this.fixtures.errors['not-found'], 404)
    return copy(item)
  }

  async reopenPortalRequest (requestID: string, _reason: string, _options: C311RequestOptions = {}): Promise<ReopenRequestResponse> {
    this.failIfNeeded()
    this.request(requestID)
    return { request_id: requestID, status: 'PENDING_APPROVAL' }
  }

  async getPublicStatus (input: AnonymousStatusLookupRequest): Promise<AnonymousStatusLookupResponse> {
    this.failIfNeeded()
    const detail = this.fixtures.public_details[input.request_number]
    return { request_detail: detail ? copy(detail) : null }
  }

  async geocode (input: GeocodeRequest): Promise<GeocodeResponse> {
    this.failIfNeeded()
    const result = this.fixtures.geocodes[input.address]
    if (!result) {
      throw new C311ApiError({ error: 'ADDRESS_NOT_FOUND', message: 'The address could not be found.', retryable: false }, 404)
    }
    return copy(result)
  }

  async listStaffRequests (query: RequestListQuery = {}): Promise<PageResponse<RequestQueueItem>> {
    this.failIfNeeded()
    return this.page(this.scenario === 'empty' ? [] : this.fixtures.queue, query)
  }

  async getStaffRequest (requestID: string): Promise<StaffServiceRequestDetail> {
    this.failIfNeeded()
    const detail = this.fixtures.details[requestID]
    if (!detail) throw new C311ApiError(this.fixtures.errors['not-found'], 404)
    return copy(detail)
  }

  async transitionStaffRequest (requestID: string, _input: RequestTransition, _options: C311RequestOptions = {}): Promise<StaffServiceRequestDetail> {
    this.failIfNeeded()
    return this.getStaffRequest(requestID)
  }

  async listReports (): Promise<PageResponse<ReportDefinition>> {
    this.failIfNeeded()
    return this.page(this.scenario === 'empty' ? [] : this.fixtures.reports)
  }

  async getReport (reportID: string): Promise<ReportDefinition> {
    this.failIfNeeded()
    const report = this.fixtures.reports.find(item => item.report_id === reportID)
    if (!report) throw new C311ApiError(this.fixtures.errors['not-found'], 404)
    return copy(report)
  }

  async createReport (input: ReportDefinition): Promise<ReportDefinition> {
    this.failIfNeeded()
    return copy(input)
  }

  async updateReport (reportID: string, input: Partial<ReportDefinition>, _options: C311RequestOptions = {}): Promise<ReportDefinition> {
    this.failIfNeeded()
    const report = await this.getReport(reportID)
    return copy({ ...report, ...input, report_id: reportID })
  }

  async runReport (_input: { definition: ReportDefinition }): Promise<Operation> {
    this.failIfNeeded()
    return {
      operation_id: 'operation-fixture-report',
      kind: 'report_run',
      status: 'PENDING',
      created_at: '2026-01-15T15:00:00.000Z',
      updated_at: '2026-01-15T15:00:00.000Z',
    }
  }

  async exportReport (_reportID: string, _options: ReportExportOptions = {}): Promise<Operation> {
    this.failIfNeeded()
    return {
      operation_id: 'operation-fixture-export',
      kind: 'report_export',
      status: 'PENDING',
      created_at: '2026-01-15T15:00:00.000Z',
      updated_at: '2026-01-15T15:00:00.000Z',
    }
  }

  async listWorkflows (): Promise<PageResponse<WorkflowDefinition>> {
    this.failIfNeeded()
    return this.page(this.scenario === 'empty' ? [] : this.fixtures.workflows)
  }

  async getWorkflow (workflowID: string): Promise<WorkflowDefinition> {
    this.failIfNeeded()
    const workflow = this.fixtures.workflows.find(item => item.workflow_id === workflowID)
    if (!workflow) throw new C311ApiError(this.fixtures.errors['not-found'], 404)
    return copy(workflow)
  }

  async createWorkflow (input: WorkflowDefinition): Promise<WorkflowDefinition> {
    this.failIfNeeded()
    return copy(input)
  }

  async updateWorkflow (workflowID: string, input: Partial<WorkflowDefinition>, _options: C311RequestOptions = {}): Promise<WorkflowDefinition> {
    this.failIfNeeded()
    const workflow = await this.getWorkflow(workflowID)
    return copy({ ...workflow, ...input, workflow_id: workflowID })
  }
}
