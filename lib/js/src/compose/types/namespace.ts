import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'

interface MetaAdminRecordList {
  columns: string[];
}

interface MetaAdmin {
  recordList: MetaAdminRecordList;
}

interface Meta {
  subtitle: string;
  description: string;
  hideSidebar: boolean;
  // Temporary icon & logo URLs
  // @todo rework this when we rework attachment management
  icon: string;
  logo: string;
  logoEnabled: boolean;
}

interface PartialNamespace extends Partial<Omit<Namespace, 'meta' | 'createdAt' | 'updatedAt' | 'deletedAt'>> {
  meta?: Partial<Meta>;
  createdAt?: string|number|Date;
  updatedAt?: string|number|Date;
  deletedAt?: string|number|Date;
}

export class Namespace {
  public namespaceID = NoID
  public name = ''
  public slug = ''

  public enabled = false

  public labels: object = {}

  public meta: object = {}

  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined

  public canCreateChart = false
  public canCreateModule = false
  public canCreatePage = false
  public canDeleteNamespace = false
  public canUpdateNamespace = false
  public canManageNamespace = false
  public canCloneNamespace = false
  public canExportNamespace = false
  public canGrant = false
  public canExportCharts = false
  public canExportModules = false

  constructor (i?: PartialNamespace) {
    this.apply(i)
  }

  clone (): Namespace {
    return new Namespace(JSON.parse(JSON.stringify(this)))
  }

  apply (n?: PartialNamespace | Namespace): void {
    if (!n) return

    Apply(this, n, CortezaID, 'namespaceID')
    Apply(this, n, String, 'name', 'slug')

    Apply(this, n, Boolean, 'enabled')

    if (IsOf(n, 'meta')) {
      this.meta = { ...n.meta }
    }

    if (IsOf(n, 'labels')) {
      this.labels = { ...n.labels }
    }

    Apply(this, n, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, n, Boolean,
      'canDeleteNamespace',
      'canUpdateNamespace',
      'canManageNamespace',
      'canCloneNamespace',
      'canExportNamespace',
      'canGrant',
      'canCreateModule',
      'canExportModules',
      'canCreatePage',
      'canCreateChart',
      'canExportCharts',
    )
  }

  /**
   * Returns resource ID
   */
  get resourceID (): string {
    return `${this.resourceType}:${this.namespaceID}`
  }

  /**
   * Resource type
   */
  get resourceType (): string {
    return 'compose:namespace'
  }

  /**
   * Calculate namespace initials
   */
  get initials (): string {
    let base = this.name || this.slug

    // if length is shorter than 3 letters, use that
    if (base.length <= 3) {
      return base
    }

    // split by space and take first letter of each word
    base = base.split(/\s+/).map(w => w[0]).filter(c => /[a-zA-Z]/.test(c)).join('')
    if (base.length > 3) {
      base = base.slice(0, 3)
    }

    return base
  }
}
