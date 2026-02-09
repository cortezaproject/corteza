import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'
import type { Typed } from './values'

/**
 * NgAutomation metadata interface
 */
interface NgAutomationMeta {
  short: string
  description?: string
}

/**
 * Expression for step arguments/results (matches backend Expr type)
 */
export interface Expr {
  argumentName?: string
  target?: string
  type: string
  value?: string | number | boolean | object | null
  expr?: string // Expression string for dynamic values
  scope?: string
}

/**
 * Trigger constraint for filtering events
 */
export interface TriggerConstraint {
  name: string
  op: string // Operator: eq, ne, like, etc.
  value: string | string[]
}

/**
 * NgAutomationTrigger - defines what initiates an automation
 * Aligned with backend JSON structure
 */
export interface NgAutomationTrigger {
  triggerID: string
  handle?: string // Human-readable identifier (e.g., trigger_1)
  enabled: boolean
  resourceType: string
  eventType: string
  constraints?: TriggerConstraint[]
  input?: Record<string, unknown>
  meta?: { short?: string; description?: string }
  createdAt?: string
  updatedAt?: string
  deletedAt?: string
}

/**
 * NgAutomationStep - defines an action in the automation flow
 * Aligned with backend JSON structure
 */
export interface NgAutomationStep {
  stepID: string
  handle?: string // Human-readable identifier (e.g., step_1)
  kind: string // Function kind: 'function', 'iterator', 'gateway', etc.
  ref: string // Function reference
  meta?: { short?: string; description?: string }
  arguments?: Expr[]
  results?: Expr[]
}

/**
 * NgAutomationPath - defines connections between steps
 * Aligned with backend JSON structure (uses IDs)
 */
export interface NgAutomationPath {
  parentID: string // Source step ID
  childID: string // Target step ID
  handle?: string // Human-readable identifier (e.g., path_1_2)
  meta?: { short?: string; description?: string; expr?: string }
}

/**
 * TAQ issue interface
 */
export interface TAQIssue {
  message: string
  severity?: 'warning' | 'error'
}

export type TAQIssueSet = TAQIssue[]

/**
 * Partial NgAutomation interface for constructor/apply
 */
type PartialNgAutomation = Partial<
  Omit<NgAutomation, 'createdAt' | 'updatedAt' | 'deletedAt' | 'meta'>
> & {
  automationID?: string
  meta?: Partial<NgAutomationMeta>
  createdAt?: string | number | Date
  updatedAt?: string | number | Date
  deletedAt?: string | number | Date
  triggers?: NgAutomationTrigger[]
  steps?: NgAutomationStep[]
  paths?: NgAutomationPath[]
}

/**
 * NgAutomation (Next-Gen Automation) class
 *
 * Represents a workflow automation definition with triggers, steps, and paths.
 * Internal naming follows backend convention; TAQ is kept as user-facing alias.
 */
export class NgAutomation {
  public automationID = NoID
  public handle = ''
  public labels: Record<string, string> = {}
  public meta: NgAutomationMeta = { short: '' }
  public enabled = false

  public trace = false

  public triggers: NgAutomationTrigger[] = []
  public steps: NgAutomationStep[] = []
  public paths: NgAutomationPath[] = []

  // How much time to keep completed sessions (in seconds)
  public keepSessions = 0

  // Initial input scope
  public scope: Typed | null = null

  // Collection of issues from the last parse
  public issues: TAQIssueSet = []

  public runAs = NoID

  public ownedBy = NoID
  public createdAt?: Date = undefined
  public createdBy = NoID
  public updatedAt?: Date = undefined
  public updatedBy = NoID
  public deletedAt?: Date = undefined
  public deletedBy = NoID

  constructor(t?: PartialNgAutomation) {
    this.apply(t)
  }

  apply(t?: PartialNgAutomation): void {
    if (IsOf(t, 'automationID')) {
      Apply(this, t, CortezaID, 'automationID')
    }

    Apply(this, t, String, 'handle')

    Apply(this, t, Boolean, 'enabled', 'trace')
    Apply(this, t, Number, 'keepSessions')

    Apply(this, t, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, t, CortezaID, 'runAs', 'ownedBy', 'createdBy', 'updatedBy', 'deletedBy')

    if (IsOf(t, 'meta')) {
      this.meta = { short: '', ...t.meta }
    }

    if (IsOf(t, 'labels')) {
      this.labels = { ...t.labels }
    }

    if (IsOf(t, 'triggers')) {
      this.triggers = [...(t.triggers || [])]
    }

    if (IsOf(t, 'steps')) {
      this.steps = [...(t.steps || [])]
    }

    if (IsOf(t, 'paths')) {
      this.paths = [...(t.paths || [])]
    }

    if (IsOf(t, 'scope')) {
      this.scope = t.scope ?? null
    }

    if (IsOf(t, 'issues')) {
      this.issues = [...(t.issues || [])]
    }
  }

  /**
   * Returns resource ID
   */
  get resourceID(): string {
    return `${this.resourceType}:${this.automationID}`
  }

  /**
   * Resource type
   */
  get resourceType(): string {
    return 'automation:ng-automation'
  }

  /**
   * Clone the NgAutomation instance
   */
  clone(): NgAutomation {
    return new NgAutomation(this)
  }

  /**
   * Serialize to backend format for API calls
   */
  toJSON(): Record<string, unknown> {
    return {
      automationID: this.automationID,
      handle: this.handle,
      labels: this.labels,
      meta: this.meta,
      enabled: this.enabled,
      trace: this.trace,
      triggers: this.triggers,
      steps: this.steps,
      paths: this.paths,
      keepSessions: this.keepSessions,
      scope: this.scope,
      runAs: this.runAs,
      ownedBy: this.ownedBy,
    }
  }
}

/**
 * TAQ is kept as a user-facing alias for NgAutomation
 */
export { NgAutomation as TAQ }
