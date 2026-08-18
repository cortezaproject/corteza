import { Apply } from '../../../../cast'
import { PageBlock, Registry, BlockIssue, isUnsetID } from '../base'
import Feed, { FeedInput } from './feed'
import { RecordFeed } from './feed-record'

const kind = 'Geometry'

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
  hideGeoSearch: boolean
}

const defaults: Readonly<Options> = Object.freeze({
  defaultView: '',
  center: [35, -30],
  feeds: [],
  zoomStarting: 2,
  zoomMin: 1,
  zoomMax: 18,
  bounds: null,
  lockBounds: false,
  refreshRate: 0,
  showRefresh: false,
  magnifyOption: '',
  displayOption: 'sameTab',
  hideGeoSearch: true,
})

export class PageBlockGeometry extends PageBlock {
  readonly kind = kind
  options: Options = { ...defaults }

  static feedResources = Object.freeze({
    record: 'compose:record',
  })

  constructor(i?: PageBlock | Partial<PageBlock>) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions(o?: Partial<Options>): void {
    if (!o) return

    this.options.feeds = (o.feeds || []).map(f => new Feed(f))
    this.options.center = o.center || []
    this.options.bounds = o.bounds || null

    Apply(this.options, o, String, 'magnifyOption', 'displayOption')
    Apply(this.options, o, Number, 'zoomStarting', 'zoomMin', 'zoomMax', 'refreshRate')
    Apply(this.options, o, Boolean, 'lockBounds', 'showRefresh', 'hideGeoSearch')
  }

  static makeFeed(f?: FeedInput): Feed {
    return new Feed(f)
  }

  static RecordFeed = RecordFeed

  validate(): Array<BlockIssue> {
    const ee = super.validate()

    if (!this.options.feeds.length) {
      ee.push({ option: 'feeds', labelKey: 'block.issue.noFeeds' })
      return ee
    }

    this.options.feeds.forEach((f, i) => {
      if (isUnsetID(f.options?.moduleID)) {
        ee.push({ option: `feeds.${i}.moduleID`, labelKey: 'block.issue.feedModule' })
      }
      if (!f.geometryField) {
        ee.push({ option: `feeds.${i}.geometryField`, labelKey: 'block.issue.feedGeometryField' })
      }
    })

    return ee
  }
}

Registry.set(kind, PageBlockGeometry)
