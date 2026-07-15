// Authoritative vocabularies for ActionLog filter selects.
// Values mirror server enums; labels are human-friendly.

const RESOURCE_LABELS = {
  'corteza::automation': 'Automation',
  'corteza::automation:ng-automation': 'Automation (NG)',
  'corteza::automation:session': 'Automation session',
  'corteza::automation:trigger': 'Automation trigger',
  'corteza::automation:workflow': 'Workflow',
  'corteza::compose': 'Compose',
  'corteza::compose:attachment': 'Compose attachment',
  'corteza::compose:chart': 'Chart',
  'corteza::compose:module': 'Module',
  'corteza::compose:module-field': 'Module field',
  'corteza::compose:namespace': 'Namespace',
  'corteza::compose:page': 'Page',
  'corteza::compose:page-layout': 'Page layout',
  'corteza::compose:record': 'Record',
  'corteza::compose:record-revision': 'Record revision',
  'corteza::system': 'System',
  'corteza::system:agent': 'Agent',
  'corteza::system:ai-conversation': 'AI conversation',
  'corteza::system:apigw-filter': 'Integration gateway filter',
  'corteza::system:apigw-route': 'Integration gateway route',
  'corteza::system:application': 'Application',
  'corteza::system:attachment': 'Attachment',
  'corteza::system:auth-client': 'Auth client',
  'corteza::system:auth-confirmed-client': 'Auth confirmed client',
  'corteza::system:auth-oa2token': 'OAuth2 token',
  'corteza::system:auth-session': 'Auth session',
  'corteza::system:chatbot': 'Chatbot',
  'corteza::system:configured-connection': 'Configured connection',
  'corteza::system:connection': 'Connection',
  'corteza::system:credential': 'Credential',
  'corteza::system:dal-connection': 'Data source',
  'corteza::system:dal-schema-alteration': 'Schema alteration',
  'corteza::system:dal-sensitivity-level': 'Sensitivity level',
  'corteza::system:data-privacy-request': 'Data privacy request',
  'corteza::system:data-privacy-request-comment': 'Privacy request comment',
  'corteza::system:knowledge-base': 'Knowledge base',
  'corteza::system:llm-provider': 'LLM provider',
  'corteza::system:notification': 'Notification',
  'corteza::system:queue': 'Queue',
  'corteza::system:queue-message': 'Queue message',
  'corteza::system:reminder': 'Reminder',
  'corteza::system:report': 'Report',
  'corteza::system:resource-translation': 'Resource translation',
  'corteza::system:role': 'Role',
  'corteza::system:role-member': 'Role member',
  'corteza::system:settings': 'Settings',
  'corteza::system:template': 'Template',
  'corteza::system:user': 'User',
  'corteza::system:user-group': 'User group',
}

export const RESOURCE_TYPES = Object.entries(RESOURCE_LABELS)
  .map(([value, label]) => ({ value, label }))
  .sort((a, b) => a.label.localeCompare(b.label))

function lookupLabel(type) {
  if (!type) return null
  return (
    RESOURCE_LABELS[type] ||
    RESOURCE_LABELS[`corteza::${type}`] ||
    RESOURCE_LABELS[type.replace(/^corteza::/, '')] ||
    null
  )
}

export function resourceLabel(value) {
  if (!value) return value
  const direct = lookupLabel(value)
  if (direct) return direct
  // action log resources are often `<type>/<id>` or `<type>/*`
  const slash = value.indexOf('/')
  if (slash !== -1) {
    const type = value.slice(0, slash)
    const rest = value.slice(slash + 1)
    const label = lookupLabel(type)
    if (label) {
      return rest && rest !== '*' ? `${label} (${rest})` : label
    }
  }
  return value
}

const ACTION_LABELS = {
  // CRUD
  create: 'Create',
  read: 'Read',
  update: 'Update',
  patch: 'Patch',
  delete: 'Delete',
  undelete: 'Undelete',
  archive: 'Archive',
  unarchive: 'Unarchive',
  enable: 'Enable',
  suspend: 'Suspend',
  unsuspend: 'Unsuspend',

  // Auth
  authenticate: 'Authenticate',
  login: 'Login',
  logout: 'Logout',
  register: 'Register',
  internalSignup: 'Internal signup',
  externalSignup: 'External signup',
  confirmEmail: 'Confirm email',
  sendEmailConfirmationToken: 'Send email confirmation',
  changePassword: 'Change password',
  setPassword: 'Set password',
  removePassword: 'Remove password',
  sendPasswordResetToken: 'Send password reset',
  exchangePasswordResetToken: 'Exchange password reset token',
  generatePasswordCreateToken: 'Generate password create token',
  emailOtpVerify: 'Verify email OTP',
  totpConfigure: 'Configure TOTP',
  totpRemove: 'Remove TOTP',
  totpValidate: 'Validate TOTP',
  issueToken: 'Issue token',
  validateToken: 'Validate token',
  exposeSecret: 'Expose secret',
  regenerateSecret: 'Regenerate secret',
  deleteAuthTokens: 'Delete auth tokens',
  accessTokensRemoved: 'Access tokens removed',
  deleteAuthSessions: 'Delete auth sessions',
  sessionsRevoke: 'Revoke sessions',
  createCredentials: 'Create credentials',
  updateCredentials: 'Update credentials',

  // Roles & members
  memberAdd: 'Add member',
  memberRemove: 'Remove member',
  members: 'List members',
  grant: 'Grant',
  impersonate: 'Impersonate',

  // Avatars / invites
  generateAvatar: 'Generate avatar',
  uploadAvatar: 'Upload avatar',
  deleteAvatar: 'Delete avatar',
  sendInviteEMail: 'Send invite email',

  // Search / discovery
  search: 'Search',
  searchSensitive: 'Search sensitive',
  searchRevisions: 'Search revisions',
  lookup: 'Lookup',

  // Notifications / dismissal
  approve: 'Approve',
  dismiss: 'Dismiss',
  undismiss: 'Undismiss',
  snooze: 'Snooze',
  markAsRead: 'Mark as read',
  markAsUnread: 'Mark as unread',
  markAllAsRead: 'Mark all as read',
  markAllAsUnread: 'Mark all as unread',

  // Workflow / iterators
  execute: 'Execute',
  run: 'Run',
  apply: 'Apply',
  preprocess: 'Preprocess',
  iteratorInvoked: 'Iterator invoked',
  iteratorIteration: 'Iterator iteration',
  iteratorClone: 'Iterator clone',
  iteratorUpdate: 'Iterator update',
  iteratorDelete: 'Iterator delete',
  iteratorUndelete: 'Iterator undelete',

  // Misc data ops
  clone: 'Clone',
  export: 'Export',
  import: 'Import',
  importInit: 'Import init',
  importRun: 'Import run',
  bulk: 'Bulk',
  organize: 'Organize',
  reorder: 'Reorder',
  merge: 'Merge',
  render: 'Render',
  report: 'Report',
  request: 'Request',
  send: 'Send',
  sign: 'Sign',
  serve: 'Serve',
  attachmentDownload: 'Download attachment',
  flagManage: 'Manage flag',
  flagManageGlobal: 'Manage global flag',
  autoPromote: 'Auto-promote',

  // Backwards-compat dotted names that may appear in older logs
  'password.change': 'Change password',
  'password.reset': 'Reset password',
}

export const COMMON_ACTIONS = Object.entries(ACTION_LABELS)
  .map(([value, label]) => ({ value, label }))
  .sort((a, b) => a.label.localeCompare(b.label))

export function actionLabel(value) {
  return ACTION_LABELS[value] || value
}

const ORIGIN_LABELS = {
  'api/rest': 'REST API',
  'api/grpc': 'gRPC API',
  auth: 'Authentication',
  automation: 'Automation',
  'app/init': 'App init',
  'app/upgrade': 'App upgrade',
  'app/activate': 'App activate',
  'app/provision': 'App provision',
  'app/run': 'App run',
}

export const ORIGINS = Object.entries(ORIGIN_LABELS)
  .map(([value, label]) => ({ value, label }))
  .sort((a, b) => a.label.localeCompare(b.label))

export function originLabel(value) {
  return ORIGIN_LABELS[value] || value
}

// Backend severity is a uint8 (0–7) syslog level. Maps each to a PrimeVue Tag
// `severity` + an i18n label key suffix (system.actionlog.list.severity.<label>).
// 0=Emergency 1=Alert 2=Critical 3=Error 4=Warning 5=Notice 6=Info 7=Debug
export const SEVERITY_MAP = {
  0: { severity: 'danger', label: 'emergency' },
  1: { severity: 'danger', label: 'alert' },
  2: { severity: 'danger', label: 'critical' },
  3: { severity: 'danger', label: 'error' },
  4: { severity: 'warn', label: 'warning' },
  5: { severity: 'success', label: 'notice' },
  6: { severity: 'info', label: 'info' },
  7: { severity: 'secondary', label: 'debug' },
}

// ---------------------------------------------------------------------------
// Helpers for the event timeline (project dashboard). Additive — the admin
// ActionLog list does not use these.
// ---------------------------------------------------------------------------

// resourceType strips the trailing `/<id>[/<id>…]` from an action-log resource,
// leaving the bare type: `corteza::compose:record/1/2/3` → `corteza::compose:record`.
export function resourceType(value) {
  if (!value) return ''
  const slash = value.indexOf('/')
  return slash === -1 ? value : value.slice(0, slash)
}

// resourceTypeLabel is resourceLabel without the ID suffix — `Record`, not
// `Record (1/2/3)`. Use where the ID is shown separately (or not at all).
export function resourceTypeLabel(value) {
  const type = resourceType(value)
  return lookupLabel(type) || type
}

// resourceID returns the last path segment of a resource — usually the ID of the
// affected resource itself. Empty when the resource carries no path.
export function resourceID(value) {
  if (!value) return ''
  const slash = value.lastIndexOf('/')
  if (slash === -1) return ''
  const id = value.slice(slash + 1)
  return id === '*' ? '' : id
}

// Past-tense verbs for the timeline sentence ("Ada *created* Agent foo").
// Only actions that read badly under the generic fallback need an entry; the
// fallback lower-cases the action's label, which is right for most of them
// ("grant" → "grant", "execute" → "execute").
const ACTION_VERBS = {
  create: 'created',
  read: 'read',
  update: 'updated',
  patch: 'patched',
  delete: 'deleted',
  undelete: 'restored',
  archive: 'archived',
  unarchive: 'unarchived',
  enable: 'enabled',
  suspend: 'suspended',
  unsuspend: 'unsuspended',
  authenticate: 'authenticated',
  login: 'signed in',
  logout: 'signed out',
  register: 'registered',
  search: 'searched',
  searchSensitive: 'searched (sensitive)',
  searchRevisions: 'searched revisions of',
  lookup: 'viewed',
  members: 'listed members of',
  memberAdd: 'added a member to',
  memberRemove: 'removed a member from',
  grant: 'changed permissions on',
  impersonate: 'impersonated',
  execute: 'executed',
  run: 'ran',
  apply: 'applied',
  clone: 'cloned',
  export: 'exported',
  import: 'imported',
  merge: 'merged',
  reorder: 'reordered',
  organize: 'organized',
  render: 'rendered',
  send: 'sent',
  sign: 'signed',
  approve: 'approved',
  dismiss: 'dismissed',
  undismiss: 'undismissed',
  snooze: 'snoozed',
  markAsRead: 'marked as read',
  markAsUnread: 'marked as unread',
  markAllAsRead: 'marked all as read',
  markAllAsUnread: 'marked all as unread',
  changePassword: 'changed the password of',
  setPassword: 'set the password of',
  removePassword: 'removed the password of',
  uploadAvatar: 'uploaded an avatar for',
  generateAvatar: 'generated an avatar for',
  deleteAvatar: 'deleted the avatar of',
}

export function actionVerb(value) {
  if (!value) return ''
  return ACTION_VERBS[value] || actionLabel(value).toLowerCase()
}
