interface Meta {
    name: '';
}
interface PartialWorkflow extends Partial<Omit<Workflow, 'createdAt' | 'updatedAt' | 'deletedAt' | 'suspendedAt'>> {
    meta?: Partial<Meta>;
    createdAt?: string | number | Date;
    updatedAt?: string | number | Date;
    deletedAt?: string | number | Date;
}
export declare class Workflow {
  workflowID: string
  handle: string
  enabled: boolean
  labels: object
  meta: object
  runAs: string
  ownedBy: string
  createdBy: string
  createdAt?: Date
  updatedAt?: Date
  deletedAt?: Date
  constructor(w?: PartialWorkflow);
  apply(w?: PartialWorkflow): void;
  /**
     * Returns resource ID
     */
  get resourceID(): string;
  /**
     * Resource type
     */
  get resourceType(): string;
}
export {}
