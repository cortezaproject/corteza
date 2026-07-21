// Project member roles. Each is a fixed preset: assigning a role to a member
// sets these capability flags (they are not edited per member).
//
// Project Members are human beings with roles to either (1) build part or all of
// the project and its AI systems, or (2) manage project governance and approve
// (or send back) the project's publishes.
//
// `resources` summarises what the role can reach; `description` is the role's
// responsibility statement — together they form the accountability framework
// (EU AI Act Article 17(m)).
//
// Tab visibility (see Wizard): a member who can request OR grant approval works
// in the Governance view; everyone else gets the Build view only.

// `labelKey`/`resourcesKey`/`descriptionKey` are i18n keys; components resolve
// them with $t for display. Capability flags drive behaviour and stay literal.
// `descriptionFreeKey` (optional) replaces `descriptionKey` in free-mode
// projects, where publish approval does not exist and the standard copy would
// reference it.
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
    descriptionFreeKey: 'project.roles.developer.descriptionFree',
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
    descriptionFreeKey: 'project.roles.junior-developer.descriptionFree',
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

// Approval-centric presets only make sense where approval exists at all —
// Free-mode projects publish directly, with no submit/approve step, so these
// presets are hidden from the Free-mode role picker (Gated-mode picker shows
// all of ROLE_PRESETS).
const APPROVAL_PRESET_IDS = ['governance-owner', 'security-owner', 'executive-authority']

// Presets to offer for a given build mode. `extraIds` keeps any already-
// assigned preset selectable even if it would otherwise be hidden (e.g. a
// member holding an approval-centric role, defensively — mode is immutable
// after a project is created, so this should not occur in practice).
export const rolePresetsForMode = (mode, extraIds = []) => {
  if (mode !== 'free') return ROLE_PRESETS
  const keep = new Set(extraIds)
  return ROLE_PRESETS.filter(r => keep.has(r.id) || !APPROVAL_PRESET_IDS.includes(r.id))
}
