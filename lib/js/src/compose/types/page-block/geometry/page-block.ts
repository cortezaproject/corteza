import { PageBlock, Registry } from '../base'
import { Apply } from '../../../../cast'
import Feed, { FeedInput } from './feed'
import { RecordFeed } from './feed-record'

const kind = 'Geometry'

type Bounds = number[][]

interface Options {
  defaultView: string;
  center: Array<number>;
  feeds: Array<Feed>;
  zoomStarting: number;
  zoomMin: number;
  zoomMax: number;
  bounds: Bounds | null;
  lockBounds: boolean;
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
})

export class PageBlockGeometry extends PageBlock {
  readonly kind = kind
  options: Options = { ...defaults }

  static feedResources = Object.freeze({
    record: 'compose:record',
  })

  constructor (i?: PageBlock | Partial<PageBlock>) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    this.options.feeds = (o.feeds || []).map(f => new Feed(f))
    this.options.center = (o.center || [])
    this.options.bounds = (o.bounds || null)

    Apply(this.options, o, Number, 'zoomStarting', 'zoomMin', 'zoomMax')
    Apply(this.options, o, Boolean, 'lockBounds')
  }

  static makeFeed (f?: FeedInput): Feed {
    return new Feed(f)
  }

  static RecordFeed = RecordFeed
}

Registry.set(kind, PageBlockGeometry)
