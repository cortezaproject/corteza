interface Meta {
    subtitle: string;
    description: string;
    hideSidebar: boolean;
    icon: string;
    logo: string;
    logoEnabled: boolean;
}
interface PartialNamespace extends Partial<Omit<Namespace, 'meta' | 'createdAt' | 'updatedAt' | 'deletedAt'>> {
    meta?: Partial<Meta>;
    createdAt?: string | number | Date;
    updatedAt?: string | number | Date;
    deletedAt?: string | number | Date;
}
export declare class Namespace {
  namespaceID: string
  name: string
  slug: string
  enabled: boolean
  labels: object
  meta: object
  createdAt?: Date
  updatedAt?: Date
  deletedAt?: Date
  canCreateChart: boolean
  canCreateModule: boolean
  canCreatePage: boolean
  canDeleteNamespace: boolean
  canUpdateNamespace: boolean
  canManageNamespace: boolean
  canCloneNamespace: boolean
  canExportNamespace: boolean
  canGrant: boolean
  canExportCharts: boolean
  canExportModules: boolean
  constructor(i?: PartialNamespace);
  clone(): Namespace;
  apply(n?: PartialNamespace | Namespace): void;
  /**
     * Returns resource ID
     */
  get resourceID(): string;
  /**
     * Resource type
     */
  get resourceType(): string;
  /**
     * Calculate namespace initials
     */
  get initials(): string;
}
export {}
