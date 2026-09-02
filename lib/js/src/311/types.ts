import type {
  ActorRole,
  ApplicationRole,
  ContactCategory,
  ContractCapability,
  ContractRoute,
  ContractScope,
  DepartmentCode,
  DistrictCode,
  ErrorCode,
  Language,
  PublicContentKey,
  HelpKey,
  OperationStatus,
  OriginClass,
  PhoneLabel,
  RequestAction,
  ReminderChannel,
  ReminderStatus,
  ServiceRequestStatus,
  ServiceType,
  SourceChannel,
  ValidationCode,
} from './enums'
import type { C311ErrorPayload } from './errors'

export type ISODateTime = string

export interface PhoneNumber {
  label: PhoneLabel
  value: string
}

export interface Address {
  line1: string
  line2?: string
  city: string
  region: string
  postal_code: string
  country: string
  primary?: boolean
}

export interface Constituent {
  constituent_id: string
  display_name: string
  login_identifier?: string
  emails: string[]
  phone_numbers: PhoneNumber[]
  addresses: Address[]
  primary_category: ContactCategory
  preferred_language: Language
  email_opt_out: boolean
  custom_fields?: Record<string, unknown>
  version?: number
  updated_at?: ISODateTime
}

export interface CurrentActor {
  actor_id: string
  display_name: string
  oidc_actor_type?: 'constituent' | null
  application_roles: ApplicationRole[]
  department_codes: DepartmentCode[]
  district_codes: DistrictCode[]
  capabilities: ContractCapability[]
  scopes: ContractScope[]
  available_routes: ContractRoute[]
}

export interface Session {
  authenticated: boolean
  actor?: CurrentActor | null
  preferred_language: Language
  expires_at?: ISODateTime | null
}

export interface C311RoleFixture {
  session: Session
  expired_session: Session
  denied_route: ContractRoute
  denied_capability: ContractCapability
  denied_scope: ContractScope
}

export interface LocalSignIn {
  login_identifier: string
  password: string
}

export interface AccountRegistration {
  display_name: string
  email: string
  login_identifier: string
  password: string
  preferred_language: Language
}

export interface AccountRegistrationAcknowledgement {
  accepted: true
}

export interface PasswordResetRequest {
  email: string
}

export interface PasswordResetConfirm {
  token: string
  password: string
}

export interface PasswordResetResponse {
  message: string
}

export interface LoginIdentifierChange {
  current_password: string
  login_identifier: string
}

export interface PasswordChange {
  current_password: string
  new_password: string
}

export interface FederatedRedirect {
  authorization_url: string
}

export interface PendingAccountLink {
  expires_at: ISODateTime
  provider_label?: string
}

export type FederatedSignInResult = {
  outcome: 'authenticated'
  session: Session
} | {
  outcome: 'link_confirmation_required'
  pending_link: PendingAccountLink
}

export interface Branding {
  organisation_name: string
  primary_colour: string
  accent_colour: string
  font_family: string
  published: boolean
  version: number
  updated_at: ISODateTime
  login_header?: string
  public_header?: string
  public_footer?: string
  logo_url?: string | null
  favicon_url?: string | null
  portal_wallpaper_url?: string | null
}

export interface ContentObject {
  content_key: PublicContentKey
  body: string
  state: 'DRAFT' | 'PUBLISHED'
  published: boolean
  version: number
  updated_at: ISODateTime
}

export interface HelpContent {
  help_key: HelpKey
  language: Language
  body: string
  version: number
  updated_at: ISODateTime
}

export interface LanguagePreference {
  language: Language
}

export interface ProfileUpdate {
  display_name?: string
  phone_numbers?: PhoneNumber[]
  addresses?: Address[]
  preferred_language?: Language
  primary_category?: ContactCategory
}

export interface RequesterInput {
  display_name: string
  email: string
  phone?: string
}

export interface LocationInput {
  address: string
  latitude?: number
  longitude?: number
}

export interface ServiceRequestLocation {
  address: Address
  latitude?: number
  longitude?: number
}

export interface AttachmentInput {
  filename: string
  media_type: string
  content_base64: string
}

export interface PortalAttachment {
  attachment_token: string
  filename: string
  media_type: string
  size: number
  expires_at: ISODateTime
}

export interface BinaryAttachment {
  content_type: string
  content_disposition: string
  body: string
}

export interface ServiceRequestCreate {
  summary: string
  description: string
  service_type: ServiceType
  requester: RequesterInput
  location?: LocationInput
  attachments?: AttachmentInput[]
  custom_fields?: Record<string, unknown>
}

export interface PortalServiceRequestCreate {
  summary: string
  description: string
  service_type: ServiceType
  requester: RequesterInput
  location?: LocationInput
  attachment_tokens?: string[]
  custom_fields?: Record<string, unknown>
}

export interface ConstituentReference {
  constituent_id: string
}

export interface ConstituentCreate {
  display_name: string
  email: string
}

export interface StaffServiceRequestCreate {
  constituent: ConstituentReference | ConstituentCreate
  request: PortalServiceRequestCreate
}

export interface ServiceRequest {
  request_id: string
  request_number?: string
  summary: string
  description: string
  service_type: ServiceType
  owning_department: DepartmentCode
  council_district?: DistrictCode
  source_channel: SourceChannel
  origin_class: OriginClass
  status: ServiceRequestStatus
  primary_requester: Constituent
  location?: ServiceRequestLocation
  duplicate_group_id?: string
  custom_fields?: Record<string, unknown>
  version: number
  created_at: ISODateTime
  updated_at: ISODateTime
}

export interface ServiceRequestResponse {
  request_id: string
  request_number: string
  status: ServiceRequestStatus
  version: number
  created_at: ISODateTime
  links: { self: string }
}

export interface RequestSummary {
  request_id: string
  request_number: string
  summary: string
  service_type: ServiceType
  status: ServiceRequestStatus
  owning_department: DepartmentCode
  updated_at: ISODateTime
}

export interface RequestQueueItem extends RequestSummary {
  origin_class: OriginClass
  source_channel: SourceChannel
  version: number
  available_actions: RequestAction[]
  council_district?: DistrictCode
  primary_assignee_id?: string | null
  duplicate_group_id?: string | null
}

export interface PublicHistoryItem {
  action: string
  occurred_at: ISODateTime
  responsible_department: DepartmentCode
}

export interface Reminder {
  reminder_id: string
  request_id: string
  title: string
  due_at: ISODateTime
  timezone: string
  recipient_staff_id: string
  channel: ReminderChannel
  status: ReminderStatus
  completed_at?: ISODateTime | null
  completed_by?: string
}

export interface StaffServiceRequestDetail {
  request: ServiceRequest
  available_actions: RequestAction[]
  primary_assignee_id?: string | null
  collaborator_ids: string[]
  reminders: Reminder[]
  history: PublicHistoryItem[]
  audit: Record<string, unknown>[]
  external_work_order?: Record<string, unknown> | null
}

export interface PublicServiceRequestDetail {
  request_number: string
  summary: string
  service_type: ServiceType
  status: ServiceRequestStatus
  owning_department: DepartmentCode
  updated_at: ISODateTime
  history: PublicHistoryItem[]
}

export interface AnonymousStatusLookupRequest {
  request_number: string
  email: string
}

export interface AnonymousStatusLookupResponse {
  request_detail: PublicServiceRequestDetail | null
}

export interface ReopenRequestResponse {
  request_id: string
  status: 'PENDING_APPROVAL'
}

export interface PageResponse<T> {
  items: T[]
  next_page_token: string | null
  total_count: number
  applied_filters: Record<string, unknown>
  sort: string[]
}

export interface ListQuery {
  page_token?: string
  page_size?: number
  filters?: Record<string, unknown>
  sort?: string
}

export interface RequestListQuery extends ListQuery {
  status?: ServiceRequestStatus
  service_type?: ServiceType
  department?: DepartmentCode
  district?: DistrictCode
  origin_class?: OriginClass
  source_channel?: SourceChannel
  assignee?: string
  collaborator?: string
  category?: ContactCategory
  created_from?: string
  created_to?: string
  duplicate_group?: string
}

export interface DraftWrite extends Partial<PortalServiceRequestCreate> {
  request_id?: string
}

export interface Operation {
  operation_id: string
  kind: string
  status: OperationStatus
  progress?: number
  result?: Record<string, unknown> | null
  error?: C311ErrorPayload | null
  created_at: ISODateTime
  updated_at: ISODateTime
  completed_at?: ISODateTime | null
}

export interface GeocodeRequest {
  address: string
}

export interface GeocodeResponse {
  address: string
  latitude: number
  longitude: number
  precision_digits: 4
  provider: 'BENCHMARK_MAP'
}

export interface ReportDefinition {
  report_id: string
  name: string
  entity: 'service_requests' | 'constituents' | 'follow_up_actions'
  columns: string[]
  filters: Record<string, unknown>
  grouping?: string | null
  sort: string[]
  version: number
  updated_at: ISODateTime
}

export interface WorkflowDefinition {
  workflow_id: string
  name: string
  trigger: 'SERVICE_REQUEST_CREATED' | 'SERVICE_REQUEST_STATUS_CHANGED'
  active: boolean
  conditions: Record<string, unknown>[]
  actions: Record<string, unknown>[]
  version: number
  updated_at: ISODateTime
}

export interface C311FixtureSet {
  fixture_id: 'contract-v1'
  contract_version: '1.0.0'
  session: Session
  role_fixtures: Record<ApplicationRole, C311RoleFixture>
  requests: ServiceRequest[]
  queue: RequestQueueItem[]
  details: Record<string, StaffServiceRequestDetail>
  public_details: Record<string, PublicServiceRequestDetail>
  drafts: Record<string, ServiceRequest | DraftWrite>
  attachments: Record<string, BinaryAttachment>
  reports: ReportDefinition[]
  workflows: WorkflowDefinition[]
  geocodes: Record<string, GeocodeResponse>
  errors: Record<string, C311ErrorPayload>
  branding?: Branding
  public_content?: Record<PublicContentKey, ContentObject>
  public_help?: Record<HelpKey, HelpContent>
}

export interface ValidationFailure {
  field: string
  code: ValidationCode
}

export interface C311EndpointResponse<T> {
  status: number
  headers?: Record<string, string>
  body?: T
}

export type C311ErrorCode = ErrorCode
export type C311ActorRole = ActorRole
