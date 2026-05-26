// Project member roles. Each is a fixed preset: assigning a role to a member
// sets these capability flags (they are not edited per member).
//
// Tab visibility (see Wizard): a member who can request OR grant approval works
// in the Governance view; everyone else gets the Build view only.

export const ROLE_PRESETS = [
  {
    id: 'developer',
    label: 'Developer',
    read: true,
    write: true,
    requestApproval: true,
    grantApproval: false,
  },
  {
    id: 'governance-owner',
    label: 'Governance Owner',
    read: true,
    write: true,
    requestApproval: false,
    grantApproval: true,
  },
  {
    id: 'security-owner',
    label: 'Security Owner',
    read: true,
    write: true,
    requestApproval: false,
    grantApproval: true,
  },
  {
    id: 'junior-developer',
    label: 'Junior Developer',
    read: true,
    write: true,
    requestApproval: false,
    grantApproval: false,
  },
]

const FALLBACK = { label: 'Member', read: false, write: false, requestApproval: false, grantApproval: false }

export const rolePreset = id => ROLE_PRESETS.find(r => r.id === id) || FALLBACK
