// Project member roles. Each is a fixed preset: assigning a role to a member
// sets these capability flags (they are not edited per member).
//
// Project Members are human beings with roles to either (1) build part or all of
// the project and its AI systems, or (2) manage project governance and make
// approvals at governance gates.
//
// `resources` summarises what the role can reach; `description` is the role's
// responsibility statement — together they form the accountability framework
// (EU AI Act Article 17(m)).
//
// Tab visibility (see Wizard): a member who can request OR grant approval works
// in the Governance view; everyone else gets the Build view only.

// `labelKey`/`resourcesKey`/`descriptionKey` are i18n keys; components resolve
// them with $t for display. Capability flags drive behaviour and stay literal.
export const ROLE_PRESETS = [
  {
    id: 'developer',
    labelKey: 'project.roles.developer.label',
    read: true,
    write: true,
    requestApproval: true,
    grantApproval: false,
    resourcesKey: 'project.roles.developer.resources',
    descriptionKey: 'project.roles.developer.description',
  },
  {
    id: 'governance-owner',
    labelKey: 'project.roles.governance-owner.label',
    read: true,
    write: true,
    requestApproval: false,
    grantApproval: true,
    resourcesKey: 'project.roles.governance-owner.resources',
    descriptionKey: 'project.roles.governance-owner.description',
  },
  {
    id: 'security-owner',
    labelKey: 'project.roles.security-owner.label',
    read: true,
    write: true,
    requestApproval: false,
    grantApproval: true,
    resourcesKey: 'project.roles.security-owner.resources',
    descriptionKey: 'project.roles.security-owner.description',
  },
  {
    id: 'junior-developer',
    labelKey: 'project.roles.junior-developer.label',
    read: true,
    write: true,
    requestApproval: false,
    grantApproval: false,
    resourcesKey: 'project.roles.junior-developer.resources',
    descriptionKey: 'project.roles.junior-developer.description',
  },
  {
    id: 'executive-authority',
    labelKey: 'project.roles.executive-authority.label',
    read: true,
    write: false,
    requestApproval: false,
    grantApproval: true,
    resourcesKey: 'project.roles.executive-authority.resources',
    descriptionKey: 'project.roles.executive-authority.description',
  },
  {
    id: 'infrastructure-administrator',
    labelKey: 'project.roles.infrastructure-administrator.label',
    read: false,
    write: false,
    requestApproval: false,
    grantApproval: false,
    resourcesKey: 'project.roles.infrastructure-administrator.resources',
    descriptionKey: 'project.roles.infrastructure-administrator.description',
  },
]

const FALLBACK = {
  labelKey: 'project.roles.fallback.label',
  read: false,
  write: false,
  requestApproval: false,
  grantApproval: false,
  resourcesKey: 'project.roles.fallback.resources',
  descriptionKey: null,
}

export const rolePreset = id => ROLE_PRESETS.find(r => r.id === id) || FALLBACK
