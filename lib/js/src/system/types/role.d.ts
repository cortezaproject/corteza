interface PartialRole extends Partial<Omit<Role, 'createdAt' | 'updatedAt' | 'deletedAt' | 'archivedAt'>> {
    createdAt?: string | number | Date;
    updatedAt?: string | number | Date;
    deletedAt?: string | number | Date;
    archivedAt?: string | number | Date;
}
interface Meta {
    description: string;
    context: MetaContext;
}
interface MetaContext {
    resourceTypes: Array<string>;
    expr: string;
}
export declare class Role {
  roleID: string
  name: string
  handle: string
  members: string[]
  labels: object
  meta: Meta
  createdAt?: Date
  updatedAt?: Date
  deletedAt?: Date
  archivedAt?: Date
  isSystem: boolean
  isClosed: boolean
  isBypass: boolean
  canGrant: boolean
  canUpdateRole: boolean
  canDeleteRole: boolean
  canManageMembersOnRole: boolean
  constructor(r?: PartialRole);
  apply(r?: PartialRole): void;
  /**
     * Returns resource ID
     */
  get resourceID(): string;
  /**
     * Resource type
     */
  get resourceType(): string;
  get isContext(): boolean;
  clone(): Role;
}
export {}
