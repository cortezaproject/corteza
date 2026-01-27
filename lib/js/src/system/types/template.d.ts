interface PartialTemplate extends Partial<Omit<Template, 'createdAt' | 'updatedAt' | 'deletedAt' | 'lastUsedAt'>> {
    createdAt?: string | number | Date;
    updatedAt?: string | number | Date;
    deletedAt?: string | number | Date;
    lastUsedAt?: string | number | Date;
}
interface Meta {
    short?: string;
    description?: string;
}
export declare class Template {
  templateID: string
  handle: string
  language: string
  type: string
  partial: boolean
  meta: Meta
  template: string
  labels: object
  ownerID: string
  createdAt?: Date
  updatedAt?: Date
  deletedAt?: Date
  lastUsedAt?: Date
  canDeleteTemplate: boolean
  constructor(r?: PartialTemplate);
  apply(r?: PartialTemplate): void;
  /**
     * Returns resource ID
     */
  get resourceID(): string;
  /**
     * Resource type
     */
  get resourceType(): string;
  clone(): Template;
}
export {}
