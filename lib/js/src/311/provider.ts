import { C311ApiError, type C311FieldError } from './errors'
import type { ErrorCode, HelpKey, IdentityProvider, Language, PublicContentKey } from './enums'
import type {
  AnonymousStatusLookupRequest,
  AnonymousStatusLookupResponse,
  AccountRegistration,
  AccountRegistrationAcknowledgement,
  BinaryAttachment,
  Branding,
  C311EndpointResponse,
  ContentObject,
  DraftWrite,
  GeocodeRequest,
  GeocodeResponse,
  FederatedRedirect,
  HelpContent,
  LanguagePreference,
  LoginIdentifierChange,
  ListQuery,
  LocalSignIn,
  Operation,
  PageResponse,
  PortalAttachment,
  PortalServiceRequestCreate,
  PasswordResetConfirm,
  PasswordResetRequest,
  PasswordResetResponse,
  PasswordChange,
  ProfileUpdate,
  Constituent,
  PublicServiceRequestDetail,
  ReportDefinition,
  RequestListQuery,
  RequestQueueItem,
  RequestSummary,
  ReopenRequestResponse,
  ServiceRequest,
  ServiceRequestCreate,
  ServiceRequestResponse,
  Session,
  FederatedSignInResult,
  StaffServiceRequestDetail,
  StaffServiceRequestCreate,
  WorkflowDefinition,
} from './types'

export interface C311RequestOptions {
  idempotencyKey?: string
  expectedVersion?: number
  headers?: Record<string, string>
}

export interface C311TransportRequest {
  method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  path: string
  query?: Record<string, string | number | boolean | string[] | Record<string, unknown> | undefined>
  body?: unknown
  headers?: Record<string, string>
  acceptedStatuses?: number[]
}

export interface C311Transport {
  request<T> (request: C311TransportRequest): Promise<T>
}

export interface C311FetchTransportOptions {
  baseURL?: string
  fetch?: typeof globalThis.fetch
  headers?: Record<string, string>
}

/**
 * Minimal REST transport used by browser pages. It deliberately owns only
 * HTTP concerns; endpoint paths, DTOs and error semantics stay in the provider.
 */
export class C311FetchTransport implements C311Transport {
  private readonly baseURL: string
  private readonly fetcher: typeof globalThis.fetch
  private readonly defaultHeaders: Record<string, string>

  constructor (options: C311FetchTransportOptions = {}) {
    this.baseURL = (options.baseURL || '').replace(/\/$/, '')
    this.fetcher = options.fetch || globalThis.fetch.bind(globalThis)
    this.defaultHeaders = { Accept: 'application/json', ...(options.headers || {}) }
  }

  async request<T> (request: C311TransportRequest): Promise<T> {
    const url = new URL(`${this.baseURL}${request.path}`, globalThis.location?.origin || 'http://localhost')
    Object.entries(request.query || {}).forEach(([key, value]) => {
      if (value === undefined) return
      if (Array.isArray(value)) value.forEach(item => url.searchParams.append(key, String(item)))
      else if (typeof value === 'object') url.searchParams.set(key, JSON.stringify(value))
      else url.searchParams.set(key, String(value))
    })

    const headers = { ...this.defaultHeaders, ...(request.headers || {}) }
    let body: BodyInit | undefined
    if (request.body !== undefined) {
      const isBlob = typeof Blob !== 'undefined' && request.body instanceof Blob
      const isFormData = typeof FormData !== 'undefined' && request.body instanceof FormData
      if (typeof request.body === 'string' || isBlob || isFormData) {
        body = request.body as BodyInit
      } else {
        headers['Content-Type'] = 'application/json'
        body = JSON.stringify(request.body)
      }
    }

    let response: Response
    try {
      response = await this.fetcher(url.toString(), {
        method: request.method,
        headers,
        body,
        credentials: 'include',
      })
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Network request failed.'
      throw new C311ApiError({ error: 'TEMPORARILY_UNAVAILABLE', message, retryable: true })
    }
    const responseHeaders: Record<string, string> = {}
    response.headers.forEach((value, key) => { responseHeaders[key] = value })

    const contentType = response.headers.get('content-type') || ''
    let responseBody: unknown
    if (response.status !== 204) {
      if (contentType.includes('json')) {
        try {
          if (typeof response.text === 'function') {
            const raw = await response.text()
            if (!raw.trim()) throw new Error('Empty JSON response')
            responseBody = JSON.parse(raw)
          } else if (typeof response.json === 'function') {
            responseBody = await response.json()
          } else {
            throw new Error('JSON response has no parser')
          }
        } catch (_error) {
          throw new C311ApiError({
            error: 'OPERATION_FAILED',
            message: `City 311 returned invalid JSON (${response.status})`,
            retryable: false,
          }, response.status, responseHeaders)
        }
      } else {
        responseBody = await response.text()
      }
    }

    if (!response.ok && request.acceptedStatuses?.includes(response.status)) return responseBody as T

    if (!response.ok) {
      const payload = (responseBody && typeof responseBody === 'object' && 'error' in responseBody && typeof responseBody.error === 'object')
        ? responseBody.error
        : responseBody
      const error = payload && typeof payload === 'object'
        ? payload as { error?: ErrorCode; message?: string; retryable?: boolean; [key: string]: unknown }
        : {}

      throw new C311ApiError({
        error: error.error || 'OPERATION_FAILED',
        message: error.message || `City 311 request failed (${response.status})`,
        retryable: error.retryable === true,
        ...(Array.isArray(error.errors) ? { errors: error.errors as C311FieldError[] } : {}),
        ...(typeof error.current_version === 'number' ? { current_version: error.current_version } : {}),
        ...(typeof error.failing_request_id === 'string' ? { failing_request_id: error.failing_request_id } : {}),
        ...(typeof error.operation_id === 'string' ? { operation_id: error.operation_id } : {}),
      }, response.status, responseHeaders)
    }

    return responseBody as T
  }
}

export interface PortalAttachmentUpload {
  file: Blob | string
  filename: string
  media_type: string
}

export interface ReportExportOptions extends C311RequestOptions {
  format?: 'CSV'
}

export interface RequestTransition {
  to_status: ServiceRequest['status']
  reason?: string
}

export interface C311Provider {
  getSession (): Promise<Session>
  signIn (input: LocalSignIn): Promise<Session>
  signOut (): Promise<void>
  registerAccount (input: AccountRegistration): Promise<AccountRegistrationAcknowledgement>
  requestPasswordReset (input: PasswordResetRequest): Promise<PasswordResetResponse>
  confirmPasswordReset (input: PasswordResetConfirm): Promise<PasswordResetResponse>
  changeLoginIdentifier (input: LoginIdentifierChange): Promise<Session>
  changePassword (input: PasswordChange): Promise<void>
  startFederatedSignIn (provider: IdentityProvider): Promise<FederatedRedirect>
  confirmAccountLink (): Promise<Session>
  completeFederatedSignIn (provider: IdentityProvider, query?: Record<string, string>): Promise<FederatedSignInResult>
  getBranding (): Promise<Branding>
  getPublicContent (contentKey: PublicContentKey): Promise<ContentObject>
  getPublicHelp (helpKey: HelpKey, language?: Language): Promise<HelpContent>
  getProfile (): Promise<Constituent>
  updateProfile (input: ProfileUpdate, options?: C311RequestOptions): Promise<Constituent>
  updateLanguage (language: Language): Promise<LanguagePreference>
  getOperation (operationID: string): Promise<Operation>

  uploadPortalAttachment (input: PortalAttachmentUpload): Promise<PortalAttachment>
  downloadAttachment (attachmentID: string): Promise<BinaryAttachment>
  createServiceRequest (input: ServiceRequestCreate, options?: C311RequestOptions): Promise<ServiceRequestResponse>
  submitPortalRequest (input: PortalServiceRequestCreate, options?: C311RequestOptions): Promise<ServiceRequestResponse>
  createStaffServiceRequest (input: StaffServiceRequestCreate): Promise<StaffServiceRequestDetail>
  createDraft (input: DraftWrite, options?: C311RequestOptions): Promise<ServiceRequest>
  getDraft (requestID: string): Promise<ServiceRequest>
  updateDraft (requestID: string, input: DraftWrite, options?: C311RequestOptions): Promise<ServiceRequest>
  deleteDraft (requestID: string, options?: C311RequestOptions): Promise<void>
  submitDraft (requestID: string, options?: C311RequestOptions): Promise<ServiceRequestResponse>
  listPortalRequests (query?: RequestListQuery): Promise<PageResponse<RequestSummary>>
  linkAnonymousRequest (input: AnonymousStatusLookupRequest): Promise<ServiceRequest>
  reopenPortalRequest (requestID: string, reason: string, options?: C311RequestOptions): Promise<ReopenRequestResponse>
  getPublicStatus (input: AnonymousStatusLookupRequest): Promise<AnonymousStatusLookupResponse>

  geocode (input: GeocodeRequest): Promise<GeocodeResponse>
  listStaffRequests (query?: RequestListQuery): Promise<PageResponse<RequestQueueItem>>
  getStaffRequest (requestID: string): Promise<StaffServiceRequestDetail>
  transitionStaffRequest (requestID: string, input: RequestTransition, options?: C311RequestOptions): Promise<StaffServiceRequestDetail>

  listReports (query?: ListQuery): Promise<PageResponse<ReportDefinition>>
  getReport (reportID: string): Promise<ReportDefinition>
  createReport (input: ReportDefinition): Promise<ReportDefinition>
  updateReport (reportID: string, input: ReportDefinition, options?: C311RequestOptions): Promise<ReportDefinition>
  runReport (input: { definition: ReportDefinition }): Promise<Operation>
  exportReport (reportID: string, options?: ReportExportOptions): Promise<Operation>

  listWorkflows (query?: ListQuery): Promise<PageResponse<WorkflowDefinition>>
  getWorkflow (workflowID: string): Promise<WorkflowDefinition>
  createWorkflow (input: WorkflowDefinition): Promise<WorkflowDefinition>
  updateWorkflow (workflowID: string, input: WorkflowDefinition, options?: C311RequestOptions): Promise<WorkflowDefinition>
}

export class C311HttpProvider implements C311Provider {
  constructor (private readonly transport: C311Transport) {}

  private requestOptions (options: C311RequestOptions = {}, includeExpectedVersion = true): Pick<C311TransportRequest, 'headers'> {
    const headers = { ...(options.headers || {}) }
    if (options.idempotencyKey) headers['Idempotency-Key'] = options.idempotencyKey
    if (includeExpectedVersion && options.expectedVersion !== undefined) headers['If-Match'] = `"${options.expectedVersion}"`
    return Object.keys(headers).length ? { headers } : {}
  }

  private request<T> (request: C311TransportRequest): Promise<T> {
    return this.transport.request<T>(request)
  }

  private listQuery (query: ListQuery = {}): C311TransportRequest['query'] {
    const { page_token, page_size, filters, sort } = query
    return {
      ...(page_token === undefined ? {} : { page_token }),
      ...(page_size === undefined ? {} : { page_size }),
      ...(filters && Object.keys(filters).length ? { filters } : {}),
      ...(sort === undefined ? {} : { sort }),
    }
  }

  private requestListQuery (query: RequestListQuery = {}): C311TransportRequest['query'] {
    const {
      page_token,
      page_size,
      filters = {},
      sort,
      ...filterFields
    } = query
    const normalizedFilters = Object.entries(filterFields).reduce<Record<string, unknown>>((out, [key, value]) => {
      if (value !== undefined) out[key] = value
      return out
    }, { ...filters })

    return this.listQuery({ page_token, page_size, filters: normalizedFilters, sort })
  }

  getSession (): Promise<Session> {
    return this.request({ method: 'GET', path: '/api/v1/session' })
  }

  signIn (input: LocalSignIn): Promise<Session> {
    return this.request({ method: 'POST', path: '/api/v1/session', body: input })
  }

  signOut (): Promise<void> {
    return this.request({ method: 'DELETE', path: '/api/v1/session' })
  }

  registerAccount (input: AccountRegistration): Promise<AccountRegistrationAcknowledgement> {
    return this.request({ method: 'POST', path: '/api/v1/accounts', body: input })
  }

  requestPasswordReset (input: PasswordResetRequest): Promise<PasswordResetResponse> {
    return this.request({ method: 'POST', path: '/api/v1/auth/password-reset/request', body: input })
  }

  confirmPasswordReset (input: PasswordResetConfirm): Promise<PasswordResetResponse> {
    return this.request({ method: 'POST', path: '/api/v1/auth/password-reset/confirm', body: input })
  }

  changeLoginIdentifier (input: LoginIdentifierChange): Promise<Session> {
    return this.request({ method: 'POST', path: '/api/v1/account/login-identifier', body: input })
  }

  changePassword (input: PasswordChange): Promise<void> {
    return this.request({ method: 'POST', path: '/api/v1/account/password', body: input })
  }

  startFederatedSignIn (provider: IdentityProvider): Promise<FederatedRedirect> {
    return this.request({ method: 'GET', path: `/api/v1/auth/${encodeURIComponent(provider)}/start` })
  }

  confirmAccountLink (): Promise<Session> {
    return this.request({ method: 'POST', path: '/api/v1/account/link/confirm', body: {} })
  }

  completeFederatedSignIn (provider: IdentityProvider, query: Record<string, string> = {}): Promise<FederatedSignInResult> {
    return this.request({ method: 'GET', path: `/api/v1/auth/${encodeURIComponent(provider)}/callback`, query })
  }

  getBranding (): Promise<Branding> {
    return this.request({ method: 'GET', path: '/api/v1/public/branding' })
  }

  getPublicContent (contentKey: PublicContentKey): Promise<ContentObject> {
    return this.request({ method: 'GET', path: `/api/v1/public/content/${encodeURIComponent(contentKey)}` })
  }

  getPublicHelp (helpKey: HelpKey, language?: Language): Promise<HelpContent> {
    return this.request({ method: 'GET', path: `/api/v1/public/help/${encodeURIComponent(helpKey)}`, query: language ? { language } : undefined })
  }

  getProfile (): Promise<Constituent> {
    return this.request({ method: 'GET', path: '/api/v1/account/profile' })
  }

  updateProfile (input: ProfileUpdate, options: C311RequestOptions = {}): Promise<Constituent> {
    return this.request({ method: 'PATCH', path: '/api/v1/account/profile', body: input, ...this.requestOptions(options) })
  }

  updateLanguage (language: Language): Promise<LanguagePreference> {
    return this.request({ method: 'PATCH', path: '/api/v1/preferences/language', body: { language } })
  }

  getOperation (operationID: string): Promise<Operation> {
    return this.request({ method: 'GET', path: `/api/v1/operations/${encodeURIComponent(operationID)}` })
  }

  uploadPortalAttachment (input: PortalAttachmentUpload): Promise<PortalAttachment> {
    const body = new FormData()
    body.append('file', input.file)
    body.append('filename', input.filename)
    body.append('media_type', input.media_type)
    return this.request({ method: 'POST', path: '/api/v1/portal/attachments', body })
  }

  downloadAttachment (attachmentID: string): Promise<BinaryAttachment> {
    return this.request({ method: 'GET', path: `/api/v1/attachments/${encodeURIComponent(attachmentID)}` })
  }

  submitPortalRequest (input: PortalServiceRequestCreate, options: C311RequestOptions = {}): Promise<ServiceRequestResponse> {
    return this.request({ method: 'POST', path: '/api/v1/portal/service-requests', body: input, ...this.requestOptions(options, false) })
  }

  createStaffServiceRequest (input: StaffServiceRequestCreate): Promise<StaffServiceRequestDetail> {
    return this.request({ method: 'POST', path: '/api/v1/staff/service-requests', body: input })
  }

  createServiceRequest (input: ServiceRequestCreate, options: C311RequestOptions = {}): Promise<ServiceRequestResponse> {
    return this.request({ method: 'POST', path: '/api/v1/service-requests', body: input, ...this.requestOptions(options, false) })
  }

  createDraft (input: DraftWrite, options: C311RequestOptions = {}): Promise<ServiceRequest> {
    return this.request({ method: 'POST', path: '/api/v1/portal/service-request-drafts', body: input, ...this.requestOptions(options, false) })
  }

  getDraft (requestID: string): Promise<ServiceRequest> {
    return this.request({ method: 'GET', path: `/api/v1/portal/service-request-drafts/${encodeURIComponent(requestID)}` })
  }

  updateDraft (requestID: string, input: DraftWrite, options: C311RequestOptions = {}): Promise<ServiceRequest> {
    return this.request({ method: 'PATCH', path: `/api/v1/portal/service-request-drafts/${encodeURIComponent(requestID)}`, body: input, ...this.requestOptions(options) })
  }

  deleteDraft (requestID: string, options: C311RequestOptions = {}): Promise<void> {
    return this.request({ method: 'DELETE', path: `/api/v1/portal/service-request-drafts/${encodeURIComponent(requestID)}`, ...this.requestOptions(options) })
  }

  submitDraft (requestID: string, options: C311RequestOptions = {}): Promise<ServiceRequestResponse> {
    return this.request({ method: 'POST', path: `/api/v1/portal/service-request-drafts/${encodeURIComponent(requestID)}/submit`, ...this.requestOptions(options) })
  }

  listPortalRequests (query: RequestListQuery = {}): Promise<PageResponse<RequestSummary>> {
    return this.request({ method: 'GET', path: '/api/v1/portal/service-requests', query: this.requestListQuery(query) })
  }

  linkAnonymousRequest (input: AnonymousStatusLookupRequest): Promise<ServiceRequest> {
    return this.request({ method: 'POST', path: '/api/v1/portal/service-requests/link', body: input })
  }

  reopenPortalRequest (requestID: string, reason: string, options: C311RequestOptions = {}): Promise<ReopenRequestResponse> {
    return this.request({ method: 'POST', path: `/api/v1/portal/service-requests/${encodeURIComponent(requestID)}/reopen`, body: { reason }, ...this.requestOptions(options, false) })
  }

  async getPublicStatus (input: AnonymousStatusLookupRequest): Promise<AnonymousStatusLookupResponse> {
    const response = await this.request<AnonymousStatusLookupResponse | Record<string, unknown>>({
      method: 'POST',
      path: '/api/v1/public/service-request-status',
      body: input,
      acceptedStatuses: [404],
    })
    return response && Object.prototype.hasOwnProperty.call(response, 'request_detail')
      ? response as AnonymousStatusLookupResponse
      : { request_detail: null }
  }

  geocode (input: GeocodeRequest): Promise<GeocodeResponse> {
    return this.request({ method: 'POST', path: '/api/v1/geocode', body: input })
  }

  listStaffRequests (query: RequestListQuery = {}): Promise<PageResponse<RequestQueueItem>> {
    return this.request({ method: 'GET', path: '/api/v1/staff/service-requests', query: this.requestListQuery(query) })
  }

  getStaffRequest (requestID: string): Promise<StaffServiceRequestDetail> {
    return this.request({ method: 'GET', path: `/api/v1/staff/service-requests/${encodeURIComponent(requestID)}` })
  }

  transitionStaffRequest (requestID: string, input: RequestTransition, options: C311RequestOptions = {}): Promise<StaffServiceRequestDetail> {
    return this.request({ method: 'POST', path: `/api/v1/staff/service-requests/${encodeURIComponent(requestID)}/transitions`, body: input, ...this.requestOptions(options) })
  }

  listReports (query: ListQuery = {}): Promise<PageResponse<ReportDefinition>> {
    return this.request({ method: 'GET', path: '/api/v1/staff/reports', query: this.listQuery(query) })
  }

  async getReport (reportID: string): Promise<ReportDefinition> {
    const seenPageTokens = new Set<string>()
    let pageToken: string | undefined

    do {
      const page = await this.listReports(pageToken ? { page_token: pageToken } : {})
      const report = page.items.find(item => item.report_id === reportID)
      if (report) return report

      pageToken = page.next_page_token || undefined
      if (pageToken && seenPageTokens.has(pageToken)) {
        throw new C311ApiError({ error: 'INVALID_PAGE_TOKEN', message: 'Report pagination returned a repeated page token.', retryable: false }, 422)
      }
      if (pageToken) seenPageTokens.add(pageToken)
    } while (pageToken)

    throw new C311ApiError({ error: 'NOT_FOUND', message: 'The requested report was not found.', retryable: false }, 404)
  }

  createReport (input: ReportDefinition): Promise<ReportDefinition> {
    return this.request({ method: 'POST', path: '/api/v1/staff/reports', body: input })
  }

  updateReport (reportID: string, input: ReportDefinition, options: C311RequestOptions = {}): Promise<ReportDefinition> {
    return this.request({ method: 'PATCH', path: `/api/v1/staff/reports/${encodeURIComponent(reportID)}`, body: input, ...this.requestOptions(options) })
  }

  runReport (input: { definition: ReportDefinition }): Promise<Operation> {
    return this.request({ method: 'POST', path: '/api/v1/staff/reports/run', body: input })
  }

  exportReport (reportID: string, options: ReportExportOptions = {}): Promise<Operation> {
    return this.request({
      method: 'POST',
      path: `/api/v1/staff/reports/${encodeURIComponent(reportID)}/export`,
      body: { format: options.format || 'CSV' },
      ...this.requestOptions(options, false),
    })
  }

  listWorkflows (query: ListQuery = {}): Promise<PageResponse<WorkflowDefinition>> {
    return this.request({ method: 'GET', path: '/api/v1/admin/workflows', query: this.listQuery(query) })
  }

  getWorkflow (workflowID: string): Promise<WorkflowDefinition> {
    return this.request({ method: 'GET', path: `/api/v1/admin/workflows/${encodeURIComponent(workflowID)}` })
  }

  createWorkflow (input: WorkflowDefinition): Promise<WorkflowDefinition> {
    return this.request({ method: 'POST', path: '/api/v1/admin/workflows', body: input })
  }

  updateWorkflow (workflowID: string, input: WorkflowDefinition, options: C311RequestOptions = {}): Promise<WorkflowDefinition> {
    return this.request({ method: 'PATCH', path: `/api/v1/admin/workflows/${encodeURIComponent(workflowID)}`, body: input, ...this.requestOptions(options) })
  }
}

export type C311TransportResult<T> = C311EndpointResponse<T>
