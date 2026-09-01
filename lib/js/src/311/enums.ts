export const CONTRACT_VERSION = '1.0.0' as const

export const SERVICE_REQUEST_STATUSES = [
  'DRAFT',
  'SUBMITTED',
  'TRIAGED',
  'ASSIGNED',
  'IN_PROGRESS',
  'RESOLVED',
  'CLOSED',
  'REOPENED',
] as const
export type ServiceRequestStatus = typeof SERVICE_REQUEST_STATUSES[number]

export const SOURCE_CHANNELS = [
  'PORTAL_ANONYMOUS',
  'PORTAL_AUTHENTICATED',
  'STAFF_IN_PERSON',
  'API',
] as const
export type SourceChannel = typeof SOURCE_CHANNELS[number]

export const ORIGIN_CLASSES = ['EXTERNAL', 'INTERNAL'] as const
export type OriginClass = typeof ORIGIN_CLASSES[number]

export const RELATIONSHIP_TYPES = [
  'PRIMARY_REQUESTER',
  'AFFECTED_RESIDENT',
  'PROPERTY_OWNER',
  'REPORTER',
  'ORGANISATION_CONTACT',
] as const
export type RelationshipType = typeof RELATIONSHIP_TYPES[number]

export const REMINDER_STATUSES = ['SCHEDULED', 'SNOOZED', 'COMPLETED', 'CANCELLED'] as const
export type ReminderStatus = typeof REMINDER_STATUSES[number]

export const REMINDER_CHANNELS = ['EMAIL', 'IN_APP'] as const
export type ReminderChannel = typeof REMINDER_CHANNELS[number]

export const REMINDER_ACTIONS = ['SNOOZE', 'COMPLETE', 'CANCEL'] as const
export type ReminderAction = typeof REMINDER_ACTIONS[number]

export const LANGUAGES = ['EN', 'ES', 'VI'] as const
export type Language = typeof LANGUAGES[number]

export const IDENTITY_PROVIDERS = ['oidc', 'saml'] as const
export type IdentityProvider = typeof IDENTITY_PROVIDERS[number]

export const PUBLIC_CONTENT_KEYS = ['HOME', 'SERVICE_CATALOGUE', 'HELP', 'FOOTER', 'TERMS'] as const
export type PublicContentKey = typeof PUBLIC_CONTENT_KEYS[number]

export const HELP_KEYS = [
  'admin.branding.publish',
  'admin.workflow.author',
  'public.request.lookup',
  'public.request.submit',
  'staff.report.create',
  'staff.request.bulk-update',
  'staff.request.reassign',
  'staff.request.triage',
] as const
export type HelpKey = typeof HELP_KEYS[number]

export const CUSTOM_FIELD_TYPES = [
  'TEXT',
  'INTEGER',
  'DECIMAL',
  'DATE',
  'DATETIME',
  'BOOLEAN',
  'SINGLE_CHOICE',
  'MULTI_CHOICE',
] as const
export type CustomFieldType = typeof CUSTOM_FIELD_TYPES[number]

export const CONTACT_CATEGORIES = [
  'RESIDENT',
  'BUSINESS',
  'BUSINESS_OWNER',
  'VETERAN',
  'NEIGHBORHOOD_ASSOCIATION',
  'GOVERNMENT',
  'OTHER',
] as const
export type ContactCategory = typeof CONTACT_CATEGORIES[number]

export const SERVICE_TYPES = ['TREE_MAINTENANCE', 'POTHOLE', 'MISSED_TRASH', 'GENERAL_INQUIRY'] as const
export type ServiceType = typeof SERVICE_TYPES[number]

export const DEPARTMENT_CODES = ['PUBLIC_WORKS', 'STREETS', 'SANITATION', 'GENERAL_SERVICES'] as const
export type DepartmentCode = typeof DEPARTMENT_CODES[number]

export const DISTRICT_CODES = ['NORTH', 'CENTRAL', 'SOUTH'] as const
export type DistrictCode = typeof DISTRICT_CODES[number]

export const CIVICWORKS_STATUSES = ['ASSIGNED', 'IN_PROGRESS', 'PARTIALLY_COMPLETED', 'COMPLETED'] as const
export type CivicWorksStatus = typeof CIVICWORKS_STATUSES[number]

export const ACTOR_ROLES = [
  'service_agent',
  'supervisor',
  'department_manager',
  'platform_administrator',
  'workflow_designer',
] as const
export type ActorRole = typeof ACTOR_ROLES[number]

export const APPLICATION_ROLES = [
  'public_visitor',
  'constituent',
  'service_agent',
  'supervisor',
  'department_manager',
  'platform_administrator',
  'workflow_designer',
  'integration_client',
] as const
export type ApplicationRole = typeof APPLICATION_ROLES[number]

export const AUDIT_ACTOR_TYPES = ['constituent', 'staff', 'integration_client', 'system'] as const
export type AuditActorType = typeof AUDIT_ACTOR_TYPES[number]

export const PHONE_LABELS = ['MOBILE', 'HOME', 'WORK'] as const
export type PhoneLabel = typeof PHONE_LABELS[number]

export const VALIDATION_CODES = [
  'REQUIRED',
  'INVALID_FORMAT',
  'TOO_SHORT',
  'TOO_LONG',
  'OUT_OF_RANGE',
  'INVALID_VALUE',
  'INACTIVE_VALUE',
  'TOO_MANY_ITEMS',
  'DUPLICATE',
  'CONFLICT',
  'LOCATION_REQUIRED',
  'COORDINATES_REQUIRED',
] as const
export type ValidationCode = typeof VALIDATION_CODES[number]

export const ERROR_CODES = [
  'UNAUTHENTICATED',
  'FORBIDDEN',
  'VALIDATION_ERROR',
  'INVALID_STATUS_TRANSITION',
  'IDEMPOTENCY_CONFLICT',
  'VERSION_CONFLICT',
  'INVALID_FILTER',
  'INVALID_PAGE_TOKEN',
  'RATE_LIMITED',
  'ADDRESS_NOT_FOUND',
  'MAP_UNAUTHENTICATED',
  'MAP_TEMPORARILY_UNAVAILABLE',
  'INVALID_RESET_TOKEN',
  'EXPIRED_RESET_TOKEN',
  'INSUFFICIENT_SCOPE',
  'INVALID_CLIENT',
  'INVALID_TOKEN',
  'TEMPORARILY_UNAVAILABLE',
  'NOT_FOUND',
  'EXPECTED_VERSION_REQUIRED',
  'INVALID_SIGNATURE',
  'OPERATION_FAILED',
] as const
export type ErrorCode = typeof ERROR_CODES[number]

export const REQUEST_ACTIONS = [
  'TRIAGE',
  'ASSIGN',
  'START_PROGRESS',
  'RESOLVE',
  'CLOSE',
  'REQUEST_REOPEN',
  'APPROVE_REOPEN',
] as const
export type RequestAction = typeof REQUEST_ACTIONS[number]

export const OPERATION_STATUSES = ['PENDING', 'RUNNING', 'SUCCEEDED', 'FAILED', 'CANCELLED'] as const
export type OperationStatus = typeof OPERATION_STATUSES[number]

export const C311_SCENARIOS = [
  'success',
  'empty',
  'forbidden',
  'not-found',
  'validation',
  'retryable',
  'terminal',
  'version-conflict',
  'public-anonymous',
  'invalid-credentials',
  'registration-validation',
  'expired-reset-token',
  'invalid-reset-token',
  'oidc-failure',
  'saml-failure',
  'account-link-conflict',
  'link-confirmation-required',
  'account-link-success',
  'account-link-cancelled',
  'identity-claims-failure',
  'federated-logout-failure',
  'branding-failure',
  'content-loading-failure',
  'help-loading-failure',
  'empty-catalogue',
  'empty-my-requests',
  'successful-login',
  'successful-registration',
  'successful-reset',
  'account-loading',
] as const
export type C311Scenario = typeof C311_SCENARIOS[number]

// Capabilities and routes are endpoint names in contract.json. Keeping these
// string-branded avoids silently inventing a second, incomplete allow-list.
export type ContractCapability = string
export type ContractRoute = string
export type ContractScope = string
