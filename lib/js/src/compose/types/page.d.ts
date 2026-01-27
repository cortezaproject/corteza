import { PageBlock } from './page-block'
interface PartialPage extends Partial<Omit<Page, 'children' | 'meta' | 'blocks' | 'createdAt' | 'updatedAt' | 'deletedAt'>> {
    children?: Array<PartialPage>;
    blocks?: PageBlock[];
    meta?: PageMeta;
    createdAt?: string | number | Date;
    updatedAt?: string | number | Date;
    deletedAt?: string | number | Date;
}
interface PageMeta {
    notifications: {
        enabled: boolean;
    };
}
interface PageConfig {
    navItem: {
        icon: {
            type: string;
            src: string;
        };
        expanded: false;
    };
}
export declare class Page {
  pageID: string
  selfID: string
  moduleID: string
  namespaceID: string
  title: string
  handle: string
  description: string
  weight: number
  labels: object
  visible: boolean
  children?: Page[]
  blocks: PageBlock[]
  config: PageConfig
  meta: PageMeta
  createdAt?: Date
  updatedAt?: Date
  deletedAt?: Date
  canUpdatePage: boolean
  canDeletePage: boolean
  canGrant: boolean
  constructor(i?: PartialPage);
  clone(): Page;
  apply(i?: PartialPage): void;
  /**
     * Returns resource ID
     */
  get resourceID(): string;
  /**
     * Resource type
     */
  get resourceType(): string;
  get isRecordPage(): boolean;
  get firstLevel(): boolean;
  /**
     * Validates page & it's blocks
     */
  validate(): Array<string>;
  export(): PartialPage;
}
export {}
