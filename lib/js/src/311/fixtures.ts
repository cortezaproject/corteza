import { CONTRACT_VERSION } from './enums'
import type { C311FixtureSet, Constituent, PublicHistoryItem, RequestQueueItem, ServiceRequest, StaffServiceRequestDetail } from './types'

export const BENCHMARK_NOW = '2026-01-15T15:00:00.000Z'
export const BENCHMARK_TIMEZONE = 'America/New_York'
export const C311_FIXTURE_ID = 'contract-v1'

const constituent: Constituent = {
  constituent_id: 'constituent-fixture-001',
  display_name: 'Alex Example',
  login_identifier: 'alex.example',
  emails: ['alex@example.test'],
  phone_numbers: [{ label: 'MOBILE', value: '+1 555 0100' }],
  addresses: [{
    line1: '100 Example Street',
    city: 'Buffalo',
    region: 'NY',
    postal_code: '14201',
    country: 'US',
    primary: true,
  }],
  primary_category: 'RESIDENT',
  preferred_language: 'EN',
  email_opt_out: false,
}

const request: ServiceRequest = {
  request_id: 'request-fixture-001',
  request_number: 'SR-2026-00001',
  summary: 'Pothole on Example Street',
  description: 'A pothole is affecting traffic near the intersection.',
  service_type: 'POTHOLE',
  owning_department: 'STREETS',
  council_district: 'NORTH',
  source_channel: 'PORTAL_AUTHENTICATED',
  origin_class: 'EXTERNAL',
  status: 'SUBMITTED',
  primary_requester: constituent,
  location: {
    address: {
      line1: '100 Example Street',
      city: 'Buffalo',
      region: 'NY',
      postal_code: '14201',
      country: 'US',
      primary: true,
    },
    latitude: 42.9001,
    longitude: -78.8801,
  },
  version: 1,
  created_at: BENCHMARK_NOW,
  updated_at: BENCHMARK_NOW,
}

const history: PublicHistoryItem[] = [
  {
    action: 'SUBMITTED',
    occurred_at: BENCHMARK_NOW,
    responsible_department: 'STREETS',
  },
]

const queueItem: RequestQueueItem = {
  request_id: request.request_id,
  request_number: request.request_number || '',
  summary: request.summary,
  service_type: request.service_type,
  status: request.status,
  owning_department: request.owning_department,
  updated_at: request.updated_at,
  origin_class: request.origin_class,
  source_channel: request.source_channel,
  version: request.version,
  available_actions: ['TRIAGE', 'ASSIGN'],
  council_district: request.council_district,
  primary_assignee_id: null,
  duplicate_group_id: null,
}

const detail: StaffServiceRequestDetail = {
  request,
  available_actions: ['TRIAGE', 'ASSIGN'],
  primary_assignee_id: null,
  collaborator_ids: [],
  reminders: [],
  history,
  audit: [],
  external_work_order: null,
}

export function createDefaultFixtureSet (): C311FixtureSet {
  return {
    fixture_id: C311_FIXTURE_ID,
    contract_version: CONTRACT_VERSION,
    session: {
      authenticated: true,
      actor: {
        actor_id: 'actor-fixture-001',
        display_name: 'Alex Example',
        oidc_actor_type: null,
        application_roles: ['constituent'],
        department_codes: [],
        district_codes: ['NORTH'],
        capabilities: ['portal_my_requests', 'portal_service_request_submit', 'attachment_download'],
        scopes: ['service_requests.write'],
        available_routes: ['session_current', 'portal_my_requests', 'portal_service_request_submit', 'attachment_download'],
      },
      preferred_language: 'EN',
      expires_at: '2026-01-15T16:00:00.000Z',
    },
    requests: [request],
    queue: [queueItem],
    details: { [request.request_id]: detail },
    public_details: {
      [request.request_number || '']: {
        request_number: request.request_number || '',
        summary: request.summary,
        service_type: request.service_type,
        status: request.status,
        owning_department: request.owning_department,
        updated_at: request.updated_at,
        history,
      },
    },
    drafts: {
      'draft-fixture-001': {
        request_id: 'draft-fixture-001',
        summary: 'Saved draft',
        description: 'A draft used for read-only UI checks.',
        service_type: 'GENERAL_INQUIRY',
        requester: {
          display_name: 'Jamie Example',
          email: 'jamie@example.test',
        },
      },
    },
    attachments: {
      'attachment-fixture-001': {
        content_type: 'text/plain',
        content_disposition: 'inline; filename="fixture.txt"',
        body: 'fixture attachment',
      },
    },
    reports: [
      {
        report_id: 'report-fixture-001',
        name: 'Request volume',
        entity: 'service_requests',
        columns: ['request_number', 'status'],
        filters: {},
        grouping: null,
        sort: ['-created_at'],
        version: 1,
        updated_at: BENCHMARK_NOW,
      },
    ],
    workflows: [
      {
        workflow_id: 'workflow-fixture-001',
        name: 'Notify department',
        trigger: 'SERVICE_REQUEST_CREATED',
        active: true,
        conditions: [],
        actions: [{ type: 'notify' }],
        version: 1,
        updated_at: BENCHMARK_NOW,
      },
    ],
    geocodes: {
      '100 Example Street, Buffalo, NY 14201': {
        address: '100 Example Street, Buffalo, NY 14201',
        latitude: 42.9001,
        longitude: -78.8801,
        precision_digits: 4,
        provider: 'BENCHMARK_MAP',
      },
    },
    errors: {
      forbidden: {
        error: 'FORBIDDEN',
        message: 'You are not authorized for this operation.',
        retryable: false,
      },
      'not-found': {
        error: 'NOT_FOUND',
        message: 'The requested resource was not found.',
        retryable: false,
      },
      validation: {
        error: 'VALIDATION_ERROR',
        message: 'One or more fields are invalid.',
        retryable: false,
        errors: [{ field: '/summary', code: 'REQUIRED' }],
      },
      retryable: {
        error: 'TEMPORARILY_UNAVAILABLE',
        message: 'The service is temporarily unavailable.',
        retryable: true,
      },
      terminal: {
        error: 'OPERATION_FAILED',
        message: 'The operation failed permanently.',
        retryable: false,
      },
      'version-conflict': {
        error: 'VERSION_CONFLICT',
        message: 'The record changed before your update.',
        retryable: false,
        current_version: 2,
      },
    },
  }
}

export function cloneFixtureSet (fixtures: C311FixtureSet = createDefaultFixtureSet()): C311FixtureSet {
  return JSON.parse(JSON.stringify(fixtures)) as C311FixtureSet
}
