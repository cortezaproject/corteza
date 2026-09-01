import { C311ApiError } from './errors'
import { APPLICATION_ROLES, CONTACT_CATEGORIES, C311_SCENARIOS, LANGUAGES, PHONE_LABELS, type ApplicationRole, type C311Scenario, type HelpKey, type IdentityProvider, type Language, type PublicContentKey } from './enums'
import { cloneFixtureSet, createDefaultFixtureSet } from './fixtures'
import type {
  AccountRegistration,
  AccountRegistrationAcknowledgement,
  AnonymousStatusLookupRequest,
  AnonymousStatusLookupResponse,
  BinaryAttachment,
  Branding,
  C311FixtureSet,
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
  Constituent,
  StaffServiceRequestDetail,
  WorkflowDefinition,
} from './types'
import type { C311Provider, C311RequestOptions, PortalAttachmentUpload, ReportExportOptions, RequestTransition } from './provider'

export interface MockC311ProviderOptions {
  scenario?: C311Scenario
  fixtures?: C311FixtureSet
  role?: ApplicationRole
  sessionVariant?: 'current' | 'expired'
}

const statusByScenario: Partial<Record<C311Scenario, number>> = {
  forbidden: 403,
  'not-found': 404,
  validation: 422,
  retryable: 503,
  'version-conflict': 409,
  'invalid-credentials': 401,
  'registration-validation': 422,
  'expired-reset-token': 422,
  'invalid-reset-token': 422,
  'oidc-failure': 503,
  'saml-failure': 503,
  'branding-failure': 503,
  'content-loading-failure': 503,
  'help-loading-failure': 503,
  'account-loading': 503,
  'identity-claims-failure': 401,
}

function copy<T> (value: T): T {
  return JSON.parse(JSON.stringify(value)) as T
}

function validMockProfileInput (input: ProfileUpdate): boolean {
  const allowed = ['display_name', 'phone_numbers', 'addresses', 'preferred_language', 'primary_category']
  if (Object.keys(input).some(key => !allowed.includes(key))) return false
  if (input.display_name !== undefined && (!input.display_name.trim() || input.display_name.length > 120)) return false
  if (input.preferred_language !== undefined && !LANGUAGES.includes(input.preferred_language)) return false
  if (input.primary_category !== undefined && !CONTACT_CATEGORIES.includes(input.primary_category)) return false
  if (input.phone_numbers !== undefined) {
    if (input.phone_numbers.length > 3 || input.phone_numbers.some(phone => !PHONE_LABELS.includes(phone.label) || typeof phone.value !== 'string')) return false
  }
  if (input.addresses !== undefined) {
    if (input.addresses.length > 5 || input.addresses.some(address => ['line1', 'city', 'region', 'postal_code', 'country'].some(field => !String((address as unknown as Record<string, unknown>)[field] || '').trim()))) return false
    if (input.addresses.length > 0 && input.addresses.filter(address => address.primary).length !== 1) return false
  }
  return true
}

export class MockC311Provider implements C311Provider {
  private readonly fixtures: C311FixtureSet
  private readonly scenario: C311Scenario
  private readonly role?: ApplicationRole
  private readonly sessionVariant: 'current' | 'expired'
  private currentSession: Session
  private resetTokenSerial = 0
  private activeResetToken: string | null = null
  private resetTokenUsed = false
  private pendingAccountLinkProvider: IdentityProvider | null = null
  private pendingAccountLinkExpiresAt: string | null = null
  private pendingAccountLinkConsumed = false
  private profile: Constituent

  constructor (options: MockC311ProviderOptions = {}) {
    this.fixtures = cloneFixtureSet(options.fixtures || createDefaultFixtureSet())
    this.scenario = options.scenario || 'success'
    this.role = options.role
    this.sessionVariant = options.sessionVariant || 'current'

    if (!C311_SCENARIOS.includes(this.scenario)) {
      throw new Error(`Unsupported City 311 fixture scenario: ${this.scenario}`)
    }
    if (this.role && !APPLICATION_ROLES.includes(this.role)) {
      throw new Error(`Unsupported City 311 fixture role: ${this.role}`)
    }
    this.currentSession = this.role
      ? copy(this.sessionVariant === 'expired' ? this.fixtures.role_fixtures[this.role].expired_session : this.fixtures.role_fixtures[this.role].session)
      : copy(this.fixtures.session)
    this.profile = copy(this.fixtures.requests[0].primary_requester)
    this.restorePendingAccountLink()
  }

  private pendingAccountLinkStorageKey (): string {
    return `c311.mock.pending.${this.scenario}.${this.role || 'public_visitor'}`
  }

  private restorePendingAccountLink (): void {
    try {
      const raw = typeof sessionStorage !== 'undefined' ? sessionStorage.getItem(this.pendingAccountLinkStorageKey()) : null
      if (!raw) return
      const pending = JSON.parse(raw) as { provider?: IdentityProvider, expires_at?: string, status?: string }
      if (pending.provider && pending.expires_at && Date.parse(pending.expires_at) > Date.now()) {
        this.pendingAccountLinkProvider = pending.provider
        this.pendingAccountLinkExpiresAt = pending.expires_at
        this.pendingAccountLinkConsumed = pending.status === 'consumed'
      } else {
        sessionStorage.removeItem(this.pendingAccountLinkStorageKey())
      }
    } catch (_error) {
      // Browser storage is optional in non-browser unit tests.
    }
  }

  private persistPendingAccountLink (): void {
    if (!this.pendingAccountLinkProvider || !this.pendingAccountLinkExpiresAt) return
    try {
      if (typeof sessionStorage !== 'undefined') sessionStorage.setItem(this.pendingAccountLinkStorageKey(), JSON.stringify({ provider: this.pendingAccountLinkProvider, expires_at: this.pendingAccountLinkExpiresAt, status: this.pendingAccountLinkConsumed ? 'consumed' : 'pending' }))
    } catch (_error) {
      // Browser storage is optional in non-browser unit tests.
    }
  }

  private clearPendingAccountLink (): void {
    this.pendingAccountLinkProvider = null
    this.pendingAccountLinkExpiresAt = null
    this.pendingAccountLinkConsumed = false
    try {
      if (typeof sessionStorage !== 'undefined') sessionStorage.removeItem(this.pendingAccountLinkStorageKey())
    } catch (_error) {
      // Browser storage is optional in non-browser unit tests.
    }
  }

  private failIfNeeded (supported: readonly C311Scenario[] = []): void {
    if (this.scenario === 'success' || this.scenario === 'empty') return
    if (!supported.includes(this.scenario)) return

    const payload = this.fixtures.errors[this.scenario]
    if (payload) {
      const headers = this.scenario === 'retryable' ? { 'Retry-After': '30' } : undefined
      throw new C311ApiError(payload, statusByScenario[this.scenario], headers)
    }
  }

  private failScenario (scenario: C311Scenario): void {
    const payload = this.fixtures.errors[scenario]
    if (!payload) return
    throw new C311ApiError(payload, statusByScenario[scenario], payload.retryable ? { 'Retry-After': '30' } : undefined)
  }

  private page<T> (items: T[], query: ListQuery = {}): PageResponse<T> {
    const { page_token: _pageToken, page_size: _pageSize, filters = {}, sort, ...filterFields } = query as RequestListQuery
    const appliedFilters = Object.entries(filterFields).reduce<Record<string, unknown>>((out, [key, value]) => {
      if (value !== undefined) out[key] = value
      return out
    }, { ...filters })

    return {
      items: copy(items),
      next_page_token: null,
      total_count: items.length,
      applied_filters: appliedFilters,
      sort: sort ? [sort] : [],
    }
  }

  private request (requestID: string): ServiceRequest {
    const request = this.fixtures.requests.find(item => item.request_id === requestID)
    if (!request) throw new C311ApiError(this.fixtures.errors['not-found'], 404)
    return request
  }

  private requestSummary (request: ServiceRequest): RequestSummary {
    return {
      request_id: request.request_id,
      request_number: request.request_number || '',
      summary: request.summary,
      service_type: request.service_type,
      status: request.status,
      owning_department: request.owning_department,
      updated_at: request.updated_at,
    }
  }

  async getSession (): Promise<Session> {
    this.failIfNeeded()
    return copy(this.currentSession)
  }

  async signIn (_input: LocalSignIn): Promise<Session> {
    if (this.scenario === 'invalid-credentials') this.failScenario('invalid-credentials')
    this.failIfNeeded(['validation'])
    this.currentSession = copy(this.fixtures.role_fixtures.constituent.session)
    return copy(this.currentSession)
  }

  async signOut (): Promise<void> {
    if (this.scenario === 'federated-logout-failure') this.failScenario('federated-logout-failure')
    this.failIfNeeded(['forbidden', 'not-found', 'validation'])
    this.currentSession = { authenticated: false, actor: null, preferred_language: 'EN', expires_at: null }
  }

  async registerAccount (_input: AccountRegistration): Promise<AccountRegistrationAcknowledgement> {
    if (this.scenario === 'registration-validation') this.failScenario('registration-validation')
    return { accepted: true }
  }

  async requestPasswordReset (_input: PasswordResetRequest): Promise<PasswordResetResponse> {
    this.failIfNeeded(['retryable', 'terminal'])
    this.resetTokenSerial += 1
    this.activeResetToken = `reset-token-fixture-${String(this.resetTokenSerial).padStart(3, '0')}`
    this.resetTokenUsed = false
    return { message: 'If the account exists, instructions have been sent.' }
  }

  async confirmPasswordReset (input: PasswordResetConfirm): Promise<PasswordResetResponse> {
    if (this.scenario === 'expired-reset-token') this.failScenario('expired-reset-token')
    if (this.scenario === 'invalid-reset-token') this.failScenario('invalid-reset-token')
    if (!this.activeResetToken && input.token === 'ephemeral-token' && ['success', 'successful-reset'].includes(this.scenario)) {
      this.activeResetToken = input.token
      this.resetTokenUsed = false
    }
    if (input.token !== this.activeResetToken || this.resetTokenUsed) this.failScenario('invalid-reset-token')
    this.resetTokenUsed = true
    return { message: 'Your password has been reset.' }
  }

  async changeLoginIdentifier (_input: LoginIdentifierChange): Promise<Session> {
    if (this.scenario === 'invalid-credentials') this.failScenario('invalid-credentials')
    this.failIfNeeded(['forbidden', 'not-found', 'validation', 'version-conflict'])
    return copy(this.currentSession)
  }

  async changePassword (_input: PasswordChange): Promise<void> {
    if (this.scenario === 'invalid-credentials') this.failScenario('invalid-credentials')
    this.failIfNeeded(['forbidden', 'not-found', 'validation', 'version-conflict'])
  }

  async startFederatedSignIn (provider: IdentityProvider): Promise<FederatedRedirect> {
    if (provider === 'saml') throw new C311ApiError({ error: 'FORBIDDEN', message: 'SAML sign-in is available on the staff surface only.', retryable: false }, 403)
    if (this.scenario === 'oidc-failure' && provider === 'oidc') this.failScenario('oidc-failure')
    return { authorization_url: `https://identity.example.test/${provider}/authorize` }
  }

  async confirmAccountLink (): Promise<Session> {
    if (!this.pendingAccountLinkProvider || !this.pendingAccountLinkExpiresAt || Date.parse(this.pendingAccountLinkExpiresAt) <= Date.now()) {
      this.clearPendingAccountLink()
      throw new C311ApiError({ error: 'VALIDATION_ERROR', message: 'The account-link confirmation is no longer valid.', retryable: false }, 422)
    }
    if (this.pendingAccountLinkConsumed) throw new C311ApiError({ error: 'VERSION_CONFLICT', message: 'The account-link confirmation was already consumed.', retryable: false }, 409)
    this.pendingAccountLinkConsumed = true
    this.persistPendingAccountLink()
    if (this.scenario === 'account-link-conflict') this.failScenario('account-link-conflict')
    if (this.scenario === 'identity-claims-failure') this.failScenario('identity-claims-failure')
    this.currentSession = copy(this.fixtures.role_fixtures.constituent.session)
    return copy(this.currentSession)
  }

  async completeFederatedSignIn (_provider: IdentityProvider, _query: Record<string, string> = {}): Promise<FederatedSignInResult> {
    if (_provider === 'saml') throw new C311ApiError({ error: 'FORBIDDEN', message: 'SAML sign-in is available on the staff surface only.', retryable: false }, 403)
    if (['access_denied', 'cancelled', 'canceled'].includes(String(_query.error || '').toLowerCase())) {
      throw new C311ApiError({ error: 'UNAUTHENTICATED', message: 'Federated sign-in was cancelled.', retryable: false }, 401)
    }
    if (this.scenario === 'oidc-failure') this.failScenario('oidc-failure')
    if (this.scenario === 'saml-failure') this.failScenario('saml-failure')
    if (this.scenario === 'identity-claims-failure') this.failScenario('identity-claims-failure')
    if (['link-confirmation-required', 'account-link-success', 'account-link-cancelled', 'account-link-conflict'].includes(this.scenario)) {
      this.pendingAccountLinkProvider = 'oidc'
      this.pendingAccountLinkExpiresAt = new Date(Date.now() + 10 * 60 * 1000).toISOString()
      this.pendingAccountLinkConsumed = false
      this.persistPendingAccountLink()
      return { outcome: 'link_confirmation_required', pending_link: { expires_at: this.pendingAccountLinkExpiresAt, provider_label: 'OIDC' } }
    }
    this.currentSession = copy(this.fixtures.role_fixtures.constituent.session)
    return { outcome: 'authenticated', session: copy(this.currentSession) }
  }

  async getBranding (): Promise<Branding> {
    if (this.scenario === 'branding-failure') this.failScenario('branding-failure')
    this.failIfNeeded(['terminal'])
    return copy(this.fixtures.branding || createDefaultFixtureSet().branding!)
  }

  async getPublicContent (contentKey: PublicContentKey): Promise<ContentObject> {
    if (this.scenario === 'content-loading-failure') this.failScenario('content-loading-failure')
    this.failIfNeeded(['terminal'])
    const content = this.fixtures.public_content?.[contentKey]
    if (!content) throw new C311ApiError(this.fixtures.errors['not-found'], 404)
    if (this.scenario === 'empty-catalogue' && contentKey === 'SERVICE_CATALOGUE') return { ...copy(content), body: '' }
    return copy(content)
  }

  async getPublicHelp (helpKey: HelpKey, language?: Language): Promise<HelpContent> {
    if (this.scenario === 'help-loading-failure') this.failScenario('help-loading-failure')
    this.failIfNeeded(['terminal'])
    const content = this.fixtures.public_help?.[helpKey]
    if (!content) throw new C311ApiError(this.fixtures.errors['not-found'], 404)
    const result = copy(content)
    if (language && language !== result.language) {
      result.language = language
      result.body = language === 'ES' ? '<p>Describe el problema y envialo a la ciudad.</p>' : language === 'VI' ? '<p>Mo ta van de va gui den thanh pho.</p>' : '<p>Describe the issue and submit it to the city.</p>'
    }
    return result
  }

  async getProfile (): Promise<Constituent> {
    if (this.scenario === 'account-loading') this.failScenario('account-loading')
    return copy(this.profile)
  }

  async updateProfile (input: ProfileUpdate, options: C311RequestOptions = {}): Promise<Constituent> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation', 'version-conflict'])
    if (options.expectedVersion === undefined) throw new C311ApiError({ error: 'EXPECTED_VERSION_REQUIRED', message: 'If-Match is required for this update.', retryable: false }, 428)
    if (options.expectedVersion !== undefined && options.expectedVersion !== this.profile.version) this.failScenario('version-conflict')
    if (!validMockProfileInput(input)) this.failScenario('validation')
    this.profile = { ...this.profile, ...input, version: (this.profile.version || 0) + 1, updated_at: new Date().toISOString() }
    return copy(this.profile)
  }

  async updateLanguage (language: Language): Promise<LanguagePreference> {
    this.currentSession = { ...this.currentSession, preferred_language: language }
    return { language }
  }

  async getOperation (operationID: string): Promise<Operation> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation'])
    if (this.scenario === 'terminal') {
      return {
        operation_id: operationID,
        kind: 'fixture',
        status: 'FAILED',
        progress: 100,
        result: null,
        error: copy(this.fixtures.errors.terminal),
        created_at: '2026-01-15T15:00:00.000Z',
        updated_at: '2026-01-15T15:00:00.000Z',
        completed_at: '2026-01-15T15:00:00.000Z',
      }
    }
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
    this.failIfNeeded(['validation'])
    return {
      attachment_token: 'attachment-token-fixture-001',
      filename: 'fixture.txt',
      media_type: 'text/plain',
      size: 17,
      expires_at: '2026-01-15T16:00:00.000Z',
    }
  }

  async downloadAttachment (attachmentID: string): Promise<BinaryAttachment> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation'])
    const attachment = this.fixtures.attachments[attachmentID]
    if (!attachment) throw new C311ApiError(this.fixtures.errors['not-found'], 404)
    return copy(attachment)
  }

  async createServiceRequest (_input: ServiceRequestCreate, _options: C311RequestOptions = {}): Promise<ServiceRequestResponse> {
    this.failIfNeeded(['forbidden', 'validation', 'version-conflict'])
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
    this.failIfNeeded(['validation'])
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
    this.failIfNeeded(['forbidden', 'not-found', 'validation'])
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
    this.failIfNeeded(['forbidden', 'not-found', 'validation'])
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
    this.failIfNeeded(['forbidden', 'not-found', 'validation', 'version-conflict'])
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

  async deleteDraft (requestID: string, _options: C311RequestOptions = {}): Promise<void> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation', 'version-conflict'])
    if (!this.fixtures.drafts[requestID]) throw new C311ApiError(this.fixtures.errors['not-found'], 404)
  }

  async submitDraft (requestID: string, _options: C311RequestOptions = {}): Promise<ServiceRequestResponse> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation', 'version-conflict'])
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

  async listPortalRequests (query: RequestListQuery = {}): Promise<PageResponse<RequestSummary>> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation', 'retryable', 'version-conflict', 'terminal'])
    const items = this.scenario === 'empty' || this.scenario === 'empty-my-requests' ? [] : this.fixtures.requests
    return this.page(items.map(request => this.requestSummary(request)), query)
  }

  async linkAnonymousRequest (input: AnonymousStatusLookupRequest): Promise<ServiceRequest> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation'])
    const item = this.fixtures.requests.find(request => request.request_number === input.request_number && request.primary_requester.emails.includes(input.email))
    if (!item) throw new C311ApiError(this.fixtures.errors['not-found'], 404)
    return copy(item)
  }

  async reopenPortalRequest (requestID: string, _reason: string, _options: C311RequestOptions = {}): Promise<ReopenRequestResponse> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation'])
    this.request(requestID)
    return { request_id: requestID, status: 'PENDING_APPROVAL' }
  }

  async getPublicStatus (input: AnonymousStatusLookupRequest): Promise<AnonymousStatusLookupResponse> {
    if (this.scenario === 'not-found') return { request_detail: null }
    this.failIfNeeded()
    const normalizedEmail = typeof input.email === 'string' ? input.email.trim().toLowerCase() : ''
    const request = this.fixtures.requests.find(item => item.request_number === input.request_number && item.primary_requester.emails.some(email => email.toLowerCase() === normalizedEmail))
    const detail = request ? this.fixtures.public_details[input.request_number] : undefined
    return { request_detail: detail ? copy(detail) : null }
  }

  async geocode (input: GeocodeRequest): Promise<GeocodeResponse> {
    if (this.scenario === 'retryable') {
      throw new C311ApiError({ error: 'MAP_TEMPORARILY_UNAVAILABLE', message: 'The mapping service is temporarily unavailable.', retryable: true }, 503, { 'Retry-After': '30' })
    }
    if (this.scenario === 'not-found') {
      throw new C311ApiError({ error: 'ADDRESS_NOT_FOUND', message: 'The address could not be found.', retryable: false }, 404)
    }
    this.failIfNeeded(['validation'])
    const result = this.fixtures.geocodes[input.address]
    if (!result) {
      throw new C311ApiError({ error: 'ADDRESS_NOT_FOUND', message: 'The address could not be found.', retryable: false }, 404)
    }
    return copy(result)
  }

  async listStaffRequests (query: RequestListQuery = {}): Promise<PageResponse<RequestQueueItem>> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation', 'retryable', 'version-conflict', 'terminal'])
    return this.page(this.scenario === 'empty' ? [] : this.fixtures.queue, query)
  }

  async getStaffRequest (requestID: string): Promise<StaffServiceRequestDetail> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation'])
    const detail = this.fixtures.details[requestID]
    if (!detail) throw new C311ApiError(this.fixtures.errors['not-found'], 404)
    return copy(detail)
  }

  async transitionStaffRequest (requestID: string, _input: RequestTransition, _options: C311RequestOptions = {}): Promise<StaffServiceRequestDetail> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation', 'version-conflict'])
    return this.getStaffRequest(requestID)
  }

  async listReports (query: ListQuery = {}): Promise<PageResponse<ReportDefinition>> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation'])
    return this.page(this.scenario === 'empty' ? [] : this.fixtures.reports, query)
  }

  async getReport (reportID: string): Promise<ReportDefinition> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation'])
    const report = this.fixtures.reports.find(item => item.report_id === reportID)
    if (!report) throw new C311ApiError(this.fixtures.errors['not-found'], 404)
    return copy(report)
  }

  async createReport (input: ReportDefinition): Promise<ReportDefinition> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation'])
    return copy(input)
  }

  async updateReport (reportID: string, input: ReportDefinition, _options: C311RequestOptions = {}): Promise<ReportDefinition> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation', 'version-conflict'])
    const report = await this.getReport(reportID)
    return copy({ ...report, ...input, report_id: reportID })
  }

  async runReport (_input: { definition: ReportDefinition }): Promise<Operation> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation'])
    return {
      operation_id: 'operation-fixture-report',
      kind: 'report_run',
      status: 'PENDING',
      created_at: '2026-01-15T15:00:00.000Z',
      updated_at: '2026-01-15T15:00:00.000Z',
    }
  }

  async exportReport (_reportID: string, _options: ReportExportOptions = {}): Promise<Operation> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation'])
    return {
      operation_id: 'operation-fixture-export',
      kind: 'report_export',
      status: 'PENDING',
      created_at: '2026-01-15T15:00:00.000Z',
      updated_at: '2026-01-15T15:00:00.000Z',
    }
  }

  async listWorkflows (query: ListQuery = {}): Promise<PageResponse<WorkflowDefinition>> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation'])
    return this.page(this.scenario === 'empty' ? [] : this.fixtures.workflows, query)
  }

  async getWorkflow (workflowID: string): Promise<WorkflowDefinition> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation'])
    const workflow = this.fixtures.workflows.find(item => item.workflow_id === workflowID)
    if (!workflow) throw new C311ApiError(this.fixtures.errors['not-found'], 404)
    return copy(workflow)
  }

  async createWorkflow (input: WorkflowDefinition): Promise<WorkflowDefinition> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation'])
    return copy(input)
  }

  async updateWorkflow (workflowID: string, input: WorkflowDefinition, _options: C311RequestOptions = {}): Promise<WorkflowDefinition> {
    this.failIfNeeded(['forbidden', 'not-found', 'validation', 'version-conflict'])
    const workflow = await this.getWorkflow(workflowID)
    return copy({ ...workflow, ...input, workflow_id: workflowID })
  }
}
