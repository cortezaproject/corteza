import { PageBlock } from '../base'
import Feed, { FeedInput } from './feed'
import { RecordFeed } from './feed-record'
type Bounds = number[][]
interface Options {
  defaultView: string
  center: Array<number>
  feeds: Array<Feed>
  zoomStarting: number
  zoomMin: number
  zoomMax: number
  bounds: Bounds | null
  lockBounds: boolean
  refreshRate: number
  showRefresh: boolean
  magnifyOption: string
  displayOption: string
}
export declare class PageBlockGeometry extends PageBlock {
  readonly kind = 'Geometry'
  options: Options
  static feedResources: Readonly<{
    record: 'compose:record'
  }>
  constructor(i?: PageBlock | Partial<PageBlock>)
  applyOptions(o?: Partial<Options>): void
  static makeFeed(f?: FeedInput): Feed
  static RecordFeed: typeof RecordFeed
}
export {}
