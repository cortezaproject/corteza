import { APPLICATION_ROLES, CONTRACT_VERSION, type ApplicationRole } from './enums'
import type { C311FixtureSet, C311RoleFixture, Constituent, CurrentActor, PublicHistoryItem, RequestQueueItem, ServiceRequest, Session, StaffServiceRequestDetail } from './types'

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

type RoleDefinition = {
  actor_id: string
  display_name: string
  departments: CurrentActor['department_codes']
  districts: CurrentActor['district_codes']
  capabilities: string[]
  scopes: string[]
  routes: string[]
  deniedScope: string
}

const roleDefinitions: Record<ApplicationRole, RoleDefinition> = {
  public_visitor: {
    actor_id: 'actor-fixture-public',
    display_name: 'Public visitor',
    departments: [],
    districts: [],
    capabilities: [],
    scopes: [],
    routes: ['public_branding_get', 'public_content_get', 'public_help_get', 'anonymous_status_lookup', 'portal_service_request_submit', 'portal_attachment_upload'],
    deniedScope: 'service_requests.write',
  },
  constituent: {
    actor_id: 'actor-fixture-001',
    display_name: 'Alex Example',
    departments: [],
    districts: ['NORTH'],
    capabilities: ['portal_my_requests', 'portal_reopen_request', 'portal_link_anonymous_request', 'attachment_download', 'profile_get', 'profile_update', 'password_change', 'login_identifier_change'],
    scopes: ['service_requests.write'],
    routes: ['session_current', 'portal_my_requests', 'portal_service_request_submit', 'portal_attachment_upload', 'portal_reopen_request', 'portal_link_anonymous_request', 'attachment_download', 'profile_get', 'profile_update', 'password_change', 'login_identifier_change'],
    deniedScope: 'workflow.execute',
  },
  service_agent: {
    actor_id: 'actor-fixture-agent',
    display_name: 'Service agent',
    departments: ['STREETS'],
    districts: ['NORTH'],
    capabilities: ['staff_request_queue', 'staff_request_detail', 'staff_request_transition', 'staff_service_request_create', 'report_catalogue', 'report_export'],
    scopes: ['service_requests.write'],
    routes: ['session_current', 'staff_request_queue', 'staff_request_detail', 'staff_request_transition', 'staff_service_request_create', 'report_catalogue', 'report_export'],
    deniedScope: 'workflow.execute',
  },
  supervisor: {
    actor_id: 'actor-fixture-supervisor',
    display_name: 'Supervisor',
    departments: ['STREETS'],
    districts: ['NORTH', 'CENTRAL'],
    capabilities: ['staff_request_queue', 'staff_request_detail', 'staff_request_transition', 'staff_request_reassign', 'staff_request_bulk', 'report_catalogue', 'report_export', 'audit_list', 'audit_export'],
    scopes: ['service_requests.write'],
    routes: ['session_current', 'staff_request_queue', 'staff_request_detail', 'staff_request_transition', 'staff_request_reassign', 'staff_request_bulk', 'report_catalogue', 'report_export', 'audit_list', 'audit_export'],
    deniedScope: 'workflow.execute',
  },
  department_manager: {
    actor_id: 'actor-fixture-manager',
    display_name: 'Department manager',
    departments: ['STREETS', 'PUBLIC_WORKS'],
    districts: ['NORTH', 'CENTRAL', 'SOUTH'],
    capabilities: ['staff_request_queue', 'staff_request_detail', 'staff_request_transition', 'staff_request_reassign', 'report_catalogue', 'report_export', 'audit_list', 'audit_export'],
    scopes: ['service_requests.write', 'crm.export'],
    routes: ['session_current', 'staff_request_queue', 'staff_request_detail', 'staff_request_transition', 'staff_request_reassign', 'report_catalogue', 'report_export', 'audit_list', 'audit_export'],
    deniedScope: 'workflow.execute',
  },
  platform_administrator: {
    actor_id: 'actor-fixture-admin',
    display_name: 'Platform administrator',
    departments: ['PUBLIC_WORKS', 'STREETS', 'SANITATION', 'GENERAL_SERVICES'],
    districts: ['NORTH', 'CENTRAL', 'SOUTH'],
    capabilities: ['staff_request_queue', 'staff_request_detail', 'staff_request_transition', 'staff_request_reassign', 'staff_request_bulk', 'report_catalogue', 'report_export', 'audit_list', 'audit_export', 'admin_branding_get', 'admin_branding_preview', 'admin_branding_publish', 'admin_branding_rollback', 'admin_branding_update', 'admin_branding_versions', 'admin_content_get', 'admin_content_list', 'admin_content_preview', 'admin_content_publish', 'admin_content_rollback', 'admin_content_update', 'admin_content_versions', 'admin_help_update', 'admin_categories_list', 'admin_categories_create', 'admin_categories_update', 'admin_custom_fields_list', 'admin_custom_fields_create', 'admin_custom_fields_update'],
    scopes: ['service_requests.write', 'crm.export'],
    routes: ['session_current', 'staff_request_queue', 'staff_request_detail', 'staff_request_transition', 'staff_request_reassign', 'staff_request_bulk', 'report_catalogue', 'report_export', 'audit_list', 'audit_export', 'admin_branding_get', 'admin_branding_preview', 'admin_branding_publish', 'admin_branding_rollback', 'admin_branding_update', 'admin_branding_versions', 'admin_content_get', 'admin_content_list', 'admin_content_preview', 'admin_content_publish', 'admin_content_rollback', 'admin_content_update', 'admin_content_versions', 'admin_help_update', 'admin_categories_list', 'admin_categories_create', 'admin_categories_update', 'admin_custom_fields_list', 'admin_custom_fields_create', 'admin_custom_fields_update'],
    deniedScope: 'workflow.execute',
  },
  workflow_designer: {
    actor_id: 'actor-fixture-workflow',
    display_name: 'Workflow designer',
    departments: ['GENERAL_SERVICES'],
    districts: ['NORTH', 'CENTRAL', 'SOUTH'],
    capabilities: ['workflow_list', 'workflow_get', 'workflow_create', 'workflow_update', 'workflow_test', 'workflow_activate', 'workflow_deactivate', 'workflow_execution_list', 'workflow_execution_get'],
    scopes: ['workflow.execute'],
    routes: ['session_current', 'workflow_list', 'workflow_get', 'workflow_create', 'workflow_update', 'workflow_test', 'workflow_activate', 'workflow_deactivate', 'workflow_execution_list', 'workflow_execution_get'],
    deniedScope: 'service_requests.write',
  },
  integration_client: {
    actor_id: 'actor-fixture-integration',
    display_name: 'Integration client',
    departments: [],
    districts: [],
    capabilities: [],
    scopes: ['service_requests.write'],
    routes: ['service_request_create', 'data_export'],
    deniedScope: 'crm.export',
  },
}

function roleSession (role: ApplicationRole, definition: RoleDefinition, expiresAt: string): Session {
  if (role === 'public_visitor') {
    return { authenticated: false, actor: null, preferred_language: 'EN', expires_at: expiresAt }
  }

  return {
    authenticated: true,
    actor: {
      actor_id: definition.actor_id,
      display_name: definition.display_name,
      oidc_actor_type: role === 'constituent' ? 'constituent' : null,
      application_roles: [role],
      department_codes: definition.departments,
      district_codes: definition.districts,
      capabilities: definition.capabilities,
      scopes: definition.scopes,
      available_routes: definition.routes,
    },
    preferred_language: 'EN',
    expires_at: expiresAt,
  }
}

function createRoleFixtures (): Record<ApplicationRole, C311RoleFixture> {
  const fixtures = {} as Record<ApplicationRole, C311RoleFixture>
  APPLICATION_ROLES.forEach(role => {
    const definition = roleDefinitions[role]
    fixtures[role] = {
      session: roleSession(role, definition, '2099-01-15T16:00:00.000Z'),
      expired_session: roleSession(role, definition, '2026-01-15T14:00:00.000Z'),
      denied_route: role === 'public_visitor' || role === 'constituent'
        ? 'staff_request_queue'
        : role === 'platform_administrator' ? 'workflow_list' : 'admin_branding_get',
        denied_capability: role === 'public_visitor' || role === 'constituent' ? 'staff_request_queue' : role === 'platform_administrator' ? 'workflow_list' : 'admin_branding_get',
        denied_scope: definition.deniedScope,
    }
  })
  return fixtures
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
        capabilities: ['portal_my_requests', 'attachment_download'],
        scopes: ['service_requests.write'],
        available_routes: ['session_current', 'portal_my_requests', 'portal_service_request_submit', 'attachment_download'],
      },
      preferred_language: 'EN',
      expires_at: '2099-01-15T16:00:00.000Z',
    },
    role_fixtures: createRoleFixtures(),
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
