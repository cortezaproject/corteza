import { Apply, HumanID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'

export type ProjectStatus =
  | 'draft'
  | 'active'
  | 'published'
  | 'archived'
  | 'suspended'

export type ProjectVisibility = 'open' | 'invite-only'

export type ProjectMode = 'free' | 'gated'

export type ProjectMemberRole =
  | 'governance-owner'
  | 'security-owner'
  | 'developer'
  | 'junior-developer'
  | 'member'
  | 'executive-authority'
  | 'infrastructure-administrator'

export type ProjectGovernanceStatus =
  | 'draft'
  | 'submitted'
  | 'approved'
  | 'changes-requested'

export type ProjectGovernanceAction =
  | 'submit'
  | 'approve'
  | 'request-changes'
  | 'reopen'
  | 'recall'

export interface ProjectGovernanceStep {
  values: Record<string, unknown>;
  status: ProjectGovernanceStatus;
  reviewNote: string;
}

export type ProjectGovernance = Record<string, ProjectGovernanceStep>

export interface ProjectDeployerCategories {
  publicAuthorityAnnex3: boolean;
  privateEssentialServices: boolean;
  insuranceBanking: boolean;
}

export interface ProjectPermittedConnection {
  id: string;
  name: string;
  connector?: string;
  type?: string;
  description?: string;
  actionIfUnavailable?: string;
  replacement?: string;
  isAiSystem?: string;
}

export interface ProjectResourceManagement {
  ai: Record<string, unknown>;
  infra: Record<string, unknown>;
  connections: ProjectPermittedConnection[];
}

interface Config {
  visibility?: ProjectVisibility;
  defaultMemberRole?: ProjectMemberRole;
  featureFlags?: Record<string, boolean>;

  namespaceID: string;
  deployerCategories: ProjectDeployerCategories;
  friaRequired: boolean;
  resourceManagement: ProjectResourceManagement;
}

const defaultConfig = (): Config => ({
  namespaceID: NoID,
  deployerCategories: {
    publicAuthorityAnnex3: false,
    privateEssentialServices: false,
    insuranceBanking: false,
  },
  friaRequired: false,
  resourceManagement: { ai: {}, infra: {}, connections: [] },
})

interface Meta {
  short: string;
  description: string;
  icon?: string;
  color?: string;
  tags?: string[];
}

const defaultMeta = (): Meta => ({
  short: '',
  description: '',
})

interface PartialProject extends Partial<
  Omit<Project, 'createdAt' | 'updatedAt' | 'deletedAt'>
> {
  createdAt?: string | number | Date;
  updatedAt?: string | number | Date;
  deletedAt?: string | number | Date;
}

export class Project {
  public projectID = NoID
  public tenantID = NoID
  public handle = ''
  public status: ProjectStatus = 'draft'

  // Chosen at creation, immutable after. Top-level column (not config JSON).
  public mode: ProjectMode = 'free'

  // Revision chain. Originals leave these at NoID/0; rootProjectID falls back
  // to the project's own ID (mirrors the backend RootProjectID()).
  public rootProjectID = NoID
  public parentRevisionID = NoID
  public revision = 0

  public config: Config = defaultConfig()
  public meta: Meta = defaultMeta()
  public governance: ProjectGovernance = {}
  public labels: object = {}

  // Payload capability flags
  public canGrant = false
  public canUpdateProject = false
  public canDeleteProject = false
  public canManageMembers = false

  public createdBy = NoID
  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined

  constructor(p?: PartialProject) {
    this.apply(p)
  }

  apply(p?: PartialProject): void {
    Apply(this, p, HumanID, 'projectID', 'tenantID', 'createdBy', 'rootProjectID', 'parentRevisionID')
    Apply(this, p, String, 'handle', 'status', 'mode')
    Apply(this, p, Number, 'revision')
    Apply(this, p, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')

    // An original revision carries no rootProjectID; it is its own root.
    if (this.rootProjectID === NoID) {
      this.rootProjectID = this.projectID
    }
    Apply(
      this,
      p,
      Boolean,
      'canGrant',
      'canUpdateProject',
      'canDeleteProject',
      'canManageMembers',
    )

    if (IsOf(p, 'config')) {
      this.config = { ...defaultConfig(), ...p.config }
      this.config.namespaceID = HumanID(this.config.namespaceID)
      this.config.resourceManagement = {
        ai: {},
        infra: {},
        connections: [],
        ...(p.config.resourceManagement || {}),
      }
    }

    if (IsOf(p, 'meta')) {
      this.meta = { ...defaultMeta(), ...p.meta }
    }

    if (IsOf(p, 'governance')) {
      this.governance = { ...p.governance }
    }

    if (IsOf(p, 'labels')) {
      this.labels = { ...p.labels }
    }
  }

  /**
   * Display name; the handle is the fallback
   */
  get name(): string {
    return this.meta.short || this.handle
  }

  get isGated(): boolean {
    return this.mode === 'gated'
  }

  get namespaceID(): string {
    return this.config.namespaceID
  }

  /**
   * Whether the project has its compose namespace yet (created server-side).
   * namespaceID is NoID until then, so guard on this rather than truthiness.
   */
  get hasNamespace(): boolean {
    return this.config.namespaceID !== NoID
  }

  get friaRequired(): boolean {
    return this.config.friaRequired
  }

  /**
   * Governance entry for a step; a virgin step is an editable draft
   */
  governanceStep(key: string): ProjectGovernanceStep {
    return this.governance[key] || { values: {}, status: 'draft', reviewNote: '' }
  }

  get resourceID(): string {
    return `${this.resourceType}:${this.projectID}`
  }

  get resourceType(): string {
    return 'system:project'
  }

  get fts(): string {
    return [this.meta.short, this.handle, this.projectID].join(' ').toLocaleLowerCase()
  }

  clone(): Project {
    return new Project(JSON.parse(JSON.stringify(this)))
  }
}

export interface ProjectCapabilities {
  canRead: boolean;
  canWrite: boolean;
  canRequestApproval: boolean;
  canGrantApproval: boolean;
}

interface PartialProjectMember extends Partial<
  Omit<ProjectMember, 'createdAt' | 'updatedAt' | 'deletedAt'>
> {
  createdAt?: string | number | Date;
  updatedAt?: string | number | Date;
  deletedAt?: string | number | Date;
}

export class ProjectMember {
  public projectMemberID = NoID
  public projectID = NoID
  public tenantID = NoID
  public userID = NoID
  public rolePreset: ProjectMemberRole = 'member'
  public invitedBy = NoID

  // Derived server-side from the role preset; never stored
  public capabilities: ProjectCapabilities = {
    canRead: false,
    canWrite: false,
    canRequestApproval: false,
    canGrantApproval: false,
  }

  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined

  constructor(m?: PartialProjectMember) {
    this.apply(m)
  }

  apply(m?: PartialProjectMember): void {
    Apply(this, m, HumanID, 'projectMemberID', 'projectID', 'tenantID', 'userID', 'invitedBy')
    Apply(this, m, String, 'rolePreset')
    Apply(this, m, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')

    if (IsOf(m, 'capabilities')) {
      this.capabilities = { ...this.capabilities, ...m.capabilities }
    }
  }

  get resourceID(): string {
    return `${this.resourceType}:${this.projectMemberID}`
  }

  get resourceType(): string {
    return 'system:project-member'
  }

  clone(): ProjectMember {
    return new ProjectMember(JSON.parse(JSON.stringify(this)))
  }
}

export interface ProjectGroupEntry {
  projectGroupID: string;
  resourceRef: string;
  createdAt?: string;
}

interface GroupMeta {
  short: string;
  description: string;
}

interface PartialProjectGroup extends Partial<
  Omit<ProjectGroup, 'createdAt' | 'updatedAt' | 'deletedAt'>
> {
  createdAt?: string | number | Date;
  updatedAt?: string | number | Date;
  deletedAt?: string | number | Date;
}

export class ProjectGroup {
  public projectGroupID = NoID
  public projectID = NoID
  public tenantID = NoID
  public handle = ''
  public meta: GroupMeta = { short: '', description: '' }

  public entries: ProjectGroupEntry[] = []

  // Payload capability flags
  public canUpdate = false
  public canDelete = false
  public canManageMembers = false

  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined

  constructor(g?: PartialProjectGroup) {
    this.apply(g)
  }

  apply(g?: PartialProjectGroup): void {
    Apply(this, g, HumanID, 'projectGroupID', 'projectID', 'tenantID')
    Apply(this, g, String, 'handle')
    Apply(this, g, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, g, Boolean, 'canUpdate', 'canDelete', 'canManageMembers')

    if (IsOf(g, 'meta')) {
      this.meta = { short: '', description: '', ...g.meta }
    }

    if (g?.entries) {
      this.entries = g.entries.map(e => ({ ...e, projectGroupID: HumanID(e.projectGroupID) }))
    }
  }

  get name(): string {
    return this.meta.short || this.handle
  }

  /**
   * Plain resource references carried by the group entries
   */
  get resourceRefs(): string[] {
    return this.entries.map(e => e.resourceRef)
  }

  get resourceID(): string {
    return `${this.resourceType}:${this.projectGroupID}`
  }

  get resourceType(): string {
    return 'system:project-group'
  }

  clone(): ProjectGroup {
    return new ProjectGroup(JSON.parse(JSON.stringify(this)))
  }
}
