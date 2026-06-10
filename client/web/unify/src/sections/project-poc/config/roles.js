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

export const ROLE_PRESETS = [
  {
    id: 'developer',
    label: 'Developer',
    read: true,
    write: true,
    requestApproval: true,
    grantApproval: false,
    resources: 'All Pages, All Gates',
    description:
      'Builds part or all of the project and its AI systems, and requests governance approval at the gates.',
  },
  {
    id: 'governance-owner',
    label: 'Governance Owner',
    read: true,
    write: true,
    requestApproval: false,
    grantApproval: true,
    resources: 'All Pages, All Gates',
    description: 'Owns project governance and grants approvals at every governance gate.',
  },
  {
    id: 'security-owner',
    label: 'Security Owner',
    read: true,
    write: true,
    requestApproval: false,
    grantApproval: true,
    resources: 'RBAC, Security, Gate 6',
    description:
      'Responsible for RBAC and security configuration; grants approval at the security gate (Gate 6).',
  },
  {
    id: 'junior-developer',
    label: 'Junior Developer',
    read: true,
    write: true,
    requestApproval: false,
    grantApproval: false,
    resources: 'Technical Architecture, Data Model, Connections, Automations',
    description:
      'Builds the technical architecture, data model, connections and automations. Cannot request or grant approvals.',
  },
  {
    id: 'executive-authority',
    label: 'Executive Authority',
    read: true,
    write: false,
    requestApproval: false,
    grantApproval: true,
    resources: 'Read all, Approve at Gate 1',
    description:
      'Has read access across the project and signs off the project at Gate 1, once the Project Summary, Quality Management System and Project Members have been defined.',
  },
  {
    id: 'infrastructure-administrator',
    label: 'Infrastructure Administrator',
    read: false,
    write: false,
    requestApproval: false,
    grantApproval: false,
    resources: '—',
    description:
      'Responsible for setting up the server infrastructure, and installing and maintaining the platform software, including platform-level software updates.',
  },
]

const FALLBACK = {
  label: 'Member',
  read: false,
  write: false,
  requestApproval: false,
  grantApproval: false,
  resources: '—',
  description: '',
}

export const rolePreset = id => ROLE_PRESETS.find(r => r.id === id) || FALLBACK
