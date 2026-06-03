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
