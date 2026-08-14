import { Apply, HumanID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'

export type ProjectStatus = 'draft' | 'active' | 'published' | 'archived' | 'suspended'

// Where a revision stands in the publish approval cycle — a separate axis from
// ProjectStatus, which says what the revision IS rather than whether anyone has
// agreed to put it live.
export type ProjectApprovalStatus = 'draft' | 'submitted' | 'approved' | 'rejected'

export type ProjectVisibility = 'open' | 'invite-only'

export type ProjectMemberRole =
  | 'governance-owner'
  | 'security-owner'
  | 'developer'
  | 'junior-developer'
  | 'member'
  | 'executive-authority'
  | 'infrastructure-administrator'

export interface ProjectDeployerCategories {
  publicAuthorityAnnex3: boolean
  privateEssentialServices: boolean
  insuranceBanking: boolean
}

export interface ProjectPermittedConnection {
  id: string
  name: string
  connector?: string
  type?: string
  description?: string
  actionIfUnavailable?: string
  replacement?: string
  isAiSystem?: string
}

export interface ProjectResourceManagement {
  ai: Record<string, unknown>
  infra: Record<string, unknown>
  connections: ProjectPermittedConnection[]
}

interface Config {
  visibility?: ProjectVisibility
  defaultMemberRole?: ProjectMemberRole
  featureFlags?: Record<string, boolean>

  namespaceID: string
  deployerCategories: ProjectDeployerCategories
  friaRequired: boolean
  resourceManagement: ProjectResourceManagement
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
  short: string
  description: string
  icon?: string
  color?: string
  tags?: string[]
}

const defaultMeta = (): Meta => ({
  short: '',
  description: '',
})

interface PartialProject extends Partial<Omit<Project, 'createdAt' | 'updatedAt' | 'deletedAt'>> {
  createdAt?: string | number | Date
  updatedAt?: string | number | Date
  deletedAt?: string | number | Date
}

export class Project {
  public projectID = NoID
  public tenantID = NoID
  public handle = ''
  public status: ProjectStatus = 'draft'

  // Revision chain. Originals leave these at NoID/0; rootProjectID falls back
  // to the project's own ID (mirrors the backend RootProjectID()).
  public rootProjectID = NoID
  public parentRevisionID = NoID
  public revision = 0

  public config: Config = defaultConfig()
  public meta: Meta = defaultMeta()
  public labels: object = {}

  // Payload capability flags
  public canGrant = false
  public canUpdateProject = false
  public canDeleteProject = false
  public canManageMembers = false
  // Branching and publishing are their own RBAC operations, not update: a
  // deployment can grant "may edit the draft" without "may put it live".
  public canReviseProject = false
  public canPublishProject = false

  public createdBy = NoID
  // Stamped by the backend on every update, including the one that publishes a
  // draft — which is what makes it the "published by" on a publish receipt.
  public updatedBy = NoID
  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined
  // Shelf state, deliberately not a status value: status is owned by publish
  // and revision, and archiving is the one lifecycle change a user makes
  // directly. Set means archived.
  public archivedAt?: Date = undefined

  // Server-owned publish approval, like status: the backend refuses to publish
  // anything not 'approved' here. approvalPlan fingerprints the deployment plan
  // the approval was granted against; publish recomputes it and sends the
  // revision back for re-approval on a mismatch.
  public approvalStatus: ProjectApprovalStatus = 'draft'
  public approvalPlan = ''
  public approvalNote = ''
  public approvalSubmittedBy = NoID
  public approvalSubmittedAt?: Date = undefined
  public approvalDecidedBy = NoID
  public approvalDecidedAt?: Date = undefined

  constructor(p?: PartialProject) {
    this.apply(p)
  }

  apply(p?: PartialProject): void {
    Apply(
      this,
      p,
      HumanID,
      'projectID',
      'tenantID',
      'createdBy',
      'updatedBy',
      'rootProjectID',
      'parentRevisionID',
      'approvalSubmittedBy',
      'approvalDecidedBy',
    )
    Apply(this, p, String, 'handle', 'status', 'approvalStatus', 'approvalPlan', 'approvalNote')
    Apply(this, p, Number, 'revision')
    Apply(
      this,
      p,
      ISO8601Date,
      'createdAt',
      'updatedAt',
      'deletedAt',
      'archivedAt',
      'approvalSubmittedAt',
      'approvalDecidedAt',
    )

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
      'canReviseProject',
      'canPublishProject',
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
  canRead: boolean
  canWrite: boolean
  canRequestApproval: boolean
  canGrantApproval: boolean
}

interface PartialProjectMember extends Partial<
  Omit<ProjectMember, 'createdAt' | 'updatedAt' | 'deletedAt'>
> {
  createdAt?: string | number | Date
  updatedAt?: string | number | Date
  deletedAt?: string | number | Date
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

export interface ProjectAiSystemEntry {
  projectAiSystemID: string
  resourceRef: string
  createdAt?: string
}

/**
 * EU AI Act Art. 6 risk classification of an AI system
 */
export type ProjectAiSystemRiskClass = 'prohibited' | 'high' | 'limited' | 'minimal'

interface AiSystemMeta {
  short: string
  description: string
  intendedPurpose: string
}

interface PartialProjectAiSystem extends Partial<
  Omit<ProjectAiSystem, 'createdAt' | 'updatedAt' | 'deletedAt'>
> {
  createdAt?: string | number | Date
  updatedAt?: string | number | Date
  deletedAt?: string | number | Date
}

export class ProjectAiSystem {
  public projectAiSystemID = NoID
  public projectID = NoID
  public tenantID = NoID
  public handle = ''
  public riskClass = ''
  public meta: AiSystemMeta = { short: '', description: '', intendedPurpose: '' }

  public entries: ProjectAiSystemEntry[] = []

  // Payload capability flags
  public canUpdate = false
  public canDelete = false
  public canManageResources = false

  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined

  constructor(g?: PartialProjectAiSystem) {
    this.apply(g)
  }

  apply(g?: PartialProjectAiSystem): void {
    Apply(this, g, HumanID, 'projectAiSystemID', 'projectID', 'tenantID')
    Apply(this, g, String, 'handle', 'riskClass')
    Apply(this, g, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, g, Boolean, 'canUpdate', 'canDelete', 'canManageResources')

    if (IsOf(g, 'meta')) {
      this.meta = { short: '', description: '', intendedPurpose: '', ...g.meta }
    }

    if (g?.entries) {
      this.entries = g.entries.map(e => ({ ...e, projectAiSystemID: HumanID(e.projectAiSystemID) }))
    }
  }

  get name(): string {
    return this.meta.short || this.handle
  }

  /**
   * Plain resource references carried by the AI system entries
   */
  get resourceRefs(): string[] {
    return this.entries.map(e => e.resourceRef)
  }

  get resourceID(): string {
    return `${this.resourceType}:${this.projectAiSystemID}`
  }

  get resourceType(): string {
    return 'system:project-ai-system'
  }

  clone(): ProjectAiSystem {
    return new ProjectAiSystem(JSON.parse(JSON.stringify(this)))
  }
}
