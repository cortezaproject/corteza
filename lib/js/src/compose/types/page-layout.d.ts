import { PageBlock } from './page-block/base'
import { Button } from './page-block/types'
export type PageLayoutInput = PageLayout | Partial<PageLayout>
interface PageLayoutConfig {
  visibility: Visibility
  buttons: {
    back: Button
    delete: Button
    new: Button
    clone: Button
    edit: Button
    submit: Button
  }
  actions: Action[]
  useTitle: boolean
}
interface Action {
  kind: string
  enabled: boolean
  placement: string
  params: unknown
  meta: ActionMeta
}
interface ActionMeta {
  label: string
  style: {
    variant: string
  }
}
interface Visibility {
  expression: string
  roles: string[]
}
interface Meta {
  title: string
  description: string
}
export declare class PageLayout {
  pageLayoutID: string
  namespaceID: string
  pageID: string
  handle: string
  weight: number
  blocks: Partial<PageBlock>[]
  config: PageLayoutConfig
  meta: Meta
  createdAt?: Date
  updatedAt?: Date
  deletedAt?: Date
  ownedBy: string
  constructor(pl?: PageLayoutInput)
  apply(pl?: PageLayoutInput): void
  clone(): PageLayout
  addAction(): void
  /**
   * Returns resource ID
   */
  get resourceID(): string
  /**
   * Resource type
   */
  get resourceType(): string
  export(): PageLayoutInput
}
export {}
