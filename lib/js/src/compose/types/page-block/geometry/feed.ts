import { Apply, NoID } from '../../../../cast'
import { IsOf } from '../../../../guards'

interface FeedOptions {
  color: string;
  prefilter: string;
  moduleID: string;
  resource: string;
  titleField: string;
  geometryField: string;
  displayMarker: boolean;
  displayPolygon: boolean;
}

export type FeedInput = Partial<Feed> | Feed

const defOptions = {
  moduleID: NoID,
  color: '#2f85cb',
  prefilter: '',
  resource: 'compose:record',
  titleField: '',
  geometryField: '',
  displayMarker: false,
  displayPolygon: true,
}

/**
 * Feed class represents an event feed for the given calendar
 */
export default class Feed {
  public resource = 'compose:record'
  public titleField = ''
  public color = '#2f85cb'
  public geometryField = ''
  public displayMarker = false
  public displayPolygon = true
  public options: FeedOptions = { ...defOptions }

  constructor (i?: FeedInput) {
    this.apply(i)
  }

  apply (i?: FeedInput): void {
    if (!i) return

    if (IsOf<Feed>(i, 'resource')) {
      Apply(this, i, String, 'resource', 'color', 'titleField', 'geometryField')
      Apply(this, i, Boolean, 'displayMarker', 'displayPolygon')

      if (i.options) {
        this.options = { ...this.options, ...i.options }
      }
    }
  }
}
